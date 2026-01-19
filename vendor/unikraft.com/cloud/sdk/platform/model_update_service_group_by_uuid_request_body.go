// This file is auto-generated. DO NOT EDIT.
// SPDX-License-Identifier: BSD-3-Clause
// Copyright (c) 2025, Unikraft GmbH.
// Licensed under the BSD-3-Clause License (the "License").
// You may not use this file except in compliance with the License.

package platform

// The property to modify.
type UpdateServiceGroupByUUIDRequestBodyProp string

const (
	UpdateServiceGroupByUUIDRequestBodyPropServices   UpdateServiceGroupByUUIDRequestBodyProp = "services"
	UpdateServiceGroupByUUIDRequestBodyPropDomains    UpdateServiceGroupByUUIDRequestBodyProp = "domains"
	UpdateServiceGroupByUUIDRequestBodyPropSoft_limit UpdateServiceGroupByUUIDRequestBodyProp = "soft_limit"
	UpdateServiceGroupByUUIDRequestBodyPropHard_limit UpdateServiceGroupByUUIDRequestBodyProp = "hard_limit"
)

// The operation to perform.
type UpdateServiceGroupByUUIDRequestBodyOp string

const (
	UpdateServiceGroupByUUIDRequestBodyOpSet UpdateServiceGroupByUUIDRequestBodyOp = "set"
	UpdateServiceGroupByUUIDRequestBodyOpAdd UpdateServiceGroupByUUIDRequestBodyOp = "add"
	UpdateServiceGroupByUUIDRequestBodyOpDel UpdateServiceGroupByUUIDRequestBodyOp = "del"
)

type UpdateServiceGroupByUUIDRequestBody struct {
	// (Optional).  A client-provided identifier for tracking this operation in the response.
	Id *string `json:"id,omitempty"`
	// The property to modify.
	Prop UpdateServiceGroupByUUIDRequestBodyProp `json:"prop"`
	// The operation to perform.
	Op UpdateServiceGroupByUUIDRequestBodyOp `json:"op"`
	// The value for the update operation:
	// - For "services": array of Service objects
	// - For "domains": array of Domain objects
	// - For "soft_limit": integer (1–65535), must be <= "hard_limit"
	// - For "hard_limit": integer (1–65535), must be >= "soft_limit"
	Value *interface{} `json:"value,omitempty"`
}
