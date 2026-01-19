// This file is auto-generated. DO NOT EDIT.
// SPDX-License-Identifier: BSD-3-Clause
// Copyright (c) 2025, Unikraft GmbH.
// Licensed under the BSD-3-Clause License (the "License").
// You may not use this file except in compliance with the License.

package platform

// The request message to create an autoscale configuration policy for a
// service.

type CreateAutoscaleConfigurationPolicyRequest struct {
	// The Name of the service to add a policy to.
	Name string                                        `json:"name"`
	Type CreateAutoscaleConfigurationPolicyRequestType `json:"type"`
}
