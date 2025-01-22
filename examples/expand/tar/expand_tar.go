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
	"archive/tar"
	"context"
	"os"
	"path/filepath"

	tarExpander "github.com/enterprise-contract/go-gather/expand/tar"
)

func main() {
	//-------------------------------------------------------------------------
	// The following code sets up a source directory containing a tar archive
	// file, "test.tar" containing a file, "test.txt", and a destination 
	// directory to expand the tar compressed file to.
	//-------------------------------------------------------------------------
	src, dst := setup()
	defer os.RemoveAll(src)
	defer os.RemoveAll(dst)

	//-------------------------------------------------------------------------
	// The following code shows how to utilize a tar expander to expand a tar
	// archive file to a destination directory
	//-------------------------------------------------------------------------

	// Set the source to the tar compressed file
	src = filepath.Join(src, "test.tar")


	// Create a new tar expander
	t := &tarExpander.TarExpander{}

	// Expand the gzip compressed file to the destination directory

	err := t.Expand(context.Background(), src, dst, 0600)
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
	src, err := os.MkdirTemp("", "tar_expander_example_src_")
	if err != nil {
		panic(err)
	}

	// Setting up a temporary directory to act as our destination dir
	dst, err := os.MkdirTemp("", "tar_expander_example_dst_")
	if err != nil {
		panic(err)
	}

	// Create a tar archive file in the src directory
	err = createTar(src, "test.tar")
	if err != nil {
		panic(err)
	}

	return src, dst
}

func createTar(src, name string) error {
	// Create a new tar archive file
	f, err := os.Create(filepath.Join(src, name))
	if err != nil {
		return err
	}
	defer f.Close()

	// Create a new tar writer
	tw := tar.NewWriter(f)
	defer tw.Close()

	// Create a new file header
	hdr := &tar.Header{
		Name: "test.txt",
		Mode: 0600,
		Size: int64(len("Hello, World!")),
	}

	// Write the file header
	if err := tw.WriteHeader(hdr); err != nil {
		return err
	}

	// Write the file contents
	if _, err := tw.Write([]byte("Hello, World!")); err != nil {
		return err
	}

	return nil
}