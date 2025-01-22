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
	"fmt"
	"os"
	"path/filepath"

	"github.com/enterprise-contract/go-gather/expand/bzip2"
)

func main() {

	//-------------------------------------------------------------------------
	// The following code sets up a source directory containing a bzip2 compressed
	// file, and a destination directory to expand the bzip2 compressed file to.
	//-------------------------------------------------------------------------

	// Set the source URL to the bzip2 compressed file
	src, dst := setup()
	defer os.RemoveAll(src)
	defer os.RemoveAll(dst)

	//-------------------------------------------------------------------------
	// The following code shows how to utilize a bzip2 expander to expand a bzip2
	// compressed file to a destination directory
	//-------------------------------------------------------------------------

	// Create a new bzip2 expander
	b := &bzip2.Bzip2Expander{}

	// Expand the bzip2 compressed file to the destination directory
	err := b.Expand(context.Background(), filepath.Join(src, "test.bz2"), dst, 0600)
	if err != nil {
		panic(err)
	}

	// Print the destination directory
	println("Destination Directory: ", dst)

	// Get the contents of the destination directory
	contents, err := os.ReadDir(dst)
	if err != nil {
		panic(err)
	}

	// Print the contents of the destination directory
	for _, entry := range contents {
		println("File Name: ", entry.Name())
	}
}

func setup() (string, string) {
	// Setting up a temporary directory to act as our source dir
	src, err := os.MkdirTemp("", "bzip2_expander_example_src_")
	if err != nil {
		panic(err)
	}

	// Setting up a temporary directory to act as our destination dir
	dst, err := os.MkdirTemp("", "bzip2_expander_example_dst_")
	if err != nil {
		panic(err)
	}

	// Create a bzip2 compressed file in the src directory
	err = createBzip2File(filepath.Join(src, "test.bz2"))
	if err != nil {
		panic(err)
	}

	return src, dst
}

func createBzip2File(filepath string) error {
	fmt.Println("Creating bzip2 file at", filepath)
	err := os.WriteFile(filepath, helloBzip2Fixture, 0600)
	if err != nil {
		return err
	}
	return nil
}

// helloBzip2Fixture is a small bzip2-encoded byte slice that decompresses to "Hello Bzip2!".
// This was generated externally to ensure its validity.
var helloBzip2Fixture = []byte{
	0x42, 0x5a, 0x68, 0x39, 0x31, 0x41, 0x59, 0x26, 0x53, 0x59, 0x8e, 0x9d,
	0x35, 0x69, 0x00, 0x00, 0x02, 0x1d, 0x80, 0x60, 0x00, 0x10, 0x00, 0x10,
	0x40, 0x02, 0x24, 0xc0, 0x10, 0x20, 0x00, 0x31, 0x00, 0xd3, 0x4d, 0x04,
	0x0d, 0x06, 0x9a, 0x11, 0xc2, 0xb1, 0x14, 0xc9, 0x78, 0xbb, 0x92, 0x29,
	0xc2, 0x84, 0x84, 0x74, 0xe9, 0xab, 0x48,
}
