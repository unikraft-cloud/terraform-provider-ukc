// This file is auto-generated. DO NOT EDIT.
// SPDX-License-Identifier: BSD-3-Clause
// Copyright (c) 2025, Unikraft GmbH.
// Licensed under the BSD-3-Clause License (the "License").
// You may not use this file except in compliance with the License.

package platform

// The scale-to-zero configuration for the instance.
//
// With conventional cloud platforms you need to keep at least one instance
// running at all times to be able to respond to incoming requests. Performing
// a just-in-time cold boot is simply too time-consuming and would create a
// response latency of multiple seconds.  This is not the case with Unikraft
// Cloud.  Instances on Unikraft Cloud are able to cold boot within
// milliseconds, which allows us to perform low-latency scale-to-zero.
//
// To enable scale-to-zero for an instance it is sufficient to add a
// `scale_to_zero` configuration block.  Unikraft Cloud will then put the
// instance into standby if there is no traffic to your service within the
// window of a cooldown period.  When there is new traffic coming in, it is
// automatically started again.
//
// If you have a heavyweight application that takes long to cold boot or has
// bad first request latency (e.g., with JIT compilation) consider to enable
// stateful scale-to-zero.
// The specific policy to use for scaling the instance to zero.
type InstanceScaleToZeroPolicy string

const (
	InstanceScaleToZeroPolicyOn   InstanceScaleToZeroPolicy = "on"
	InstanceScaleToZeroPolicyOff  InstanceScaleToZeroPolicy = "off"
	InstanceScaleToZeroPolicyIdle InstanceScaleToZeroPolicy = "idle"
)

type InstanceScaleToZero struct {
	// Indicates whether scale-to-zero is enabled for the instance.
	Enabled *bool `json:"enabled,omitempty"`
	// The specific policy to use for scaling the instance to zero.
	Policy *InstanceScaleToZeroPolicy `json:"policy,omitempty"`
	// Whether the instance should be stateful when scaled to zero. If set to
	// true, the instance will retain its state (e.g., RAM contents) when scaled
	// to zero.  This is useful for instances that need to maintain their state
	// across scale-to-zero operations.  If set to false, the instance will lose
	// its state when scaled to zero, and it will be restarted from scratch when
	// scaled back up.
	Stateful *bool `json:"stateful,omitempty"`
	// The cooldown time in milliseconds before the instance can be scaled to
	// zero again.  This is useful to prevent rapid scaling to zero and back up,
	// which can lead to performance issues or resource exhaustion.
	CooldownTimeMs *int32 `json:"cooldown_time_ms,omitempty"`
}
