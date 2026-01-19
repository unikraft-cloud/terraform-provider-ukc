// This file is auto-generated. DO NOT EDIT.
// SPDX-License-Identifier: BSD-3-Clause
// Copyright (c) 2025, Unikraft GmbH.
// Licensed under the BSD-3-Clause License (the "License").
// You may not use this file except in compliance with the License.

package platform

// (Optional).  The service group configuration when creating an instance.
//
// When creating an instance, either a previously created (persistent) service
// group can be referenced (either through its name or UUID), or a new
// (ephemeral) service group can be created for the instance by specifying the
// list of services it should expose and optionally the domains it should use.
// Not used by template instances.

type CreateInstanceRequestServiceGroup struct {
	// (Optional).  Reference an existing (persistent) service group by its
	// UUID.  Mutually exclusive with name.
	Uuid *string `json:"uuid,omitempty"`
	// (Optional).  Reference an existing (persistent) service group by its
	// name.  Mutually exclusive with UUID.
	Name *string `json:"name,omitempty"`
	// If no existing (persistent) service group is specified via its
	// identifier, a new (ephemeral) service group can be created.  In addition
	// to the services it must expose, you can specify which domains it should
	// use too.
	Domains []CreateInstanceRequestDomain `json:"domains,omitempty"`
	// If no existing service group identifier is provided, one or more new
	// (ephemeral, non-persistent) service(s) can be created with the following
	// definitions.
	Services []Service `json:"services,omitempty"`
	// The soft limit for the number of services that can be created in this
	// service group.
	SoftLimit *uint32 `json:"soft_limit,omitempty"`
	// The hard limit for the number of services that can be created in this
	// service group.
	HardLimit *uint32 `json:"hard_limit,omitempty"`
}
