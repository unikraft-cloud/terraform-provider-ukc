// This file is auto-generated. DO NOT EDIT.
// SPDX-License-Identifier: BSD-3-Clause
// Copyright (c) 2025, Unikraft GmbH.
// Licensed under the BSD-3-Clause License (the "License").
// You may not use this file except in compliance with the License.

package platform

// AutoscalePolicy defines the autoscale policy for a service.
// Right now it contains fields from both the `ondemand` and `step` policies.
// They are marked both as optional, so only one of them should be set at a
// time. This is a current limitation of the API design.
// Metric to use for the step policy.
type AutoscalePolicyMetric string

const (
	AutoscalePolicyMetricCpu AutoscalePolicyMetric = "cpu"
)

// The type of adjustment to be made in the step policy.
type AutoscalePolicyAdjustmentType string

const (
	AutoscalePolicyAdjustmentTypeChange     AutoscalePolicyAdjustmentType = "change"
	AutoscalePolicyAdjustmentTypeExact      AutoscalePolicyAdjustmentType = "exact"
	AutoscalePolicyAdjustmentTypePercentage AutoscalePolicyAdjustmentType = "percentage"
)

type AutoscalePolicy struct {
	// The name of the policy.
	Name *string `json:"name,omitempty"`
	// If the policy is enabled.
	Enabled *bool `json:"enabled,omitempty"`
	// Metric to use for the step policy.
	Metric *AutoscalePolicyMetric `json:"metric,omitempty"`
	// The type of adjustment to be made in the step policy.
	AdjustmentType *AutoscalePolicyAdjustmentType `json:"adjustment_type,omitempty"`
	// The steps for the step policy.
	// Each step defines an adjustment value and optional bounds.
	Steps []AutoscalePolicyStep `json:"steps,omitempty"`
}
