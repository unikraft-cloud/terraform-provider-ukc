// This file is auto-generated. DO NOT EDIT.
// SPDX-License-Identifier: BSD-3-Clause
// Copyright (c) 2025, Unikraft GmbH.
// Licensed under the BSD-3-Clause License (the "License").
// You may not use this file except in compliance with the License.

package platform

// A single request item to stop an instance.

type StopInstancesRequestItem struct {
	// The UUID of the instance to stop.  Mutually exclusive with name.
	Uuid *string `json:"uuid,omitempty"`
	// The name of the instance to stop.  Mutually exclusive with UUID.
	Name *string `json:"name,omitempty"`
	// Whether to immediately force stop the instance.
	Force *bool `json:"force,omitempty"`
	// Timeout for draining connections in milliseconds.  The instance does not
	// receive new connections in the draining phase.  The instance is stopped
	// when the last connection has been closed or the timeout expired.  The
	// maximum timeout may vary.  Use -1 for the largest possible value.
	//
	// Note: This endpoint does not block.  Use the wait endpoint for the
	// instance to reach the stopped state.
	DrainTimeoutMs *uint64 `json:"drain_timeout_ms,omitempty"`
}
