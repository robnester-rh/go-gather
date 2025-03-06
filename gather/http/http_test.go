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

package http

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestHTTPGatherer_Matcher(t *testing.T) {
	g := &HTTPGatherer{}

	testCases := []struct {
		name string
		uri  string
		want bool
	}{
		{"http scheme", "http://example.com/file.txt", true},
		{"https scheme", "https://example.com/file.txt", true},
		{"no scheme", "example.com/file.txt", false},
		{"ftp scheme", "ftp://example.com/file.txt", false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := g.Matcher(tc.uri)
			if got != tc.want {
				t.Errorf("Matcher(%q) = %v, want %v", tc.uri, got, tc.want)
			}
		})
	}
}

func TestHTTPGatherer_Gather_Success(t *testing.T) {
	testData := "Hello from test server!"
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(testData))
	})
	server := httptest.NewServer(handler)
	defer server.Close()

	g := NewHTTPGatherer()

	tempDir := t.TempDir()
	dest := filepath.Join(tempDir, "downloaded_file.txt")

	ctx := context.Background()
	meta, err := g.Gather(ctx, server.URL+"/subdir/file.txt", dest)
	if err != nil {
		t.Fatalf("Gather returned unexpected error: %v", err)
	}

	fileContent, err := os.ReadFile(dest)
	if err != nil {
		t.Fatalf("failed to read downloaded file: %v", err)
	}
	if string(fileContent) != testData {
		t.Errorf("expected content %q, got %q", testData, string(fileContent))
	}

	httpMeta, ok := meta.(HTTPMetadata)
	if !ok {
		t.Fatalf("expected *HTTPMetadata, got %T", meta)
	}
	if httpMeta.URI != server.URL+"/subdir/file.txt" {
		t.Errorf("expected URI=%s, got %s", server.URL+"/subdir/file.txt", httpMeta.URI)
	}
	if httpMeta.Path != dest {
		t.Errorf("expected Path=%s, got %s", dest, httpMeta.Path)
	}
	if httpMeta.ResponseCode != http.StatusOK {
		t.Errorf("expected 200, got %d", httpMeta.ResponseCode)
	}
	if httpMeta.Size != int64(len(testData)) {
		t.Errorf("expected size=%d, got %d", len(testData), httpMeta.Size)
	}
	if httpMeta.Timestamp == "" {
		t.Error("expected non-empty timestamp")
	}
}

func TestHTTPGatherer_Gather_NoScheme(t *testing.T) {
	g := NewHTTPGatherer()
	ctx := context.Background()

	tempDir := t.TempDir()
	dest := filepath.Join(tempDir, "file.txt")

	_, err := g.Gather(ctx, "example.com/file.txt", dest)
	if err == nil {
		t.Fatal("expected an error when no scheme is provided, got nil")
	}
	if !strings.Contains(err.Error(), "no source scheme provided") {
		t.Errorf("expected error mentioning missing scheme, got %v", err)
	}
}

func TestHTTPGatherer_Gather_NoPath(t *testing.T) {
	g := NewHTTPGatherer()
	ctx := context.Background()

	tempDir := t.TempDir()
	dest := filepath.Join(tempDir, "file.txt")

	// Provide a URL with a scheme but no path
	_, err := g.Gather(ctx, "http://example.com", dest)
	if err == nil {
		t.Fatal("expected error when URL has no path, got nil")
	}
	if !strings.Contains(err.Error(), "specify a path to a file to download") {
		t.Errorf("expected error about specifying a path, got %v", err)
	}
}

func TestHTTPGatherer_Gather_Non200(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.NotFound(w, r)
	})
	server := httptest.NewServer(handler)
	defer server.Close()

	g := NewHTTPGatherer()
	tempDir := t.TempDir()
	dest := filepath.Join(tempDir, "file.txt")

	ctx := context.Background()
	_, err := g.Gather(ctx, server.URL+"/missing-file.txt", dest)
	if err == nil {
		t.Fatal("expected an error for non-200 response, got nil")
	}
	if !strings.Contains(err.Error(), "received non-200 response code") {
		t.Errorf("expected error about non-200 response, got %v", err)
	}
}

func TestHTTPGatherer_Gather_EmptyDirDestination(t *testing.T) {
	testData := "Test data"
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(testData))
	})
	server := httptest.NewServer(handler)
	defer server.Close()

	g := NewHTTPGatherer()

	tempDir := t.TempDir()

	dest := filepath.Join(tempDir, "someDir") + "/"

	ctx := context.Background()
	srcURL := server.URL + "/download-me.bin"
	meta, err := g.Gather(ctx, srcURL, dest)
	if err != nil {
		t.Fatalf("Gather returned error: %v", err)
	}

	httpMeta := meta.(HTTPMetadata)
	expectedPath := filepath.Join(dest, "download-me.bin")
	if httpMeta.Path != expectedPath {
		t.Errorf("expected path=%s, got %s", expectedPath, httpMeta.Path)
	}

	fileContent, err := os.ReadFile(expectedPath)
	if err != nil {
		t.Fatalf("failed to read downloaded file: %v", err)
	}
	if string(fileContent) != testData {
		t.Errorf("expected content=%q, got %q", testData, string(fileContent))
	}
}

func TestHTTPGatherer_Gather_CanceledContext(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Simulate a slow response so we can cancel the context
		time.Sleep(2 * time.Second)
		_, _ = w.Write([]byte("Large amount of data..."))
	})
	server := httptest.NewServer(handler)
	defer server.Close()

	g := NewHTTPGatherer()
	tempDir := t.TempDir()
	dest := filepath.Join(tempDir, "file.txt")

	// Create a context and cancel it immediately
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, err := g.Gather(ctx, server.URL+"/slow-file", dest)
	if err == nil {
		t.Fatal("expected an error due to context cancellation, got nil")
	}
	if ctx.Err() == nil {
		t.Errorf("expected context to be canceled, got nil")
	}
}
