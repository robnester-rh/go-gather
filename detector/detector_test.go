// Copyright The Enterprise Contract Contributors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
// SPDX-License-Identifier: Apache-2.0
package detector

import "testing"

// TestFileDetector tests the FileDetector function.
func TestFileDetector(t *testing.T) {
	testCases := []struct {
		name     string
		uri      string
		expected bool
	}{
		{
			name:     "valid file URI with file:: scheme",
			uri:      "file::/tmp/somefile.txt",
			expected: true,
		},
		{
			name:     "valid file URI with file:// scheme",
			uri:      "file:///tmp/somefile.txt",
			expected: true,
		},
		{
			name:     "valid file path (no scheme)",
			uri:      "/tmp/somefile.txt",
			expected: true,
		},
		{
			name:     "git URI should not match as file",
			uri:      "git@github.com:enterprise-contract/go-gather.git",
			expected: false,
		},
		{
			name:     "oci URI should not match as file",
			uri:      "oci://registry/repository:tag",
			expected: false,
		},
		{
			name:     "empty string",
			uri:      "",
			expected: false,
		},
	}

	for _, tc := range testCases {
		tc := tc // capture range variable
		t.Run(tc.name, func(t *testing.T) {
			got := FileDetector(tc.uri)
			if got != tc.expected {
				t.Errorf("FileDetector(%q) = %v, want %v", tc.uri, got, tc.expected)
			}
		})
	}
}

// TestGitDetector tests the GitDetector function.
func TestGitDetector(t *testing.T) {
	testCases := []struct {
		name     string
		uri      string
		expected bool
	}{
		{
			name:     "valid git SSH URI",
			uri:      "git@github.com:enterprise-contract/go-gather.git",
			expected: true,
		},
		{
			name:     "valid git HTTPS URI",
			uri:      "https://github.com/enterprise-contract/go-gather.git",
			expected: true,
		},
		{
			name:     "valid git:// URI",
			uri:      "git://example.com/enterprise-contract/go-gather",
			expected: true,
		},
		{
			name:     "valid git:: URI",
			uri:      "git::example.com/enterprise-contract/go-gather",
			expected: true,
		},
		{
			name:     "file URI should not match as git",
			uri:      "file:///tmp/somefile.txt",
			expected: false,
		},
		{
			name:     "file path with git scheme should match as git",
			uri:      "git:///tmp/somefile",
			expected: true,
		},
		{
			name:     "oci URI should not match as git",
			uri:      "oci://registry/repository:tag",
			expected: false,
		},
		{
			name:     "empty string",
			uri:      "",
			expected: false,
		},
	}

	for _, tc := range testCases {
		tc := tc // capture range variable
		t.Run(tc.name, func(t *testing.T) {
			got := GitDetector(tc.uri)
			if got != tc.expected {
				t.Errorf("GitDetector(%q) = %v, want %v", tc.uri, got, tc.expected)
			}
		})
	}
}

// TestHttpDetector checks if HttpDetector correctly identifies HTTP URLs.
func TestHttpDetector(t *testing.T) {
	testCases := []struct {
		name     string
		uri      string
		expected bool
	}{
		{
			name:     "valid HTTP URL",
			uri:      "http://example.com/somepath",
			expected: true,
		},
		{
			name:     "valid HTTPS URL",
			uri:      "https://example.com/somepath",
			expected: true,
		},
		{
			name:     "known web git repos should not match as HTTP",
			uri:      "https://github.com/somepath",
			expected: false,
		},
		{
			name:     "file URI should not match as HTTP",
			uri:      "file:///tmp/somefile.txt",
			expected: false,
		},
		{
			name:     "git URI should not match as HTTP",
			uri:      "git@github.com:enterprise-contract/go-gather.git",
			expected: false,
		},
		{
			name:     "OCI URI should not match as HTTP",
			uri:      "oci://registry/repository:tag",
			expected: false,
		},
		{
			name:     "empty string",
			uri:      "",
			expected: false,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			got := HttpDetector(tc.uri)
			if got != tc.expected {
				t.Errorf("HttpDetector(%q) = %v, want %v", tc.uri, got, tc.expected)
			}
		})
	}
}

// TestOciDetector tests the OciDetector function.
func TestOciDetector(t *testing.T) {
	testCases := []struct {
		name     string
		uri      string
		expected bool
	}{
		{
			name:     "valid OCI URI",
			uri:      "oci://registry/repository:tag",
			expected: true,
		},
		{
			name:     "valid OCI URI",
			uri:      "oci::registry/repository:tag",
			expected: true,
		},
		{
			name:     "valid OCI URI",
			uri:      "localhost:32433/repository:tag",
			expected: true,
		},
		{
			name:     "file URI should not match as OCI",
			uri:      "file:///tmp/somefile.txt",
			expected: false,
		},
		{
			name:     "git URI should not match as OCI",
			uri:      "git@github.com:enterprise-contract/go-gather.git",
			expected: false,
		},
		{
			name:     "empty string",
			uri:      "",
			expected: false,
		},
	}

	for _, tc := range testCases {
		tc := tc // capture range variable
		t.Run(tc.name, func(t *testing.T) {
			got := OciDetector(tc.uri)
			if got != tc.expected {
				t.Errorf("OciDetector(%q) = %v, want %v", tc.uri, got, tc.expected)
			}
		})
	}
}
