// This file is auto-generated. DO NOT EDIT.
// SPDX-License-Identifier: BSD-3-Clause
// Copyright (c) 2025, Unikraft GmbH.
// Licensed under the BSD-3-Clause License (the "License").
// You may not use this file except in compliance with the License.

package platform

// Additional limits

type QuotasLimits struct {
	// Minimum amount of memory assigned to live instances in megabytes
	MinMemoryMb *int64 `json:"min_memory_mb,omitempty"`
	// Maximum amount of memory assigned to live instances in megabytes
	MaxMemoryMb *int64 `json:"max_memory_mb,omitempty"`
	// Minimum size of a volume in megabytes
	MinVolumeMb *int64 `json:"min_volume_mb,omitempty"`
	// Maximum size of a volume in megabytes
	MaxVolumeMb *int64 `json:"max_volume_mb,omitempty"`
	// Minimum size of an autoscale group
	MinAutoscaleSize *int64 `json:"min_autoscale_size,omitempty"`
	// Maximum size of an autoscale group
	MaxAutoscaleSize *int64 `json:"max_autoscale_size,omitempty"`
	// Minimum number of vCPUs
	MinVcpus *int64 `json:"min_vcpus,omitempty"`
	// Maximum number of vCPUs
	MaxVcpus *int64 `json:"max_vcpus,omitempty"`
}
