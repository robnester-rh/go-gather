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
	"archive/zip"
	"context"
	"os"
	"path/filepath"

	zipExpander "github.com/enterprise-contract/go-gather/expand/zip"
)

func main() {
	//-------------------------------------------------------------------------
	// The following code creates a source directory containing a zip file, and a
	// destination directory to expand the zip file to.
	//-------------------------------------------------------------------------
	src, dst := setup()
	defer os.RemoveAll(src)
	defer os.RemoveAll(dst)

	//-------------------------------------------------------------------------
	// The following code shows how to utilize a zip expander to expand a zip
	// compressed file to a destination directory
	//-------------------------------------------------------------------------

	// Set the source to the zip compressed file
	src = filepath.Join(src, "test.zip")

	// Create a new zip expander
	z := &zipExpander.ZipExpander{}

	// Expand the zip compressed file to the destination directory
	err := z.Expand(context.Background(), src, dst, 0600)
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
	src, err := os.MkdirTemp("", "zip_expander_example_src_")
	if err != nil {
		panic(err)
	}

	// Setting up a temporary directory to act as our destination dir
	dst, err := os.MkdirTemp("", "zip_expander_example_dst_")
	if err != nil {
		panic(err)
	}

	// Create a zip compressed file in the src directory
	err = createZipFile(filepath.Join(src, "test.zip"))
	if err != nil {
		panic(err)
	}

	return src, dst
}

func createZipFile(filepath string) error {
	// Create a zip file in the src directory
	f, err := os.Create(filepath)
	if err != nil {
		return err
	}

	// Write a file to the zip file
	z := zip.NewWriter(f)
	defer z.Close()

	// Create a new file in the zip file
	fw, err := z.Create("test.txt")
	if err != nil {
		return err
	}

	// Write to the file in the zip file
	_, err = fw.Write([]byte("Hi there"))
	if err != nil {
		return err
	}

	return nil
}