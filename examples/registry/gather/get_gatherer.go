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
	"reflect"

	"github.com/conforma/go-gather/registry"
)

func main(){
	//-------------------------------------------------------------------------
	// The following code shows how to utilize the GetGatherer function from the
	// registry package to get the appropriate gatherer for a given source URL
	//-------------------------------------------------------------------------

	// Set the source URL to the file
	srcs := []string{"file:///tmp/test.bz2", "https://www.google.com/index.html", "git://github.com/conforma/go-gather.git", "oci::quay.io/konflux-ci/tekton-catalog/task-buildah-remote-oci-ta:0.3"}

	for _, src := range srcs {
		// Get the gatherer for the given source URL
		gatherer, err := registry.GetGatherer(src)
		if err != nil {
			panic(err)
		}

		// Check the type of the gatherer returned
		println("Gatherer Type: ", reflect.TypeOf(gatherer).String())
	}
}