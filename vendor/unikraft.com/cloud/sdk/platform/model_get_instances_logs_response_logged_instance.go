// This file is auto-generated. DO NOT EDIT.
// SPDX-License-Identifier: BSD-3-Clause
// Copyright (c) 2025, Unikraft GmbH.
// Licensed under the BSD-3-Clause License (the "License").
// You may not use this file except in compliance with the License.

package platform

// State of the instance when the logs were retrieved.
type GetInstancesLogsResponseLoggedInstanceState string

const (
	GetInstancesLogsResponseLoggedInstanceStateStopped  GetInstancesLogsResponseLoggedInstanceState = "stopped"
	GetInstancesLogsResponseLoggedInstanceStateStarting GetInstancesLogsResponseLoggedInstanceState = "starting"
	GetInstancesLogsResponseLoggedInstanceStateRunning  GetInstancesLogsResponseLoggedInstanceState = "running"
	GetInstancesLogsResponseLoggedInstanceStateDraining GetInstancesLogsResponseLoggedInstanceState = "draining"
	GetInstancesLogsResponseLoggedInstanceStateStopping GetInstancesLogsResponseLoggedInstanceState = "stopping"
	GetInstancesLogsResponseLoggedInstanceStateTemplate GetInstancesLogsResponseLoggedInstanceState = "template"
	GetInstancesLogsResponseLoggedInstanceStateStandby  GetInstancesLogsResponseLoggedInstanceState = "standby"
)

type GetInstancesLogsResponseLoggedInstance struct {
	// The UUID of the instance.
	Uuid *string `json:"uuid,omitempty"`
	// The name of the instance.
	Name *string `json:"name,omitempty"`
	// Base64 encoded log output of the instance.
	Output    *string                                          `json:"output,omitempty"`
	Available *GetInstancesLogsResponseLoggedInstanceAvailable `json:"available,omitempty"`
	Range     *GetInstancesLogsResponseLoggedInstanceRange     `json:"range,omitempty"`
	// State of the instance when the logs were retrieved.
	State *GetInstancesLogsResponseLoggedInstanceState `json:"state,omitempty"`
	// An optional message providing additional information about the status.
	// This field is useful when the status is not `success`.
	Message *string `json:"message,omitempty"`
	// An optional error code providing additional information about the status.
	// This field is useful when the status is not `success`.
	Error *int32 `json:"error,omitempty"`
}
