// This file is auto-generated. DO NOT EDIT.
// SPDX-License-Identifier: BSD-3-Clause
// Copyright (c) 2025, Unikraft GmbH.
// Licensed under the BSD-3-Clause License (the "License").
// You may not use this file except in compliance with the License.

package platform

// The request message to create an autoscale configuration for a service group
// based on its UUID.

type CreateAutoscaleConfigurationByServiceGroupUUIDRequest struct {
	// The UUID of the service to create a configuration for.
	// Mutually exclusive with name.
	Uuid *string `json:"uuid,omitempty"`
	// The minimum number of instances to keep running.
	MinSize *int64 `json:"min_size,omitempty"`
	// The maximum number of instances to keep running.
	MaxSize *int64 `json:"max_size,omitempty"`
	// The warmup time in milliseconds for new instances.
	WarmupTimeMs *int64 `json:"warmup_time_ms,omitempty"`
	// The cooldown time in milliseconds for the autoscale configuration.
	CooldownTimeMs *int64                                                           `json:"cooldown_time_ms,omitempty"`
	CreateArgs     *CreateAutoscaleConfigurationByServiceGroupUUIDRequestCreateArgs `json:"create_args,omitempty"`
	// The policies to apply to the autoscale configuration.
	Policies []AutoscalePolicy `json:"policies,omitempty"`
}
