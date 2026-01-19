// This file is auto-generated. DO NOT EDIT.
// SPDX-License-Identifier: BSD-3-Clause
// Copyright (c) 2025, Unikraft GmbH.
// Licensed under the BSD-3-Clause License (the "License").
// You may not use this file except in compliance with the License.

package platform

// The response data for this request.

type DeleteCertificatesResponseData struct {
	// The certificate(s) which were deleted by the request.
	Certificates []DeleteCertificatesResponseDeletedCertificate `json:"certificates,omitempty"`
}
