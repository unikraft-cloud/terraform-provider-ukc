// This file is auto-generated. DO NOT EDIT.
// SPDX-License-Identifier: BSD-3-Clause
// Copyright (c) 2025, Unikraft GmbH.
// Licensed under the BSD-3-Clause License (the "License").
// You may not use this file except in compliance with the License.

package platform

// The response data for this request.

type GetInstancesMetricsResponseData struct {
	// The instance which this requested metrics for.  Note: only one instance
	// can be specified in the request, so this will always contain a single
	// entry.
	Instances []GetInstancesMetricsResponseInstanceMetrics `json:"instances,omitempty"`
}
