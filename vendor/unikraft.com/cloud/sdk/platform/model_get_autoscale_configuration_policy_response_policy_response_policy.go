// This file is auto-generated. DO NOT EDIT.
// SPDX-License-Identifier: BSD-3-Clause
// Copyright (c) 2025, Unikraft GmbH.
// Licensed under the BSD-3-Clause License (the "License").
// You may not use this file except in compliance with the License.

package platform

// The policy which was retrieved by the request.
// Metric to use for the step policy.
type GetAutoscaleConfigurationPolicyResponsePolicyResponsePolicyMetric string

const (
	GetAutoscaleConfigurationPolicyResponsePolicyResponsePolicyMetricCpu GetAutoscaleConfigurationPolicyResponsePolicyResponsePolicyMetric = "cpu"
)

// The type of adjustment to be made in the step policy.
type GetAutoscaleConfigurationPolicyResponsePolicyResponsePolicyAdjustmentType string

const (
	GetAutoscaleConfigurationPolicyResponsePolicyResponsePolicyAdjustmentTypeChange     GetAutoscaleConfigurationPolicyResponsePolicyResponsePolicyAdjustmentType = "change"
	GetAutoscaleConfigurationPolicyResponsePolicyResponsePolicyAdjustmentTypeExact      GetAutoscaleConfigurationPolicyResponsePolicyResponsePolicyAdjustmentType = "exact"
	GetAutoscaleConfigurationPolicyResponsePolicyResponsePolicyAdjustmentTypePercentage GetAutoscaleConfigurationPolicyResponsePolicyResponsePolicyAdjustmentType = "percentage"
)

type GetAutoscaleConfigurationPolicyResponsePolicyResponsePolicy struct {
	// The name of the policy.
	Name *string `json:"name,omitempty"`
	// If the policy is enabled.
	Enabled *bool `json:"enabled,omitempty"`
	// Metric to use for the step policy.
	Metric *GetAutoscaleConfigurationPolicyResponsePolicyResponsePolicyMetric `json:"metric,omitempty"`
	// The type of adjustment to be made in the step policy.
	AdjustmentType *GetAutoscaleConfigurationPolicyResponsePolicyResponsePolicyAdjustmentType `json:"adjustment_type,omitempty"`
	// The steps for the step policy.
	// Each step defines an adjustment value and optional bounds.
	Steps []AutoscalePolicyStep `json:"steps,omitempty"`
}
