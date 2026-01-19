// This file is auto-generated. DO NOT EDIT.
// SPDX-License-Identifier: BSD-3-Clause
// Copyright (c) 2025, Unikraft GmbH.
// Licensed under the BSD-3-Clause License (the "License").
// You may not use this file except in compliance with the License.

package platform

// A service connects a public-facing port to an internal destination port on
// which an application instance listens on.  Additional handlers can be defined
// for each published port in order to define how the service will handle
// incoming connections and forward traffic from the Internet to your
// application.  For example, a service can be configured to terminate TLS
// connections, redirect HTTP traffic, or enable HTTP mode for load balancing.
// Connection handlers to use for the service.  Handlers define how the
// service will handle incoming connections and forward traffic from the
// Internet to your application.  For example, a service can be configured
// to terminate TLS connections, redirect HTTP traffic, or enable HTTP mode
// for load balancing.  You configure the handlers for every published
// service port individually.
type ServiceHandlers string

const (
	ServiceHandlersTls      ServiceHandlers = "tls"
	ServiceHandlersHttp     ServiceHandlers = "http"
	ServiceHandlersRedirect ServiceHandlers = "redirect"
)

type Service struct {
	// This is the public-facing port that the service will be accessible from
	// on the Internet.
	Port uint32 `json:"port"`
	// The port number that the instance is listening on.  This is the internal
	// port which Unikraft Cloud will forward traffic to.
	DestinationPort *uint32 `json:"destination_port,omitempty"`
	// Connection handlers to use for the service.  Handlers define how the
	// service will handle incoming connections and forward traffic from the
	// Internet to your application.  For example, a service can be configured
	// to terminate TLS connections, redirect HTTP traffic, or enable HTTP mode
	// for load balancing.  You configure the handlers for every published
	// service port individually.
	Handlers []ServiceHandlers `json:"handlers,omitempty"`
}
