// This file is auto-generated. DO NOT EDIT.
// SPDX-License-Identifier: BSD-3-Clause
// Copyright (c) 2025, Unikraft GmbH.
// Licensed under the BSD-3-Clause License (the "License").
// You may not use this file except in compliance with the License.

package platform

// The request message for creating a new service group.

type CreateServiceGroupRequest struct {
	// Name of the service group.  This is a human-readable name that can be used
	// to identify the service group.  The name must be unique within the context
	// of your account.  If no name is specified, a random name is generated for
	// you.  The name can also be used to identify the service group in API calls.
	Name *string `json:"name,omitempty"`
	// Description of exposed services.
	Services []Service `json:"services,omitempty"`
	// Description of domains associated with the service group.
	Domains []CreateServiceGroupRequestDomain `json:"domains,omitempty"`
	// The soft limit is used by the Unikraft Cloud load balancer to decide when
	// to wake up another standby instance.
	//
	// For example, if the soft limit is set to 5 and the service consists of 2
	// standby instances, one of the instances receives up to 5 concurrent
	// requests.  The 6th parallel requests wakes up the second instance.  If
	// there are no more standby instances to wake up, the number of requests
	// assigned to each instance will exceed the soft limit.  The load balancer
	// makes sure that when the number of in-flight requests goes down again,
	// instances are put into standby as fast as possible.
	SoftLimit *uint64 `json:"soft_limit,omitempty"`
	// The hard limit defines the maximum number of concurrent requests that an
	// instance assigned to the this service can handle.
	//
	// The load balancer will never assign more requests to a single instance.  In
	// case there are no other instances available, excess requests fail (i.e.,
	// they are blocked and not queued).
	HardLimit *uint64 `json:"hard_limit,omitempty"`
}
