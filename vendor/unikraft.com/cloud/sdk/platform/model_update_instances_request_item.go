// This file is auto-generated. DO NOT EDIT.
// SPDX-License-Identifier: BSD-3-Clause
// Copyright (c) 2025, Unikraft GmbH.
// Licensed under the BSD-3-Clause License (the "License").
// You may not use this file except in compliance with the License.

package platform

// A single update operation to be applied to an instance.
// The property to modify.
type UpdateInstancesRequestItemProp string

const (
	UpdateInstancesRequestItemPropImage         UpdateInstancesRequestItemProp = "image"
	UpdateInstancesRequestItemPropArgs          UpdateInstancesRequestItemProp = "args"
	UpdateInstancesRequestItemPropEnv           UpdateInstancesRequestItemProp = "env"
	UpdateInstancesRequestItemPropMemory_mb     UpdateInstancesRequestItemProp = "memory_mb"
	UpdateInstancesRequestItemPropVcpus         UpdateInstancesRequestItemProp = "vcpus"
	UpdateInstancesRequestItemPropScale_to_zero UpdateInstancesRequestItemProp = "scale_to_zero"
	UpdateInstancesRequestItemPropTags          UpdateInstancesRequestItemProp = "tags"
	UpdateInstancesRequestItemPropDelete_lock   UpdateInstancesRequestItemProp = "delete_lock"
)

// The operation to perform on the property.
type UpdateInstancesRequestItemOp string

const (
	UpdateInstancesRequestItemOpSet UpdateInstancesRequestItemOp = "set"
	UpdateInstancesRequestItemOpAdd UpdateInstancesRequestItemOp = "add"
	UpdateInstancesRequestItemOpDel UpdateInstancesRequestItemOp = "del"
)

type UpdateInstancesRequestItem struct {
	// (Optional).  A client-provided identifier for tracking this operation in
	// the response.
	Id *string `json:"id,omitempty"`
	// The UUID of the instance to update. Mutually exclusive with name.
	Uuid *string `json:"uuid,omitempty"`
	// The name of the instance to update. Mutually exclusive with UUID.
	Name *string `json:"name,omitempty"`
	// The property to modify.
	Prop UpdateInstancesRequestItemProp `json:"prop"`
	// The operation to perform on the property.
	Op UpdateInstancesRequestItemOp `json:"op"`
	// The value for the update operation. The type depends on the property and operation:
	// - For "image": string
	// - For "args": string or array of strings
	// - For "env": object (for SET/ADD) or string/array of strings (for DEL)
	// - For "memory_mb": integer
	// - For "vcpus": integer
	// - For "scale_to_zero": object with cooldown_time_ms, policy, and stateful fields
	// - For "tags": array of strings
	// - For "delete_lock": boolean
	Value *interface{} `json:"value,omitempty"`
}
