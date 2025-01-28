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

package expand

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
)

/* package expander provides an interface for expanders to implement. Expanders are used to expand compressed files. */

type Expander interface {
	Expand(ctx context.Context, source string, destination string, umask os.FileMode) error
	Matcher(extension string) bool
}

var expanders []Expander

type ExpandOptions struct{}

func GetExpander(extension string) Expander {
	for _, expander := range expanders {
		if expander.Matcher(extension) {
			return expander
		}
	}
	return nil
}

func RegisterExpander(e Expander) {
	expanders = append(expanders, e)
}

// Known magic numbers for common compressed file formats
var magicNumbers = map[string][]byte{
	"gzip":  {0x1f, 0x8b},
	"zip":   {0x50, 0x4b, 0x03, 0x04},
	"bzip2": {0x42, 0x5a, 0x68},
	"xz":    {0xfd, 0x37, 0x7a, 0x58, 0x5a, 0x00},
	"7z":    {0x37, 0x7a, 0xbc, 0xaf, 0x27, 0x1c},
}

func IsCompressedFile(filePath string) (bool, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return false, fmt.Errorf("could not open file: %w", err)
	}
	defer file.Close()

	// Read the first few bytes
	header := make([]byte, 10) // maximum length of magic numbers
	_, err = file.Read(header)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) || errors.Is(err, os.ErrPermission) {
			return false, fmt.Errorf("could not read file header: %w", err)
		}
		return false, nil
	}

	// Check against known magic numbers
	for _, magic := range magicNumbers {
		if len(header) >= len(magic) && bytes.Equal(header[:len(magic)], magic) {
			return true, nil
		}
	}

	return false, nil
}

// IsTarFile checks whether the file at filePath is a tar archive by reading
// the standard tar magic bytes at offset 257 ("ustar\0" or "ustar ").
func IsTarFile(filePath string) (bool, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return false, fmt.Errorf("could not open file: %w", err)
	}
	defer f.Close()

	// Move to the position where the tar magic string should appear.
	_, err = f.Seek(257, io.SeekStart)
	if err != nil {
		return false, fmt.Errorf("could not seek in file: %w", err)
	}

	// Read the 6 bytes of the magic string ("ustar\0" or "ustar ").
	magic := make([]byte, 6)
	n, err := f.Read(magic)
	if err != nil && err != io.EOF {
		return false, fmt.Errorf("could not read magic bytes: %w", err)
	}

	// If we didn't get enough bytes, it can't be a valid tar.
	if n < 6 {
		return false, nil
	}

	// Check if we have "ustar" at the start (POSIX tar magic).
	if bytes.HasPrefix(magic, []byte("ustar")) {
		return true, nil
	}

	return false, nil
}