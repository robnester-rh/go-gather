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
