// This file is auto-generated. DO NOT EDIT.
// SPDX-License-Identifier: BSD-3-Clause
// Copyright (c) 2025, Unikraft GmbH.
// Licensed under the BSD-3-Clause License (the "License").
// You may not use this file except in compliance with the License.

package platform

// UUID or name of the instance to detach the volume from.

type DetachVolumeByUUIDRequestBodyFrom struct {
	// The UUID of the instance that the volume is detached from.
	Uuid *string `json:"uuid,omitempty"`
	// The name of the instance that the volume is detached from.
	Name *string `json:"name,omitempty"`
}
