// Copyright The Enterprise Contract Contributors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
// SPDX-License-Identifier: Apache-2.0

package file

import (
	"archive/zip"
	"bytes"
	"context"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/enterprise-contract/go-gather/expand"
	_ "github.com/enterprise-contract/go-gather/expand/zip" // Register zip expander
)

func TestFileGatherer_Matcher(t *testing.T) {
	fg := &FileGatherer{}

	tests := []struct {
		name string
		uri  string
		want bool
	}{
		{"file:// prefix", "file://some/path", true},
		{"file:: prefix", "file::another/path", true},
		{"absolute path", "/etc/hosts", true},
		{"relative path dot", "./myfile", true},
		{"relative path dotdot", "../myfile", true},
		{"no match", "http://example.com/file.txt", false},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := fg.Matcher(tc.uri)
			if got != tc.want {
				t.Errorf("Matcher(%q) = %v, want %v", tc.uri, got, tc.want)
			}
		})
	}
}

func TestFileGatherer_Gather_File(t *testing.T) {
	fg := &FileGatherer{}

	tempDir := t.TempDir()
	srcFile := filepath.Join(tempDir, "test.txt")
	dstFile := filepath.Join(tempDir, "dest.txt")

	content := []byte("Hello from FileGatherer!")
	if err := os.WriteFile(srcFile, content, 0600); err != nil {
		t.Fatalf("failed to create source file: %v", err)
	}

	ctx := context.Background()
	meta, err := fg.Gather(ctx, srcFile, dstFile)
	if err != nil {
		t.Fatalf("Gather returned an unexpected error: %v", err)
	}

	if _, err := os.Stat(dstFile); os.IsNotExist(err) {
		t.Fatalf("expected file %s to exist, but it does not", dstFile)
	}

	fsMeta, ok := meta.(*FSMetadata)
	if !ok {
		t.Fatalf("expected FSMetadata, got %T", meta)
	}
	if fsMeta.Path != dstFile {
		t.Errorf("expected metadata path=%s, got %s", dstFile, fsMeta.Path)
	}
	if fsMeta.Size != int64(len(content)) {
		t.Errorf("expected metadata size=%d, got %d", len(content), fsMeta.Size)
	}
	if fsMeta.Timestamp == "" {
		t.Error("expected timestamp to be set, got empty string")
	}
}

func TestFileGatherer_Gather_Directory(t *testing.T) {
	fg := &FileGatherer{}

	tempDir := t.TempDir()
	srcDir := filepath.Join(tempDir, "source_dir")
	dstDir := filepath.Join(tempDir, "dest_dir")

	if err := os.Mkdir(srcDir, 0755); err != nil {
		t.Fatalf("failed to create source directory: %v", err)
	}
	file1 := filepath.Join(srcDir, "file1.txt")
	file2 := filepath.Join(srcDir, "file2.txt")
	if err := os.WriteFile(file1, []byte("file1"), 0600); err != nil {
		t.Fatalf("failed to write file1: %v", err)
	}
	if err := os.WriteFile(file2, []byte("file2"), 0600); err != nil {
		t.Fatalf("failed to write file2: %v", err)
	}

	ctx := context.Background()
	meta, err := fg.Gather(ctx, srcDir, dstDir)
	if err != nil {
		t.Fatalf("Gather returned an unexpected error: %v", err)
	}

	copied1 := filepath.Join(dstDir, "file1.txt")
	copied2 := filepath.Join(dstDir, "file2.txt")
	if _, err := os.Stat(copied1); os.IsNotExist(err) {
		t.Fatalf("expected file %s to exist, but it does not", copied1)
	}
	if _, err := os.Stat(copied2); os.IsNotExist(err) {
		t.Fatalf("expected file %s to exist, but it does not", copied2)
	}

	fsMeta, ok := meta.(*FSMetadata)
	if !ok {
		t.Fatalf("expected FSMetadata, got %T", meta)
	}
	if fsMeta.Path != dstDir {
		t.Errorf("expected metadata path=%s, got %s", dstDir, fsMeta.Path)
	}
	if fsMeta.Size <= 0 {
		t.Errorf("expected size > 0, got %d", fsMeta.Size)
	}
	if fsMeta.Timestamp == "" {
		t.Error("expected timestamp to be set, got empty string")
	}
}

func TestFileGatherer_Gather_NotExist(t *testing.T) {
	fg := &FileGatherer{}

	tempDir := t.TempDir()
	srcFile := filepath.Join(tempDir, "no_such_file.txt")
	dstFile := filepath.Join(tempDir, "dest.txt")

	ctx := context.Background()
	_, err := fg.Gather(ctx, srcFile, dstFile)
	if err == nil {
		t.Fatal("expected an error for non-existent source, got nil")
	}
	if !strings.Contains(err.Error(), "source file does not exist") {
		t.Errorf("expected error about 'source file does not exist', got %v", err)
	}
}

func TestFileGatherer_Gather_Cancel(t *testing.T) {
	fg := &FileGatherer{}

	tempDir := t.TempDir()
	srcFile := filepath.Join(tempDir, "test.txt")
	dstFile := filepath.Join(tempDir, "dest.txt")
	if err := os.WriteFile(srcFile, []byte("some content"), 0600); err != nil {
		t.Fatalf("failed to create source file: %v", err)
	}

	// Create a context that's already canceled
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, err := fg.Gather(ctx, srcFile, dstFile)
	if err == nil {
		t.Fatal("expected an error due to context cancellation, got nil")
	}
	if ctx.Err() != context.Canceled {
		t.Errorf("expected context to be canceled, got %v", ctx.Err())
	}
}

func TestFileGatherer_Gather_Compressed(t *testing.T) {
	fg := &FileGatherer{}

	tempDir := t.TempDir()
	srcZip := filepath.Join(tempDir, "test.zip")
	dstDir := filepath.Join(tempDir, "extracted")

	if err := createZipFile(srcZip, "hello.txt", "Hello Zip"); err != nil {
		t.Fatalf("failed to create test zip file: %v", err)
	}

	if expand.GetExpander("zip") == nil {
		t.Skip("no zip expander registered. Register a ZipExpander first or remove this test.")
	}

	ctx := context.Background()
	meta, err := fg.Gather(ctx, srcZip, dstDir)
	if err != nil {
		t.Fatalf("Gather returned an unexpected error: %v", err)
	}

	extractedFile := filepath.Join(dstDir, "hello.txt")
	if _, err := os.Stat(extractedFile); os.IsNotExist(err) {
		t.Fatalf("expected %s to exist, but it does not", extractedFile)
	}

	fsMeta, ok := meta.(*FSMetadata)
	if !ok {
		t.Fatalf("expected FSMetadata, got %T", meta)
	}
	if fsMeta.Path != dstDir {
		t.Errorf("expected path=%s, got %s", dstDir, fsMeta.Path)
	}
	if fsMeta.Size <= 0 {
		t.Errorf("expected size > 0, got %d", fsMeta.Size)
	}
	if fsMeta.Timestamp == "" {
		t.Error("expected timestamp to be set, got empty string")
	}
}

func createZipFile(zipPath, fileName, content string) error {
	out, err := os.Create(zipPath)
	if err != nil {
		return err
	}
	defer out.Close()

	zipWriter := zip.NewWriter(out)
	defer zipWriter.Close()

	writer, err := zipWriter.Create(fileName)
	if err != nil {
		return err
	}
	_, err = io.Copy(writer, bytes.NewBufferString(content))
	return err
}
