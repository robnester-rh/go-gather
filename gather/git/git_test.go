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

package git

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
)

func TestGitGatherer_Matcher(t *testing.T) {
	gg := GitGatherer{}

	testCases := []struct {
		name string
		uri  string
		want bool
	}{
		{"git@ domain", "git@github.com:org/repo.git", true},
		{"git protocol double colon", "git::github.com/org/repo", true},
		{"git protocol slash slash", "git://github.com/org/repo.git", true},
		{"unknown protocol double colon", "s3::github.com/org/repo", false},
		{"dot git suffix", "https://github.com/org/repo.git", true},
		{"match github.com", "github.com/org/repo", true},
		{"not match githubusercontent.com", "https://raw.githubusercontent.com/foo/bar", false},
		{"other prefix", "svn://some/repo", false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := gg.Matcher(tc.uri)
			if got != tc.want {
				t.Errorf("Matcher(%q) = %v, want %v", tc.uri, got, tc.want)
			}
		})
	}
}

func TestGitGatherer_Gather_CanceledContext(t *testing.T) {
	gg := GitGatherer{}

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, err := gg.Gather(ctx, "git::github.com/org/repo", "/tmp/dest/dir")
	if err == nil {
		t.Fatal("expected context cancellation error, got nil")
	}
	if ctx.Err() == nil {
		t.Errorf("expected context to be canceled, but ctx.Err() is nil")
	}
}

func TestGitGatherer_Gather_InvalidRef(t *testing.T) {
	gg := GitGatherer{}
	sourceDir := t.TempDir()
	repoPath, _ := initLocalGitRepo(t, sourceDir)

	invalidRefURI := fmt.Sprintf("git::%s?ref=refs/heads/nope", repoPath)

	destDir := t.TempDir()
	ctx := context.Background()

	_, err := gg.Gather(ctx, invalidRefURI, destDir)
	if err == nil {
		t.Fatal("expected error for invalid ref, got nil")
	}
	if !strings.Contains(err.Error(), "error cloning repository: reference not found") {
		t.Errorf("expected 'error cloning repository: reference not found' error, got %v", err)
	}
}

func initLocalGitRepo(t *testing.T, repoDir string) (string, string) {
	t.Helper()

	repo, err := git.PlainInit(repoDir, false)
	if err != nil {
		t.Fatalf("failed to init local git repo in %s: %v", repoDir, err)
	}

	readmePath := filepath.Join(repoDir, "README.md")
	content := []byte("# Test Repo\n")
	if err := os.WriteFile(readmePath, content, 0600); err != nil {
		t.Fatalf("failed to write file in local repo: %v", err)
	}

	w, err := repo.Worktree()
	if err != nil {
		t.Fatalf("failed to get worktree: %v", err)
	}
	if _, err = w.Add("README.md"); err != nil {
		t.Fatalf("failed to add README.md to index: %v", err)
	}
	commit, err := w.Commit("Initial commit", &git.CommitOptions{
		Author: &object.Signature{
			Name:  "Tester",
			Email: "tester@example.com",
			When:  time.Now(),
		},
	})
	if err != nil {
		t.Fatalf("failed to commit: %v", err)
	}

	return repoDir, commit.String()
}

func TestProcessUrl(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		wantRef   string
		wantSub   string
		wantDepth string
		wantErr   bool
	}{
		{
			name:  "https github URL",
			input: "https://github.com/org/repo",
		},
		{
			name:    "https with ref",
			input:   "https://github.com/org/repo?ref=v1.0",
			wantRef: "v1.0",
		},
		{
			name:    "https with ref and subdir",
			input:   "https://github.com/org/repo?ref=main//subdir",
			wantRef: "main",
			wantSub: "subdir",
		},
		{
			name:      "https with depth",
			input:     "https://github.com/org/repo?depth=1",
			wantDepth: "1",
		},
		{
			name:    "git:: prefix with ref",
			input:   "git::https://github.com/org/repo?ref=abc123",
			wantRef: "abc123",
		},
		{
			name:    "git@ SSH URL",
			input:   "git@github.com:org/repo",
			wantSub: "github.com",
		},
		{
			name:    "path with subdir via double slash",
			input:   "https://github.com/org/repo//policies/base",
			wantSub: "policies/base",
		},
		{
			name:  "file path",
			input: "file:///tmp/local-repo",
		},
		{
			name:  "relative file path",
			input: "./local-repo",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			src, ref, subdir, depth, err := processUrl(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("processUrl(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
				return
			}
			if err != nil {
				return
			}
			if src == "" {
				t.Errorf("processUrl(%q) returned empty src", tt.input)
			}
			if ref != tt.wantRef {
				t.Errorf("processUrl(%q) ref = %q, want %q", tt.input, ref, tt.wantRef)
			}
			if subdir != tt.wantSub {
				t.Errorf("processUrl(%q) subdir = %q, want %q", tt.input, subdir, tt.wantSub)
			}
			if depth != tt.wantDepth {
				t.Errorf("processUrl(%q) depth = %q, want %q", tt.input, depth, tt.wantDepth)
			}
		})
	}
}
