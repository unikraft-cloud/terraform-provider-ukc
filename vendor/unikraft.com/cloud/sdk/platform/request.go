// This file is auto-generated. DO NOT EDIT.
// SPDX-License-Identifier: BSD-3-Clause
// Copyright (c) 2025, Unikraft GmbH.
// Licensed under the BSD-3-Clause License (the "License").
// You may not use this file except in compliance with the License.

package platform

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime"
	"net/http"
	"net/url"
	"strings"
	"time"

	"unikraft.com/cloud/sdk/pkg/httpclient"
	"unikraft.com/cloud/sdk/pkg/sse"
)

// Request is the utility structure for performing individual requests to
// a service location at KraftCloud.
type Request struct {
	copts      *ClientOptions
	metro      string
	httpClient httpclient.HTTPClient
	timeout    time.Duration
}

// NewRequestFromDefaultOptions is a constructor method which uses the
// prebuilt set of options as part of the request.
func NewRequestFromDefaultOptions(copts *ClientOptions) *Request {
	return &Request{
		copts: copts,
	}
}

// WithMetro returns a Request that uses the given metro in API requests.
func (r *Request) WithMetro(m string) *Request {
	rcpy := r.clone()
	rcpy.metro = m
	return rcpy
}

// WithTimeout returns a Request that uses the specified timeout
// duration in API requests.
func (r *Request) WithTimeout(to time.Duration) *Request {
	rcpy := r.clone()
	rcpy.timeout = to
	return rcpy
}

// WithHTTPClient returns a Request which performs API requests using
// the given HTTPClient.
func (r *Request) WithHTTPClient(hc httpclient.HTTPClient) *Request {
	rcpy := r.clone()
	rcpy.httpClient = hc
	return rcpy
}

// Metro returns the metro that this request will perform against.
func (r *Request) Metro() string {
	return r.metro
}

// GetBearerToken uses the pre-defined token to construct the header used for
// authenticating requests.
func (r *Request) GetBearerToken() string {
	return "Bearer " + r.copts.Token()
}

// GetToken uses the pre-defined token to construct the Bearer token used
// for authenticating requests.
func (r *Request) GetToken() string {
	return r.copts.Token()
}

// GetUserAgent uses the pre-defined token to construct the Bearer token used
// for authenticating requests.
func (r *Request) GetUserAgent() string {
	return r.copts.UserAgent()
}

// clone returns a shallow copy of r.
func (r *Request) clone() *Request {
	rcpy := *r
	return &rcpy
}

type (
	// RequestCallback is a function that can be used to modify the request
	// before it is sent.
	RequestCallback func(*http.Request)

	// ResponseCallback is a function that can be used to modify the response
	// after it is received.
	ResponseCallback func(*http.Request, *http.Response)
)

// requestModifiers contains headers, callbacks, etc. that will be added to
// each request
type requestModifiers struct {
	headers http.Header

	requestCallbacks  []RequestCallback
	responseCallbacks []ResponseCallback

	// additionalQueryParameters, if specified will be appended to the list of
	// query parameters included with the request.
	additionalQueryParameters url.Values
}

// RequestOption is a functional parameter used to modify a request
type RequestOption func(*requestModifiers) error

// WithRequestCallbacks sets callbacks which will be invoked before the next
// request; these callbacks will be appended to the list of the callbacks set at
// the client-level.
func WithRequestCallbacks(callbacks ...RequestCallback) RequestOption {
	return func(m *requestModifiers) error {
		m.requestCallbacks = callbacks
		return nil
	}
}

// WithResponseCallbacks sets callbacks which will be invoked after a successful
// response within the next request; these callbacks will be appended to the
// list of the callbacks set at the client-level.
func WithResponseCallbacks(callbacks ...ResponseCallback) RequestOption {
	return func(m *requestModifiers) error {
		m.responseCallbacks = callbacks
		return nil
	}
}

// WithQueryParameters appends the given list of query parameters to the
// request. This is primarily intended to be used with client.Read,
// client.ReadRaw, client.List, and client.Delete methods.
func WithQueryParameters(parameters url.Values) RequestOption {
	return func(m *requestModifiers) error {
		m.additionalQueryParameters = parameters
		return nil
	}
}

