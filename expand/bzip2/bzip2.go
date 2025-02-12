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

package bzip2

import (
	"compress/bzip2"
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/conforma/go-gather/expand"
	"github.com/conforma/go-gather/internal/helpers"
)

var pathExpanderFunc = helpers.ExpandPath

type Bzip2Expander struct {
	FileSizeLimit int64
}

func (b *Bzip2Expander) Expand(ctx context.Context, src, dst string, umask os.FileMode) error {
	src, err := pathExpanderFunc(src)
	if err != nil {
		return fmt.Errorf("failed to expand source path: %w", err)
	}
	dst, err = pathExpanderFunc(dst)
	if err != nil {
		return fmt.Errorf("failed to expand destination path: %w", err)
	}

	input, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("failed to open bzip2 file %q: %w", src, err)
	}
	defer input.Close()

	bzipReader := bzip2.NewReader(input)

	// Ensure the parent directory of dst exists
	if err := os.MkdirAll(dst, umask); err != nil {
		return fmt.Errorf("failed to create directory %q: %w", filepath.Dir(dst), err)
	}

	baseName := strings.TrimSuffix(filepath.Base(src), filepath.Ext(src))

	fpath := filepath.Join(dst, baseName)
	// Create or truncate the output file
	outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return fmt.Errorf("failed to create file %q: %w", dst, err)
	}
	defer outFile.Close()

	const bufferSize = 32 * 1024 // 32 KB
	buffer := make([]byte, bufferSize)

	// Track total decompressed size to avoid decompression bombs.
	var totalBytes int64
	for {
		n, err := bzipReader.Read(buffer)
		if n > 0 {
			if totalBytes+int64(n) > b.FileSizeLimit && b.FileSizeLimit > 0 {
				return fmt.Errorf("decompressed file exceeds size limit of %d bytes", b.FileSizeLimit)
			}
			if _, writeErr := outFile.Write(buffer[:n]); writeErr != nil {
				return fmt.Errorf("failed to write decompressed data: %w", writeErr)
			}
			totalBytes += int64(n)
		}
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("error during decompression: %w", err)
		}
	}

	return nil
}

// Matcher checks if the extension matches supported formats.
func (b *Bzip2Expander) Matcher(extension string) bool {
	return (strings.Contains(extension, "bz2") || strings.Contains(extension, "bzip2")) && !strings.Contains(extension, "tar")
}

func init() {
	expand.RegisterExpander(&Bzip2Expander{})
}
