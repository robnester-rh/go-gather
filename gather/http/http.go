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

// Package http implements a Gatherer for HTTP and HTTPS sources.
package http

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/conforma/go-gather/gather"
	"github.com/conforma/go-gather/internal/helpers"
	"github.com/conforma/go-gather/metadata"
)

// Option configures an HTTPGatherer.
type Option func(*HTTPGatherer)

// WithTransport sets the http.RoundTripper used by the HTTPGatherer's client.
func WithTransport(t http.RoundTripper) Option {
	return func(g *HTTPGatherer) { g.Client.Transport = t }
}

// HTTPGatherer gathers resources over HTTP/HTTPS.
type HTTPGatherer struct {
	Client http.Client
}

// HTTPMetadata holds metadata about a gathered HTTP resource.
type HTTPMetadata struct {
	URI          string
	Path         string
	ResponseCode int
	Size         int64
	Timestamp    string
}

// NewHTTPGatherer returns an HTTPGatherer with a default 30-second timeout.
func NewHTTPGatherer(opts ...Option) *HTTPGatherer {
	g := &HTTPGatherer{
		Client: http.Client{Timeout: 30 * time.Second},
	}
	for _, opt := range opts {
		opt(g)
	}
	return g
}

// Gather downloads a file from rawSource via HTTP and writes it to dst.
func (h *HTTPGatherer) Gather(ctx context.Context, rawSource, dst string) (meta metadata.Metadata, err error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	src, err := url.Parse(rawSource)
	if err != nil {
		return nil, fmt.Errorf("failed to parse source URI: %w", err)
	}

	// Check if the source scheme is provided
	if src.Scheme == "" {
		return nil, fmt.Errorf("no source scheme provided")
	}

	// Check if the source filename is provided
	if src.Path == "" {
		return nil, fmt.Errorf("specify a path to a file to download")
	}

	// Get the source filename
	sourceFileName := filepath.Base(src.Path)

	// Expand the destination path
	dst, err = helpers.ExpandPath(dst)
	if err != nil {
		return nil, fmt.Errorf("failed to expand destination path: %w", err)
	}

	// Check if the destination has a trailing slash.
	// If it does, append the source filename to the destination path.
	if strings.HasSuffix(dst, "/") {
		dst = filepath.Join(dst, sourceFileName)
	} else {
		// If it doesn't, append the source filename to the destination path.
		if filepath.Ext(dst) == "" {
			dst = filepath.Join(dst, "/", sourceFileName)
		}
	}

	// Create a new HTTP request
	req, err := http.NewRequestWithContext(ctx, "GET", rawSource, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}

	// Set the User-Agent header
	req.Header.Set("User-Agent", "Go-Gather")

	// Perform the HTTP request
	resp, err := h.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to download from URL: %w", err)
	}
	defer resp.Body.Close()

	// Check if the response code is "ok"
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("received non-200 response code: %d", resp.StatusCode)
	}

	// Create the destination file
	err = os.MkdirAll(filepath.Dir(dst), 0755)
	if err != nil {
		return nil, fmt.Errorf("failed to create destination directory: %w", err)
	}
	outFile, err := os.Create(dst)
	if err != nil {
		return nil, fmt.Errorf("failed to create destination file: %w", err)
	}
	defer func() {
		if cerr := outFile.Close(); cerr != nil && err == nil {
			meta = nil
			err = fmt.Errorf("failed to close destination file: %w", cerr)
		}
	}()

	bytesWritten, err := io.Copy(outFile, resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to write to destination file: %w", err)
	}

	return HTTPMetadata{
		URI:          rawSource,
		Path:         dst,
		ResponseCode: resp.StatusCode,
		Size:         bytesWritten,
		Timestamp:    time.Now().Format(time.RFC3339),
	}, nil
}

// Matcher returns true if the URI uses an HTTP or HTTPS scheme and is not a known git host.
func (h *HTTPGatherer) Matcher(uri string) bool {
	u, err := url.Parse(uri)
	if err != nil {
		return false
	}
	if !strings.EqualFold(u.Scheme, "http") && !strings.EqualFold(u.Scheme, "https") {
		return false
	}
	host := strings.ToLower(u.Hostname())
	if host == "" {
		return false
	}
	for _, vendor := range []string{"github.com", "gitlab.com", "bitbucket.org"} {
		if host == vendor || strings.HasSuffix(host, "."+vendor) {
			return false
		}
	}
	return true
}

// Get returns the HTTPMetadata value.
func (h HTTPMetadata) Get() interface{} {
	return h
}

// GetPinnedURL returns an http:: prefixed URL for the given address.
func (h HTTPMetadata) GetPinnedURL(u string) (string, error) {
	if len(u) == 0 {
		return "", fmt.Errorf("empty URL")
	}
	for _, scheme := range []string{"http://", "https://", "http::"} {
		u = strings.TrimPrefix(u, scheme)
	}
	return "http::" + u, nil
}

func init() {
	gather.RegisterGatherer(NewHTTPGatherer())
}
