// This file is auto-generated. DO NOT EDIT.
// SPDX-License-Identifier: BSD-3-Clause
// Copyright (c) 2025, Unikraft GmbH.
// Licensed under the BSD-3-Clause License (the "License").
// You may not use this file except in compliance with the License.

package platform

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"strings"
)

// The response for an invocation to an API on the Unikraft Cloud Platform.
//
// It uses standard HTTP response codes to indicate success or failure.  In
// addition, the response body contains more details about the result of the
// operation in a JSON object.  On success the data member contains an array of
// objects with one object per result.  The array is named according to the type
// of object that the request is operating on.  For example, when working with
// instances the response contains an instances array.
type Response[T any] struct {
	// Status contains the top-level information about a server response, and
	// returns either `success`, `partial_success` or `error`.
	Status string `json:"status,omitempty"`

	// Message contains the error message either on `partial_success` or `error`.
	Message string `json:"message,omitempty"`

	// Errors are the list of errors that have occurred.
	Errors []ResponseError `json:"errors,omitempty"`

	// On a successful response, the data element is returned with relevant
	// information.
	Data *T `json:"data,omitempty"`

	// Buffer holding the raw API response body.
	body bytes.Buffer

	// The underlying HTTP response, if available.
	raw *http.Response

	// The channel to stream events from the response body.
	// This is only used for streaming responses (e.g., SSE).
	events chan *Response[T]
}

// RawBody returns the raw API response body.
func (r *Response[T]) RawBody() []byte {
	b := make([]byte, r.body.Len())
	_, _ = r.body.Read(b) // cannot fail on a bytes.Buffer
	return b
}

// Events returns a channel that streams events from the response body.
func (r *Response[T]) Events() (<-chan *Response[T], error) {
	if r == nil || r.events == nil {
		return nil, errors.New("no events available in response")
	}

	return r.events, nil
}

func (r *Response[T]) Error() string {
	if r == nil || r.Status == "success" {
		return ""
	}

	// Use reflection to determine if the response data type is an array containing
	// individual sub-errors.  This will be re-worked in the future when
	// a top-level errors attribute is properly populated.
	var errs []string

	v := reflect.ValueOf(r.Data)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	if v.Kind() == reflect.Struct && v.NumField() == 1 {
		innerData := v.Field(0)
		if innerData.IsValid() && innerData.Kind() == reflect.Slice {
			for i := 0; i < innerData.Len(); i++ {
				item := innerData.Index(i)
				if item.Kind() == reflect.Ptr {
					item = item.Elem()
				}

				status := item.FieldByName("Status")
				if status.Kind() == reflect.Ptr {
					status = status.Elem()
				}
				message := item.FieldByName("Message")
				if message.Kind() == reflect.Ptr {
					message = message.Elem()
				}

				if status.IsValid() && message.IsValid() && status.String() == "error" {
					errs = append(errs, message.String())
				}
			}
		}
	}

	if len(errs) > 0 {
		return fmt.Sprintf("API error: status code %d: %s: %s", r.raw.StatusCode, r.Message, strings.Join(errs, "; "))
	}

	return fmt.Sprintf("API error: status code %d: %s", r.raw.StatusCode, r.Message)
}
