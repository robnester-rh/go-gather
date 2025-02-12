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

package detector

import (
	"github.com/conforma/go-gather/gather/file"
	"github.com/conforma/go-gather/gather/git"
	"github.com/conforma/go-gather/gather/http"
	"github.com/conforma/go-gather/gather/oci"
)

// FileDetector checks if the URI is a file path.
func FileDetector(uri string) bool {
	f := file.FileGatherer{}
	return f.Matcher(uri)
}

// GitDetector checks if the URI is a git repository.
func GitDetector(uri string) bool {
	g := git.GitGatherer{}
	return g.Matcher(uri)
}

// HttpDetector checks if the URI is an HTTP URL.
func HttpDetector(uri string) bool {
	h := http.HTTPGatherer{}
	return h.Matcher(uri)
}

// OciDetector checks if the URI is an OCI registry.
func OciDetector(uri string) bool {
	o := oci.OCIGatherer{}
	return o.Matcher(uri)
}
