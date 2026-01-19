// This file is auto-generated. DO NOT EDIT.
// SPDX-License-Identifier: BSD-3-Clause
// Copyright (c) 2025, Unikraft GmbH.
// Licensed under the BSD-3-Clause License (the "License").
// You may not use this file except in compliance with the License.

package platform

// Reference to the instance to attach the volume to.

type AttachVolumesRequestInstanceID struct {
	// The UUID of the instance that the volume is attached to.
	Uuid *string `json:"uuid,omitempty"`
	// The name of the instance that the volume is attached to.
	Name *string `json:"name,omitempty"`
}
