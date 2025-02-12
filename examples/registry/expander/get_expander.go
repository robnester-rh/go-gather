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

func main() {

	//-------------------------------------------------------------------------
	// The following code shows how to utilize the GetExpander function from the
	// registry package to get the appropriate expander for a given source URL
	//-------------------------------------------------------------------------

	// Set the source URL to the file
	src := "file:///tmp/test.bz2"

	// Get the expander for the given source URL
	expander := registry.GetExpander(src)
	if expander == nil {
		panic("No expander found for the given source URL")
	}

	// Check the type of the expander
	println("Expander Type: ", reflect.TypeOf(expander).String())

}
