// This file is auto-generated. DO NOT EDIT.
// SPDX-License-Identifier: BSD-3-Clause
// Copyright (c) 2025, Unikraft GmbH.
// Licensed under the BSD-3-Clause License (the "License").
// You may not use this file except in compliance with the License.

package platform

// A volume defines a storage which can be attached to the instance.
//
// Volumes can be used to store persistent data which should remain available
// even if the instance is stopped or restarted.

type InstanceVolume struct {
	// The UUID of the volume.
	//
	// This is a unique identifier for the volume that is generated when the
	// volume is created.  The UUID is used to reference the volume in API calls
	// and can be used to identify the volume in all API calls that require a
	// volume identifier.
	Uuid *string `json:"uuid,omitempty"`
	// The name of the volume.
	//
	// This is a human-readable name that can be used to identify the volume.
	// The name must be unique within the context of your account.  The name can
	// also be used to identify the volume in API calls.
	Name *string `json:"name,omitempty"`
	// The mount point of the volume in the instance.  This is the directory in
	// the instance where the volume will be mounted.
	At *string `json:"at,omitempty"`
	// Whether the volume is read-only or not.
	Readonly *bool `json:"readonly,omitempty"`
}