// WithCustomHeaders sets custom headers for the next request; these headers
// will be appended to any custom headers set at the client level.
func WithCustomHeaders(headers http.Header) RequestOption {
	return func(m *requestModifiers) error {
		m.headers = headers
		return nil
	}
}

// doRequest performs the request and handles the return media type differently
// depending on the media type:
//
//   - application/json: hydrates a target type with response body and closes the
//     body.
//   - text/event-stream: returns a channel of events which will be closed when
//     the context is done or the connection is closed.  The channel will contain
//     pointers to Response items, which are decoded from the event data.
func doRequest[T any](ctx context.Context, req *Request, method, path string, query url.Values, reqBody io.Reader, target *Response[T], ropts ...RequestOption) error {
	var m string
	var u *url.URL
	var err error

	if req.metro != "" {
		m = req.metro
	} else {
		m = req.copts.DefaultMetro()
	}

	// If the metro contains a full URL, quantified by the presence of a scheme,
	// we assume it is a full URL to a metro, otherwise place it within the
	// well-defined format URL.
	if !strings.Contains(m, "://") {
		m = fmt.Sprintf(BaseV1FormatURL, m)
	}
	u, err = url.Parse(m)
	if err != nil {
		return fmt.Errorf("error constructing URL: %w", err)
	}

	u = u.JoinPath(path)
	u.RawQuery = query.Encode()

	httpReq, err := http.NewRequestWithContext(ctx, method, u.String(), reqBody)
	if err != nil {
		return fmt.Errorf("error creating the request: %w", err)
	}

	httpReq.Header.Set("Authorization", req.GetBearerToken())
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Accept", "application/json")
	httpReq.Header.Set("User-Agent", req.GetUserAgent())

	hc := req.copts.HTTPClient()
	if req.httpClient != nil {
		hc = req.httpClient
	}

	resp, err := hc.Do(httpReq)
	if err != nil {
		return fmt.Errorf("error performing the request: %w", err)
	}

	// Store the raw response in the target.
	target.raw = resp

	ct, _, err := mime.ParseMediaType(resp.Header.Get("Content-Type"))
	if err != nil {
		return fmt.Errorf("parsing media type: %v", err)
	}

	switch ct {
	case "text/event-stream":
		if resp.StatusCode != 200 {
			defer resp.Body.Close()

			respBody, err := io.ReadAll(resp.Body)
			if err != nil {
				return fmt.Errorf("reading response body: %w", err)
			}

			if _, err = io.Copy(&target.body, bytes.NewReader(respBody)); err != nil {
				return fmt.Errorf("storing response body: %w", err)
			}

			bodyReader := bytes.NewReader(respBody)

			if err := json.NewDecoder(bodyReader).Decode(target); err != nil {
				return fmt.Errorf("error parsing response: %w", err)
			}

			return target // implements error
		}

		target.events = make(chan *Response[T])

		// Start a goroutine to process the event stream.
		go func() {
			defer close(target.events)

			// Create an SSE reader.
			reader := sse.NewReader(resp.Body)

			// Process events until context is done or connection is closed
			for {
				select {
				case <-resp.Request.Context().Done():
					return
				default:
					event, err := reader.ReadEvent()
					if err != nil {
						return
					}

					// Skip if event data is empty or it's a keep-alive
					if len(event.Data) == 0 || bytes.Equal(event.Data, []byte(":keep-alive")) {
						continue
					}

					var item Response[T]
					if err := json.Unmarshal(event.Data, &item); err != nil {
						continue
					}

					// Send the decoded item to the channel
					select {
					case target.events <- &item:
					case <-resp.Request.Context().Done():
						return
					}
				}
			}
		}()

	default: // case "application/json":
		defer resp.Body.Close()

		respBody, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("reading response body: %w", err)
		}

		if _, err = io.Copy(&target.body, bytes.NewReader(respBody)); err != nil {
			return fmt.Errorf("storing response body: %w", err)
		}

		bodyReader := bytes.NewReader(respBody)

		if err := json.NewDecoder(bodyReader).Decode(target); err != nil {
			return fmt.Errorf("error parsing response: %w", err)
		}

		if target.Status != "success" {
			return target // implements error
		}
	}

	return nil
}
