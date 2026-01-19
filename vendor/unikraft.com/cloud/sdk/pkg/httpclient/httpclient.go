// SPDX-License-Identifier: BSD-3-Clause
// Copyright (c) 2022, Unikraft GmbH and The KraftKit Authors.
// Licensed under the BSD-3-Clause License (the "License").
// You may not use this file except in compliance with the License.

// Package httpclient provides an HTTP client interface and default
// implementations.
package httpclient

import (
	"crypto/tls"
	"net/http"
	"net/url"

	"golang.org/x/net/http/httpproxy"
)

// HTTPClient interface abstracts a generic HTTP request issuing client.
type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// NewHTTPClient creates a default Go HTTP client.
func NewHTTPClient() *http.Client {
	// We disable KeepAlive due to issues with the proxy in front of the API.
	return &http.Client{
		Transport: &http.Transport{
			DisableKeepAlives: true,
			Proxy: func(req *http.Request) (*url.URL, error) {
				return httpproxy.FromEnvironment().ProxyFunc()(req.URL)
			},
		},
	}
}

// NewInsecureHTTPClient creates a default Go HTTP client with insecure checks
// skipped.
func NewInsecureHTTPClient() *http.Client {
	return &http.Client{
		Transport: &http.Transport{
			DisableKeepAlives: true,
			Proxy: func(req *http.Request) (*url.URL, error) {
				return httpproxy.FromEnvironment().ProxyFunc()(req.URL)
			},
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true, // Allow insecure connections
			},
		},
	}
}
