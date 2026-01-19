// This file is auto-generated. DO NOT EDIT.
// SPDX-License-Identifier: BSD-3-Clause
// Copyright (c) 2025, Unikraft GmbH.
// Licensed under the BSD-3-Clause License (the "License").
// You may not use this file except in compliance with the License.

package platform

// The current state of the instance.
type StopInstancesResponseStoppedInstanceState string

const (
	StopInstancesResponseStoppedInstanceStateStopped  StopInstancesResponseStoppedInstanceState = "stopped"
	StopInstancesResponseStoppedInstanceStateStarting StopInstancesResponseStoppedInstanceState = "starting"
	StopInstancesResponseStoppedInstanceStateRunning  StopInstancesResponseStoppedInstanceState = "running"
	StopInstancesResponseStoppedInstanceStateDraining StopInstancesResponseStoppedInstanceState = "draining"
	StopInstancesResponseStoppedInstanceStateStopping StopInstancesResponseStoppedInstanceState = "stopping"
	StopInstancesResponseStoppedInstanceStateTemplate StopInstancesResponseStoppedInstanceState = "template"
	StopInstancesResponseStoppedInstanceStateStandby  StopInstancesResponseStoppedInstanceState = "standby"
)

// The previous state of the instance before the stop operation was invoked.
type StopInstancesResponseStoppedInstancePreviousState string

const (
	StopInstancesResponseStoppedInstancePreviousStateStopped  StopInstancesResponseStoppedInstancePreviousState = "stopped"
	StopInstancesResponseStoppedInstancePreviousStateStarting StopInstancesResponseStoppedInstancePreviousState = "starting"
	StopInstancesResponseStoppedInstancePreviousStateRunning  StopInstancesResponseStoppedInstancePreviousState = "running"
	StopInstancesResponseStoppedInstancePreviousStateDraining StopInstancesResponseStoppedInstancePreviousState = "draining"
	StopInstancesResponseStoppedInstancePreviousStateStopping StopInstancesResponseStoppedInstancePreviousState = "stopping"
	StopInstancesResponseStoppedInstancePreviousStateTemplate StopInstancesResponseStoppedInstancePreviousState = "template"
	StopInstancesResponseStoppedInstancePreviousStateStandby  StopInstancesResponseStoppedInstancePreviousState = "standby"
)

type StopInstancesResponseStoppedInstance struct {
	// The UUID of the instance.
	Uuid *string `json:"uuid,omitempty"`
	// The name of the instance.
	Name *string `json:"name,omitempty"`
	// The current state of the instance.
	State *StopInstancesResponseStoppedInstanceState `json:"state,omitempty"`
	// The previous state of the instance before the stop operation was invoked.
	PreviousState *StopInstancesResponseStoppedInstancePreviousState `json:"previous_state,omitempty"`
	// An optional message providing additional information about the status.
	// This field is useful when the status is not `success`.
	Message *string `json:"message,omitempty"`
	// An optional error code providing additional information about the status.
	// This field is useful when the status is not `success`.
	Error *int32 `json:"error,omitempty"`
}
