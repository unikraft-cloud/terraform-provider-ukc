// This file is auto-generated. DO NOT EDIT.
// SPDX-License-Identifier: BSD-3-Clause
// Copyright (c) 2025, Unikraft GmbH.
// Licensed under the BSD-3-Clause License (the "License").
// You may not use this file except in compliance with the License.

package platform

// A single update operation to be applied to a service group
// The property to modify.
type UpdateServiceGroupsRequestItemProp string

const (
	UpdateServiceGroupsRequestItemPropServices   UpdateServiceGroupsRequestItemProp = "services"
	UpdateServiceGroupsRequestItemPropDomains    UpdateServiceGroupsRequestItemProp = "domains"
	UpdateServiceGroupsRequestItemPropSoft_limit UpdateServiceGroupsRequestItemProp = "soft_limit"
	UpdateServiceGroupsRequestItemPropHard_limit UpdateServiceGroupsRequestItemProp = "hard_limit"
)

// The operation to perform.
type UpdateServiceGroupsRequestItemOp string

const (
	UpdateServiceGroupsRequestItemOpSet UpdateServiceGroupsRequestItemOp = "set"
	UpdateServiceGroupsRequestItemOpAdd UpdateServiceGroupsRequestItemOp = "add"
	UpdateServiceGroupsRequestItemOpDel UpdateServiceGroupsRequestItemOp = "del"
)

type UpdateServiceGroupsRequestItem struct {
	// (Optional).  A client-provided identifier for tracking this operation in the response.
	Id *string `json:"id,omitempty"`
	// The UUID of the service group to update.  Mutually exclusive with name.
	Uuid *string `json:"uuid,omitempty"`
	// The name of the service group to update.  Mutually exclusive with UUID.
	Name *string `json:"name,omitempty"`
	// The property to modify.
	Prop UpdateServiceGroupsRequestItemProp `json:"prop"`
	// The operation to perform.
	Op UpdateServiceGroupsRequestItemOp `json:"op"`
	// The value for the update operation:
	// - For "services": array of Service objects
	// - For "domains": array of Domain objects
	// - For "soft_limit": integer (1–65535), must be <= "hard_limit"
	// - For "hard_limit": integer (1–65535), must be >= "soft_limit"
	Value *interface{} `json:"value,omitempty"`
}
