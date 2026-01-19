// This file is auto-generated. DO NOT EDIT.
// SPDX-License-Identifier: BSD-3-Clause
// Copyright (c) 2025, Unikraft GmbH.
// Licensed under the BSD-3-Clause License (the "License").
// You may not use this file except in compliance with the License.

package platform

// A domain name.
//
// Domain names are completely specified with all labels in the hierarchy of the
// DNS, having no parts omitted.  The domain can be associated with an existing
// certificate by specifying the certificate's name or UUID.  If no certificate
// is specified and a FQDN is provided, Unikraft Cloud will automatically
// generate a new certificate for the domain based on Let's Encrypt and seek to
// accomplish a DNS-01 challenge.

type Domain struct {
	// Publicly accessible domain name.  If this name ends in a period `.` it must
	// be a valid Full Qualified Domain Name (FQDN), otherwise it will become a
	// subdomain of the target metro.
	Fqdn *string `json:"fqdn,omitempty"`
	// Use an existing certificate for the domain.  If this field is
	// specified, the domain must be associated with a valid certificate.
	Certificate *Certificate `json:"certificate,omitempty"`
}
