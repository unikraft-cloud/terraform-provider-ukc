// This file is auto-generated. DO NOT EDIT.
// SPDX-License-Identifier: BSD-3-Clause
// Copyright (c) 2025, Unikraft GmbH.
// Licensed under the BSD-3-Clause License (the "License").
// You may not use this file except in compliance with the License.

package platform

// The request message for converting one or more instances to templates.

type CreateTemplateInstancesRequest struct {
	// The list of IDs of the instances to convert to templates.
	Ids []NameOrUUID `json:"ids"`
}
