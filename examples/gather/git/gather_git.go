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

package main

import (
	"context"
	"os"

	"github.com/enterprise-contract/go-gather/gather/git"
)

func main() {
	// -------------------------------------------------------------------------
	// The following code shows how to gather the contents of a git repository
	// to a destination directory using the git gatherer
	// -------------------------------------------------------------------------

	// Set the source URL to the git repository
	gitSrc := "git://github.com/enterprise-contract/go-gather.git"

	// Create a temporary directory to act as our destination directory
	dst, err := os.MkdirTemp("", "git_gather_example_dst_")
	if err != nil {
		panic(err)
	}
	defer os.RemoveAll(dst)

	// Create a new git gatherer
	g := &git.GitGatherer{}

	// Gather the contents of the git repository to the destination directory,
	m, err := g.Gather(context.Background(), gitSrc, dst)
	if err != nil {
		panic(err)
	}

	// Do a type assertion for ease of use
	metadata := m.(*git.GitMetadata)

	// Print the metadata
	println("Destination Path: ", metadata.Path)
	println("Author", metadata.Author)
	println("Latest Commit", metadata.LatestCommit)
	println("Timestamp", metadata.Timestamp)
}