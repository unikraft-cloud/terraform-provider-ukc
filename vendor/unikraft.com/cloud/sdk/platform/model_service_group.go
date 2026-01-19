// This file is auto-generated. DO NOT EDIT.
// SPDX-License-Identifier: BSD-3-Clause
// Copyright (c) 2025, Unikraft GmbH.
// Licensed under the BSD-3-Clause License (the "License").
// You may not use this file except in compliance with the License.

package platform

import "time"

// A service group on Unikraft Cloud is used to describe how your application
// exposes its functionality to the outside world.  Once defined, assigning an
// instance to the service will make it accessible from the Internet.
//
// An application, running as an instance, may expose one or more ports, e.g. it
// listens on port 80 because your application exposes a HTTP web service. This,
// along with a set of additional metadata defines how the "service" is
// configured and accessed.  For example, a service may be configured to use
// TLS, or be bound to a specific domain name.
//
// When an instance is assigned to a service group, it immediately becomes
// accessible over the Internet on the exposed public port, using the set DNS
// name, and is routed to the set destination port.
//
// Note: If you do not specify a DNS name when you create a service and you
// indicate that the application exposes some ports, Unikraft Cloud will
// generates a random DNS name for you.  Unikraft Cloud also supports custom
// domains like www.example.com and wildcard domains like *.example.com.

type ServiceGroup struct {
	// The UUID of the service group.
	//
	// This is a unique identifier for the service group that is generated when
	// the service group is created.  The UUID is used to reference the service in
	// API calls and can be used to identify the service group in all API calls
	// that require an identifier.
	Uuid *string `json:"uuid,omitempty"`
	// The name of the service group.
	//
	// This is a human-readable name that can be used to identify the service
	// group. The name must be unique within the context of your account.  The
	// name can also be used to identify the service in API calls.
	Name *string `json:"name,omitempty"`
	// The time the service was created.
	CreatedAt *time.Time `json:"created_at,omitempty"`
	// Indicates if the service will stay remain even after the last instance
	// detached.  If this is set to false, the service will be deleted when the
	// last instance detached from it.  If this is set to true, the service will
	// remain and can be reused by other instances.  This is useful if you want to
	// keep the service configuration, e.g., the published ports, handlers, and
	// domains, even if there are no instances assigned to it.
	Persistent *bool `json:"persistent,omitempty"`
	// Indicates if the service has autoscale enabled.  See the associated
	// autoscale documentation for more information about how to set this up.
	// Autoscale policies can be set up after the service has been created.
	Autoscale *bool `json:"autoscale,omitempty"`
	// The soft limit is used by the Unikraft Cloud load balancer to decide when
	// to wake up another standby instance.  For example, if the soft limit is set
	// to 5 and the service consists of 2 standby instances, one of the instances
	// receives up to 5 concurrent requests.  The 6th parallel requests wakes up
	// the second instance.  If there are no more standby instances to wake up,
	// the number of requests assigned to each instance will exceed the soft
	// limit.  The load balancer makes sure that when the number of in-flight
	// requests goes down again, instances are put into standby as fast as
	// possible.
	SoftLimit *uint64 `json:"soft_limit,omitempty"`
	// The hard limit defines the maximum number of concurrent requests that an
	// instance assigned to the this service can handle.  The load balancer will
	// never assign more requests to a single instance.  In case there are no
	// other instances available, excess requests fail (i.e., they are blocked and
	// not queued).
	HardLimit *uint64 `json:"hard_limit,omitempty"`
	// List of published network ports for this service and the destination port
	// to which Unikraft Cloud will forward traffic to.  Additional handlers can
	// be defined for each published port in order to define how the service will
	// handle incoming connections and forward traffic from the Internet to your
	// application.  For example, a service can be configured to terminate TLS
	// connections, redirect HTTP traffic, or enable HTTP mode for load balancing.
	Services []Service `json:"services,omitempty"`
	// List of domains associated with the service.  Domains are used to access
	// the service over the Internet.
	Domains []Domain `json:"domains,omitempty"`
	// List of instances assigned to the service.
	Instances []ServiceGroupInstance `json:"instances,omitempty"`
	// An optional field representing the status of the request.  This field is
	// only set when this message object is used as a response message.
	Status *ResponseStatus `json:"status,omitempty"`
	// An optional message providing additional information about the status.
	// This field is only set when this message object is used as a response
	// message, and is useful when the status is not `success`.
	Message *string `json:"message,omitempty"`
	// An optional error code providing additional information about the status.
	// This field is only set when this message object is used as a response
	// message, and is useful when the status is not `success`.
	Error *int32 `json:"error,omitempty"`
}
