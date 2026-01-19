// This file is auto-generated. DO NOT EDIT.
// SPDX-License-Identifier: BSD-3-Clause
// Copyright (c) 2025, Unikraft GmbH.
// Licensed under the BSD-3-Clause License (the "License").
// You may not use this file except in compliance with the License.

package platform

// Wait parameters.
// The desired state to wait for.  Default is `running`.
type WaitInstanceByUUIDRequestBodyState string

const (
	WaitInstanceByUUIDRequestBodyStateStopped  WaitInstanceByUUIDRequestBodyState = "stopped"
	WaitInstanceByUUIDRequestBodyStateStarting WaitInstanceByUUIDRequestBodyState = "starting"
	WaitInstanceByUUIDRequestBodyStateRunning  WaitInstanceByUUIDRequestBodyState = "running"
	WaitInstanceByUUIDRequestBodyStateDraining WaitInstanceByUUIDRequestBodyState = "draining"
	WaitInstanceByUUIDRequestBodyStateStopping WaitInstanceByUUIDRequestBodyState = "stopping"
	WaitInstanceByUUIDRequestBodyStateTemplate WaitInstanceByUUIDRequestBodyState = "template"
	WaitInstanceByUUIDRequestBodyStateStandby  WaitInstanceByUUIDRequestBodyState = "standby"
)

type WaitInstanceByUUIDRequestBody struct {
	// The desired state to wait for.  Default is `running`.
	State *WaitInstanceByUUIDRequestBodyState `json:"state,omitempty"`
	// Timeout in milliseconds to wait for the instance to reach the desired
	// state.  If the timeout is reached, the request will fail with an error.
	// A value of -1 means to wait indefinitely until the instance reaches the
	// desired state.
	TimeoutMs *int64 `json:"timeout_ms,omitempty"`
}
