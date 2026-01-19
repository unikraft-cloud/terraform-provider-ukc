// This file is auto-generated. DO NOT EDIT.
// SPDX-License-Identifier: BSD-3-Clause
// Copyright (c) 2025, Unikraft GmbH.
// Licensed under the BSD-3-Clause License (the "License").
// You may not use this file except in compliance with the License.

package platform

// A volume defines a storage volume that can be attached to the instance.

type CreateInstanceRequestVolume struct {
	// The UUID of an existing volume.
	//
	// If this is the only specified field, then it will look up an existing
	// volume by this UUID.
	Uuid *string `json:"uuid,omitempty"`
	// The name of the volume.
	//
	// If this is the only specified field, then it will look up an existing
	// volume by this name.  If the volume does not exist, the request will
	// fail.  If a new volume is intended to be created, then this field must be
	// specified along with the size in MiB and the mount point in the instance.
	Name *string `json:"name,omitempty"`
	// The size of the volume when creating a new volume.
	//
	// When creating a new volume as part of the instance create request,
	// specify the size of the volume in MiB.
	SizeMb *int64 `json:"size_mb,omitempty"`
	// The mount point for the volume in the instance.
	At string `json:"at"`
	// Whether the volume is read-only.
	//
	// If this field is set to true, the volume will be mounted as read-only in
	// the instance.  This field is optional and defaults to false and is only
	// applicable when using an existing volume.
	Readonly *bool `json:"readonly,omitempty"`
}
