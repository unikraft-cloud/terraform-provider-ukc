// This file is auto-generated. DO NOT EDIT.
// SPDX-License-Identifier: BSD-3-Clause
// Copyright (c) 2025, Unikraft GmbH.
// Licensed under the BSD-3-Clause License (the "License").
// You may not use this file except in compliance with the License.

package platform

import (
	"unikraft.com/cloud/sdk/pkg/httpclient"
)

// ClientOption is an option function used during initialization of a client.
type ClientOption func(*ClientOptions)

type ClientOptions struct {
	defaultMetro string

	token         string
	allowInsecure bool
	userAgent     string
	httpClient    httpclient.HTTPClient
}

// SetDefaultMetro sets the default metro.
func (opts *ClientOptions) SetDefaultMetro(metro string) {
	opts.defaultMetro = metro
}

// DefaultMetro retrieves the default metro.
func (opts *ClientOptions) DefaultMetro() string {
	return opts.defaultMetro
}

// SetToken sets the token to use for authentication with the API.
func (opts *ClientOptions) SetToken(token string) {
	opts.token = token
}

// Token returns the token used for authentication with the API.
func (opts *ClientOptions) Token() string {
	return opts.token
}

// SetAllowInsecure sets whether to allow insecure connections (e.g., HTTP
// instead of HTTPS) when making API requests.
func (opts *ClientOptions) SetAllowInsecure(allow bool) {
	opts.allowInsecure = allow
}

// AllowInsecure returns whether insecure connections are allowed when making
// API requests.
func (opts *ClientOptions) AllowInsecure() bool {
	return opts.allowInsecure
}

// SetUserAgent sets the user agent string to use for API requests.
func (opts *ClientOptions) SetUserAgent(ua string) {
	opts.userAgent = ua
}

// UserAgent returns the user agent string used for API requests.
func (opts *ClientOptions) UserAgent() string {
	return opts.userAgent
}

// SetHTTPClient sets the HTTP client to use for making API requests.
func (opts *ClientOptions) SetHTTPClient(client httpclient.HTTPClient) {
	opts.httpClient = client
}

// HTTPClient returns the HTTP client used for making API requests.
func (opts *ClientOptions) HTTPClient() httpclient.HTTPClient {
	return opts.httpClient
}

// WithDefaultMetro sets the default metro for the client.
func WithDefaultMetro(metro string) ClientOption {
	return func(client *ClientOptions) {
		client.SetDefaultMetro(metro)
	}
}

// WithToken sets the access token of the client connecting to KraftCloud.
func WithToken(token string) ClientOption {
	return func(client *ClientOptions) {
		client.SetToken(token)
	}
}

// WithHTTPClient sets the HTTP client which is used when connecting to the
// API.
func WithHTTPClient(httpClient httpclient.HTTPClient) ClientOption {
	return func(client *ClientOptions) {
		client.SetHTTPClient(httpClient)
	}
}

// WithUserAgent configures the to use the given user agent string in API
// requests.
func WithUserAgent(ua string) ClientOption {
	return func(client *ClientOptions) {
		client.SetUserAgent(ua)
	}
}

// WithAllowInsecure configures the client to skip checks of the API
// certificates.
func WithAllowInsecure(allow bool) ClientOption {
	return func(client *ClientOptions) {
		client.SetAllowInsecure(allow)
	}
}
