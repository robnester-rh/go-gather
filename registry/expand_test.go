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

package registry_test

import (
	"testing"

	"github.com/conforma/go-gather/registry"
)

func TestGetExpanderTar(t *testing.T) {
	e := registry.GetExpander("tar")
	if e == nil {
		t.Fatal("expected non-nil tar expander")
	}
}

func TestGetExpanderZip(t *testing.T) {
	e := registry.GetExpander("zip")
	if e == nil {
		t.Fatal("expected non-nil zip expander")
	}
}

func TestGetExpanderBzip2(t *testing.T) {
	e := registry.GetExpander("bzip2")
	if e == nil {
		t.Fatal("expected non-nil bzip2 expander")
	}
}

func TestGetExpanderUnknown(t *testing.T) {
	e := registry.GetExpander("unknown")
	if e != nil {
		t.Fatal("expected nil for unknown expander")
	}
}
