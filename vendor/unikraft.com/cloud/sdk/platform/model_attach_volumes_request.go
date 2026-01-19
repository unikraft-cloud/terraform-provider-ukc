// This file is auto-generated. DO NOT EDIT.
// SPDX-License-Identifier: BSD-3-Clause
// Copyright (c) 2025, Unikraft GmbH.
// Licensed under the BSD-3-Clause License (the "License").
// You may not use this file except in compliance with the License.

package platform

// The request message for attaching one or more volume(s) to instances by
// their UUID(s) or name(s).

type AttachVolumesRequest struct {
	// The UUID of the volume to attach. Mutually exclusive with name.
	// Exactly one of uuid or name must be provided.
	Uuid *string `json:"uuid,omitempty"`
	// The name of the volume to attach. Mutually exclusive with UUID.
	// Exactly one of uuid or name must be provided.
	Name     *string                      `json:"name,omitempty"`
	AttachTo AttachVolumesRequestAttachTo `json:"attach_to"`
	// Path of the mountpoint.
	//
	// The path must be absolute, not contain `.` and `..` components, and not
	// contain colons (`:`). The path must point to an empty directory. If the
	// directory does not exist, it is created.
	At string `json:"at"`
	// Whether the volume should be mounted read-only.
	Readonly *bool `json:"readonly,omitempty"`
}
