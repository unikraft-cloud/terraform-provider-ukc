// This file is auto-generated. DO NOT EDIT.
// SPDX-License-Identifier: BSD-3-Clause
// Copyright (c) 2025, Unikraft GmbH.
// Licensed under the BSD-3-Clause License (the "License").
// You may not use this file except in compliance with the License.

package platform

type AutoscalePolicyStep struct {
	// The adjustment value for the step.
	Adjustment *int64 `json:"adjustment,omitempty"`
	// Lower bound for the step.
	LowerBound *int64 `json:"lower_bound,omitempty"`
	// Upper bound for the step.
	UpperBound *int64 `json:"upper_bound,omitempty"`
}
