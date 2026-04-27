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

// Package oci implements a Gatherer for OCI registry sources.
package oci

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"
	"time"

	"oras.land/oras-go/v2"
	"oras.land/oras-go/v2/content/file"
	"oras.land/oras-go/v2/registry"
	"oras.land/oras-go/v2/registry/remote"
	"oras.land/oras-go/v2/registry/remote/retry"

	"github.com/conforma/go-gather/gather"
	r "github.com/conforma/go-gather/internal/oci/registry"
	"github.com/conforma/go-gather/metadata"
)

// OCIGatherer gathers artifacts from OCI-compliant registries.
type OCIGatherer struct {
	// transport is the http.RoundTripper for registry requests. If nil, http.DefaultTransport is used.
	transport http.RoundTripper
}

// OCIMetadata holds metadata about a gathered OCI artifact.
type OCIMetadata struct {
	Path      string
	Digest    string
	Timestamp string
}

// Option configures an OCIGatherer.
type Option func(*OCIGatherer)

// WithTransport sets the http.RoundTripper used for OCI registry requests.
// Callers are responsible for wrapping the transport with retry if desired.
func WithTransport(t http.RoundTripper) Option {
	return func(g *OCIGatherer) { g.transport = t }
}

// NewOCIGatherer creates an OCIGatherer with the given options.
func NewOCIGatherer(opts ...Option) *OCIGatherer {
	g := &OCIGatherer{}
	for _, opt := range opts {
		opt(g)
	}
	return g
}

var orasCopy = oras.Copy

// localhostHostRegexp matches "localhost" only when it appears as a hostname (after a scheme prefix).
var localhostHostRegexp = regexp.MustCompile(`(^|://|::)(localhost)([:/?#]|$)`)

var ociRegistryPatterns = []*regexp.Regexp{
	regexp.MustCompile(`(^|\.)azurecr\.io$`),
	regexp.MustCompile(`(^|\.)gcr\.io$`),
	regexp.MustCompile(`^registry\.gitlab\.com$`),
	regexp.MustCompile(`(^|\.)pkg\.dev$`),
	regexp.MustCompile(`^[0-9]{12}\.dkr\.ecr\.[a-z0-9-]+\.amazonaws\.com$`),
	regexp.MustCompile(`^quay\.io$`),
	regexp.MustCompile(`^(?:::1|127\.0\.0\.1|(?i:localhost))$`),
}

// Gather pulls an OCI artifact from source into the dst directory.
func (o *OCIGatherer) Gather(ctx context.Context, source, dst string) (meta metadata.Metadata, err error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	transport := o.transport
	if transport == nil {
		transport = retry.NewTransport(http.DefaultTransport)
	}

	if localhostHostRegexp.MatchString(source) {
		source = localhostHostRegexp.ReplaceAllString(source, "${1}127.0.0.1${3}")
	}

	// Parse the source URI
	repo := ociURLParse(source)

	// Get the artifact reference
	ref, err := registry.ParseReference(repo)
	if err != nil {
		return nil, fmt.Errorf("failed to parse reference: %w", err)
	}

	// If the reference is empty, set it to "latest"
	if ref.Reference == "" {
		ref.Reference = "latest"
		repo = ref.String()
	}

	// Create the repository client
	src, err := remote.NewRepository(repo)
	if err != nil {
		return nil, fmt.Errorf("failed to create repository client: %w", err)
	}

	// Setup the client for the repository
	if err := r.SetupClient(src, transport); err != nil {
		return nil, fmt.Errorf("failed to setup repository client: %w", err)
	}

	// Create the destination directory
	if err := os.MkdirAll(dst, 0755); err != nil {
		return nil, fmt.Errorf("failed to create directory: %w", err)
	}

	// Create the file store
	fileStore, err := file.New(dst)
	if err != nil {
		return nil, fmt.Errorf("file store: %w", err)
	}
	defer func() {
		if cerr := fileStore.Close(); cerr != nil && err == nil {
			meta = nil
			err = fmt.Errorf("failed to close OCI file store: %w", cerr)
		}
	}()

	// Copy the artifact to the file store
	a, err := orasCopy(ctx, src, repo, fileStore, "", oras.DefaultCopyOptions)
	if err != nil {
		return nil, fmt.Errorf("pulling policy: %w", err)
	}

	return &OCIMetadata{
		Digest:    a.Digest.String(),
		Path:      dst,
		Timestamp: time.Now().Format(time.RFC3339),
	}, nil
}

// Matcher returns true if the URI uses an OCI scheme or matches a known OCI registry.
func (o *OCIGatherer) Matcher(uri string) bool {
	prefixes := []string{"oci://", "oci::"}
	for _, prefix := range prefixes {
		if strings.HasPrefix(uri, prefix) {
			return true
		}
	}
	// Check if the input matches any known OCI registry
	return containsOCIRegistry(uri)
}

// Get returns the OCIMetadata value.
func (o OCIMetadata) Get() interface{} {
	return o
}

// GetDigest returns the digest string of the pulled OCI artifact.
func (o OCIMetadata) GetDigest() string {
	return o.Digest
}

// GetPinnedURL returns an oci:: URL pinned to the artifact digest.
func (o OCIMetadata) GetPinnedURL(u string) (string, error) {
	if len(u) == 0 {
		return "", fmt.Errorf("empty URL")
	}
	if o.Digest == "" {
		return "", fmt.Errorf("image digest not set")
	}
	for _, scheme := range []string{"oci::", "oci://", "https://"} {
		u = strings.TrimPrefix(u, scheme)
	}
	parts := strings.Split(u, "@")
	if len(parts) > 1 {
		u = parts[0]
	}
	return fmt.Sprintf("oci::%s@%s", u, o.Digest), nil
}

// containsOCIRegistry checks if the input string's hostname matches a known OCI registry.
func containsOCIRegistry(src string) bool {
	host := extractHost(src)
	for _, matchRegistry := range ociRegistryPatterns {
		if matchRegistry.MatchString(host) {
			return true
		}
	}
	return false
}

// extractHost returns the lowercase hostname (without port) from a URI,
// handling scheme prefixes like "oci://", "oci::", "https://", and bare "host/path" forms.
func extractHost(src string) string {
	// Strip go-getter style "scheme::" prefixes
	if idx := strings.Index(src, "::"); idx != -1 {
		src = src[idx+2:]
	}

	// Ensure a scheme so url.Parse treats the host correctly
	if !strings.Contains(src, "://") {
		src = "oci://" + src
	}

	u, err := url.Parse(src)
	if err != nil || u.Host == "" {
		return strings.ToLower(src)
	}
	return strings.ToLower(u.Hostname())
}

func ociURLParse(source string) string {
	if strings.Contains(source, "::") {
		source = strings.Split(source, "::")[1]
	}

	scheme, src, found := strings.Cut(source, "://")
	if !found {
		src = scheme
	}
	return src
}

func init() {
	gather.RegisterGatherer(NewOCIGatherer(
		WithTransport(retry.NewTransport(http.DefaultTransport)),
	))
}
