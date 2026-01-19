// This file is auto-generated. DO NOT EDIT.
// SPDX-License-Identifier: BSD-3-Clause
// Copyright (c) 2025, Unikraft GmbH.
// Licensed under the BSD-3-Clause License (the "License").
// You may not use this file except in compliance with the License.

package platform

type DetachVolumesResponseDetachedVolume struct {
	// The status of the response.
	Status *ResponseStatus `json:"status,omitempty"`
	// The UUID of the volume that was detached.
	Uuid *string `json:"uuid,omitempty"`
	// The name of the volume that was detached.
	Name *string `json:"name,omitempty"`
	// An optional message providing additional information about the status.
	// This field is useful when the status is not `success`.
	Message *string `json:"message,omitempty"`
	// An optional error code providing additional information about the status.
	// This field is useful when the status is not `success`.
	Error *int32 `json:"error,omitempty"`
}
