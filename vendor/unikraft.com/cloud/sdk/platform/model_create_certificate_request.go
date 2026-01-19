// This file is auto-generated. DO NOT EDIT.
// SPDX-License-Identifier: BSD-3-Clause
// Copyright (c) 2025, Unikraft GmbH.
// Licensed under the BSD-3-Clause License (the "License").
// You may not use this file except in compliance with the License.

package platform

// The request message for creating/uploading a new certificate.

type CreateCertificateRequest struct {
	// The common name (CN) of the certificate.
	Cn string `json:"cn"`
	// The chain of the certificate.
	Chain string `json:"chain"`
	// The private key of the certificate.
	Pkey string `json:"pkey"`
	// The name of the certificate.
	//
	// This is a human-readable name that can be used to identify the certificate.
	// The name must be unique within the context of your account.  If no name is
	// specified, a random name is generated for you.  The name can also be used
	// to identify the certificate in API calls.
	Name *string `json:"name,omitempty"`
}
