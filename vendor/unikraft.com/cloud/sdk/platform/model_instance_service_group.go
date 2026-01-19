// This file is auto-generated. DO NOT EDIT.
// SPDX-License-Identifier: BSD-3-Clause
// Copyright (c) 2025, Unikraft GmbH.
// Licensed under the BSD-3-Clause License (the "License").
// You may not use this file except in compliance with the License.

package platform

// The service group configuration for the instance.

type InstanceServiceGroup struct {
	// The UUID of the service group.
	//
	// This is a unique identifier for the service group that is generated when
	// the service is created.  The UUID is used to reference the service group
	// in API calls and can be used to identify the service in all API calls
	// that require an service identifier.
	Uuid *string `json:"uuid,omitempty"`
	// The name of the service group.
	//
	// This is a human-readable name that can be used to identify the service
	// group.  The name is unique within the context of your account.  The name
	// can also be used to identify the service group in API calls.
	Name *string `json:"name,omitempty"`
	// The domain configuration for the service group.
	Domains []ServiceGroupInstanceDomain `json:"domains,omitempty"`
}
