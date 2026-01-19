// This file is auto-generated. DO NOT EDIT.
// SPDX-License-Identifier: BSD-3-Clause
// Copyright (c) 2025, Unikraft GmbH.
// Licensed under the BSD-3-Clause License (the "License").
// You may not use this file except in compliance with the License.

package platform

type AttachVolumeByUUIDRequestBody struct {
	// UUID or name of the instance to attach the volume to.
	AttachTo BodyInstanceID `json:"attach_to"`
	// Path of the mountpoint.
	//
	// The path must be absolute, not contain `.` and `..` components, and not
	// contain colons (`:`). The path must point to an empty directory. If the
	// directory does not exist, it is created.
	At string `json:"at"`
	// Whether the volume should be mounted read-only.
	Readonly *bool `json:"readonly,omitempty"`
}
