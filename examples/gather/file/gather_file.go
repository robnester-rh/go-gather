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

	"github.com/conforma/go-gather/gather/file"
)


func main(){

	//-------------------------------------------------------------------------
	// These lines of code create a source directory containing 5 files, and a
	// destination directory to gather the contents of the source directory to.
	//-------------------------------------------------------------------------
	src, dst := setup("file_gather_example")
	defer os.RemoveAll(src)
	defer os.RemoveAll(dst)


	//-------------------------------------------------------------------------
	// The following code shows how to gather the contents of a source directory
	// to a destination directory using the file gatherer
	//-------------------------------------------------------------------------

	// Create a new file gatherer
	g := &file.FileGatherer{}

	// Gather the contents of the source directory to the destination directory, 
	// returning a metadata struct
	m, err := g.Gather(context.Background(), src, dst)
	if err != nil {
		panic(err)
	}
	// Do a type assertion for ease of use
	metadata := m.(*file.FSMetadata)

	// Print the metadata
	fmt.Println("Destination Path: ", metadata.Path)
	fmt.Println("Directory Size: ", metadata.Size)
}


func setup(base string) (string, string) {
	// Setting up a temporary directory to act as our source dir
	src, err := os.MkdirTemp("", fmt.Sprintf("%s_src_", base))
	if err != nil {
		panic(err)
	}

	// Setting up a temporary directory to act as our destination dir
	dst, err := os.MkdirTemp("", fmt.Sprintf("%s_dst_", base))
	if err != nil {
		panic(err)
	}
	// Create five files in the src directory, each containing the string "Hi there"
	err = createFiles(src, base, 5)
	if err != nil {
		panic(err)
	}
	return src, dst
}

func createFiles(src string, fileBase string, fileCount int) error {
	for i := 0; i < fileCount; i++ {			
		f, err := os.CreateTemp(src, fmt.Sprintf("%s_file_", fileBase))
		if err != nil {
			return err
		}
		defer f.Close()
		f.Write([]byte(fmt.Sprintf("Hi there %d\n", i)))
	}
	return nil
}