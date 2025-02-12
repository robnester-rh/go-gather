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

	"github.com/conforma/go-gather/gather/oci"
)

func main() {
	// -------------------------------------------------------------------------
	// The following code shows how to gather the contents of an OCI registry to a
	// destination directory using the oci gatherer
	// -------------------------------------------------------------------------

	// Set the source URL to the OCI registry
	src := "oci::quay.io/konflux-ci/tekton-catalog/task-buildah-remote-oci-ta:0.3"

	// create a temporary directory to act as our destination directory
	dst, err := os.MkdirTemp("", "oci_gather_example_dst_")
	if err != nil {
		panic(err)
	}
	defer os.RemoveAll(dst)

	// Create a new oci gatherer
	o := &oci.OCIGatherer{}

	// Gather the contents of the OCI registry to the destination directory,
	m, err := o.Gather(context.Background(), src, dst)
	if err != nil {
		panic(err)
	}

	// Do a type assertion for ease of use
	metadata := m.(*oci.OCIMetadata)

	// Print the metadata
	println("Destination Path: ", metadata.Path)
	println("Diget", metadata.Digest)
	println("Timestamp", metadata.Timestamp)
	
}
