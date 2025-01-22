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

	"github.com/enterprise-contract/go-gather/gather/http"
)

func main(){
	// -------------------------------------------------------------------------
	// The following code shows how to gather the contents of a http repository
	// to a destination directory using the http gatherer
	// -------------------------------------------------------------------------

	// Set the source URL to the http repository
	src := "https://www.google.com/index.html"

	// create a temporary directory to act as our destination directory
	dst, err := os.MkdirTemp("", "http_gather_example_dst_")
	if err != nil {
		panic(err)
	}
	defer os.RemoveAll(dst)

	h := &http.HTTPGatherer{}

	// Gather the contents of the http repository to the destination directory,
	m, err := h.Gather(context.Background(), src, dst)
	if err != nil {
		panic(err)
	}

	// Do a type assertion for ease of use
	metadata := m.(*http.HTTPMetadata)

	// Print the metadata
	println("Destination Path: ", metadata.Path)
	println("HTTP Response", metadata.ResponseCode)
	println("URI", metadata.URI)
	println("Size", metadata.Size)
	println("Timestamp", metadata.Timestamp)

}