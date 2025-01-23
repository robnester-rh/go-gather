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
	"os"
	"path/filepath"
	"testing"
)

// mockExpander is a trivial expander that returns true if its "keyword" is in the extension.
type mockExpander struct {
	keyword string
}

func (m *mockExpander) Expand(ctx context.Context, source, destination string, umask os.FileMode) error {
	// no-op
	return nil
}

func (m *mockExpander) Matcher(extension string) bool {
	return bytes.Contains([]byte(extension), []byte(m.keyword))
}

// TestRegisterAndGetExpander ensures we can register and retrieve an expander via GetExpander.
func TestRegisterAndGetExpander(t *testing.T) {
	oldExpanders := expanders
	expanders = nil
	defer func() { expanders = oldExpanders }()

	mockFoo := &mockExpander{keyword: "foo"}
	mockBar := &mockExpander{keyword: "bar"}

	RegisterExpander(mockFoo)
	RegisterExpander(mockBar)

	if got := GetExpander("my.foo"); got != mockFoo {
		t.Errorf("expected mockFoo, got %#v", got)
	}
	if got := GetExpander("some.bar"); got != mockBar {
		t.Errorf("expected mockBar, got %#v", got)
	}

	if got := GetExpander("nope.zip"); got != nil {
		t.Errorf("expected nil, got %#v", got)
	}
}

// TestIsCompressedFile checks that known magic numbers are correctly recognized.
func TestIsCompressedFile(t *testing.T) {
	tests := []struct {
		name           string
		magic          []byte
		wantCompressed bool
	}{
		{
			name:           "gzip magic",
			magic:          []byte{0x1f, 0x8b},
			wantCompressed: true,
		},
		{
			name:           "zip magic",
			magic:          []byte{0x50, 0x4b, 0x03, 0x04},
			wantCompressed: true,
		},
		{
			name:           "tar magic",
			magic:          []byte{0x75, 0x73, 0x74, 0x61, 0x72}, // "ustar"
			wantCompressed: true,
		},
		{
			name:           "bzip2 magic",
			magic:          []byte{0x42, 0x5a, 0x68}, // "BZh"
			wantCompressed: true,
		},
		{
			name:           "xz magic",
			magic:          []byte{0xfd, 0x37, 0x7a, 0x58, 0x5a, 0x00},
			wantCompressed: true,
		},
		{
			name:           "7z magic",
			magic:          []byte{0x37, 0x7a, 0xbc, 0xaf, 0x27, 0x1c},
			wantCompressed: true,
		},
		{
			name:           "unknown magic",
			magic:          []byte{0xde, 0xad, 0xbe, 0xef},
			wantCompressed: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tmpFile, err := os.CreateTemp("", "test-magic-*.bin")
			if err != nil {
				t.Fatalf("failed to create temp file: %v", err)
			}
			defer os.Remove(tmpFile.Name())
			defer tmpFile.Close()

			if _, err := tmpFile.Write(tc.magic); err != nil {
				t.Fatalf("failed to write magic bytes: %v", err)
			}
			if _, err := tmpFile.Seek(0, 0); err != nil {
				t.Fatalf("failed to seek temp file: %v", err)
			}

			got, err := IsCompressedFile(tmpFile.Name())
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tc.wantCompressed {
				t.Errorf("IsCompressedFile(%q) = %v, want %v", tc.name, got, tc.wantCompressed)
			}
		})
	}
}

// TestIsCompressedFile_NonExistent checks that an error is returned when the file does not exist.
func TestIsCompressedFile_NonExistent(t *testing.T) {
	fakePath := filepath.Join(os.TempDir(), "no_such_file_12345.bin")
	got, err := IsCompressedFile(fakePath)
	if err == nil {
		t.Fatalf("expected error when checking a non-existent file, but got nil")
	}
	if got {
		t.Errorf("expected false for isCompressed, got true")
	}
}

// TestIsCompressedFile_EmptyFile checks if an empty file doesn't match any known magic number.
func TestIsCompressedFile_EmptyFile(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "empty-*.bin")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())
	defer tmpFile.Close()

	got, err := IsCompressedFile(tmpFile.Name())
	if err != nil {
		t.Fatalf("unexpected error for empty file: %v", err)
	}
	if got {
		t.Errorf("expected false for empty file, got true")
	}
}
