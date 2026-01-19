// This file is auto-generated. DO NOT EDIT.
// SPDX-License-Identifier: BSD-3-Clause
// Copyright (c) 2025, Unikraft GmbH.
// Licensed under the BSD-3-Clause License (the "License").
// You may not use this file except in compliance with the License.

package platform

// The domain configuration for the service group.
//
// A domain defines a publicly accessible domain name for the instance.  If
// the domain name ends with a period `.`, it must be a valid Fully Qualified
// Domain Name (FQDN), otherwise it will become a subdomain of the target
// metro.  The domain can be associated with an existing certificate by
// specifying the certificate's name or UUID.  If no certificate is specified
// and a FQDN is provided, Unikraft Cloud will automatically generate a new
// certificate for the domain based on Let's Encrypt and seek to accomplish a
// DNS-01 challenge.

type CreateInstanceRequestDomain struct {
	// Publicly accessible domain name.
	//
	// If this name ends in a period `.` it must be a valid Full Qualified
	// Domain Name (FQDN), e.g. `example.com.`; otherwise it will become a
	// subdomain of the target metro, e.g. `example` becomes
	// `example.fra0.unikraft.app`.
	Name string `json:"name"`
	// A reference to an existing certificate which can be used for the
	// specified domain.  If unspecified, Unikraft Cloud will
	// automatically generate a new certificate for the domain based on Let's
	// Encrypt and seek to accomplish a DNS-01 challenge.
	Certificate *NameOrUUID `json:"certificate,omitempty"`
}
