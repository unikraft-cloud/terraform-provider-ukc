// This file is auto-generated. DO NOT EDIT.
// SPDX-License-Identifier: BSD-3-Clause
// Copyright (c) 2025, Unikraft GmbH.
// Licensed under the BSD-3-Clause License (the "License").
// You may not use this file except in compliance with the License.

package platform

import "time"

// A volume represents a storage device that can be attached to an instance.
// Current state of the volume.
type VolumeState string

const (
	VolumeStateUninitialized VolumeState = "uninitialized"
	VolumeStateInitializing  VolumeState = "initializing"
	VolumeStateAvailable     VolumeState = "available"
	VolumeStateIdle          VolumeState = "idle"
	VolumeStateMounted       VolumeState = "mounted"
	VolumeStateBusy          VolumeState = "busy"
	VolumeStateError         VolumeState = "error"
)

type Volume struct {
	// The UUID of the volume.
	//
	// This is a unique identifier for the volume that is generated when the
	// volume is created.  The UUID is used to reference the volume in
	// API calls and can be used to identify the volume in all API calls that
	// require an identifier.
	Uuid *string `json:"uuid,omitempty"`
	// The name of the volume.
	//
	// This is a human-readable name that can be used to identify the volume.
	// The name must be unique within the context of your account.  The name can
	// also be used to identify the volume in API calls.
	Name *string `json:"name,omitempty"`
	// The time the volume was created.
	CreatedAt *time.Time `json:"created_at,omitempty"`
	// Current state of the volume.
	State *VolumeState `json:"state,omitempty"`
	// The size of the volume in megabytes.
	SizeMb *uint64 `json:"size_mb,omitempty"`
	// Indicates if the volume will stay alive when the last instance is deleted
	// that this volume is attached to.
	Persistent *bool `json:"persistent,omitempty"`
	// List of instances that this volume is attached to.
	AttachedTo []VolumeInstanceID `json:"attached_to,omitempty"`
	// List of instances that have this volume mounted.
	MountedBy []VolumeVolumeInstanceMount `json:"mounted_by,omitempty"`
	// The tags associated with the volume.
	Tags []string `json:"tags,omitempty"`
	// An optional field representing the status of the request.  This field is
	// only set when this message object is used as a response message.
	Status *ResponseStatus `json:"status,omitempty"`
	// An optional message providing additional information about the status.
	// This field is only set when this message object is used as a response
	// message, and is useful when the status is not `success`.
	Message *string `json:"message,omitempty"`
	// An optional error code providing additional information about the status.
	// This field is only set when this message object is used as a response
	// message, and is useful when the status is not `success`.
	Error *int32 `json:"error,omitempty"`
}
