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

func TestGetGathererForFile(t *testing.T) {
	g, err := registry.GetGatherer("file:///tmp/test")
	if err != nil {
		t.Fatalf("expected gatherer for file URI, got error: %v", err)
	}
	if g == nil {
		t.Fatal("expected non-nil gatherer")
	}
}

func TestGetGathererForGit(t *testing.T) {
	g, err := registry.GetGatherer("git::https://github.com/example/repo")
	if err != nil {
		t.Fatalf("expected gatherer for git URI, got error: %v", err)
	}
	if g == nil {
		t.Fatal("expected non-nil gatherer")
	}
}

func TestGetGathererForHTTP(t *testing.T) {
	g, err := registry.GetGatherer("https://example.com/file.txt")
	if err != nil {
		t.Fatalf("expected gatherer for HTTP URI, got error: %v", err)
	}
	if g == nil {
		t.Fatal("expected non-nil gatherer")
	}
}

func TestGetGathererForOCI(t *testing.T) {
	g, err := registry.GetGatherer("oci://quay.io/example/image:latest")
	if err != nil {
		t.Fatalf("expected gatherer for OCI URI, got error: %v", err)
	}
	if g == nil {
		t.Fatal("expected non-nil gatherer")
	}
}

func TestGetGathererUnknownURI(t *testing.T) {
	_, err := registry.GetGatherer("unknown://something")
	if err == nil {
		t.Fatal("expected error for unknown URI scheme")
	}
}
