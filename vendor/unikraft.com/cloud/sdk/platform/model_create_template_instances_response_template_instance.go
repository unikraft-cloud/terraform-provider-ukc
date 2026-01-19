// This file is auto-generated. DO NOT EDIT.
// SPDX-License-Identifier: BSD-3-Clause
// Copyright (c) 2025, Unikraft GmbH.
// Licensed under the BSD-3-Clause License (the "License").
// You may not use this file except in compliance with the License.

package platform

// The current state of the instance.
type CreateTemplateInstancesResponseTemplateInstanceState string

const (
	CreateTemplateInstancesResponseTemplateInstanceStateStopped  CreateTemplateInstancesResponseTemplateInstanceState = "stopped"
	CreateTemplateInstancesResponseTemplateInstanceStateStarting CreateTemplateInstancesResponseTemplateInstanceState = "starting"
	CreateTemplateInstancesResponseTemplateInstanceStateRunning  CreateTemplateInstancesResponseTemplateInstanceState = "running"
	CreateTemplateInstancesResponseTemplateInstanceStateDraining CreateTemplateInstancesResponseTemplateInstanceState = "draining"
	CreateTemplateInstancesResponseTemplateInstanceStateStopping CreateTemplateInstancesResponseTemplateInstanceState = "stopping"
	CreateTemplateInstancesResponseTemplateInstanceStateTemplate CreateTemplateInstancesResponseTemplateInstanceState = "template"
	CreateTemplateInstancesResponseTemplateInstanceStateStandby  CreateTemplateInstancesResponseTemplateInstanceState = "standby"
)

type CreateTemplateInstancesResponseTemplateInstance struct {
	// The status of this particular template instance creation operation.
	Status *string `json:"status,omitempty"`
	// The UUID of the template instance that was created.
	Uuid *string `json:"uuid,omitempty"`
	// The name of the template instance that was created.
	Name *string `json:"name,omitempty"`
	// The current state of the instance.
	State *CreateTemplateInstancesResponseTemplateInstanceState `json:"state,omitempty"`
	// An optional message providing additional information about the status.
	// This field is useful when the status is not `success`.
	Message *string `json:"message,omitempty"`
	// An optional error code providing additional information about the status.
	// This field is useful when the status is not `success`.
	Error *int32 `json:"error,omitempty"`
}
