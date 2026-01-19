// This file is auto-generated. DO NOT EDIT.
// SPDX-License-Identifier: BSD-3-Clause
// Copyright (c) 2025, Unikraft GmbH.
// Licensed under the BSD-3-Clause License (the "License").
// You may not use this file except in compliance with the License.

package platform

// The instance template used for the autoscale configuration.
// Only if enabled is true.

type GetAutoscaleConfigurationsResponseServiceGroupTemplate struct {
	// The name of the template used for the autoscale configuration.
	Name *string `json:"name,omitempty"`
	// The UUID of the template used for the autoscale configuration.
	Uuid *string `json:"uuid,omitempty"`
}
