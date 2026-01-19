// This file is auto-generated. DO NOT EDIT.
// SPDX-License-Identifier: BSD-3-Clause
// Copyright (c) 2025, Unikraft GmbH.
// Licensed under the BSD-3-Clause License (the "License").
// You may not use this file except in compliance with the License.

package platform

// The property to modify.
type UpdateInstanceByUUIDRequestBodyProp string

const (
	UpdateInstanceByUUIDRequestBodyPropImage         UpdateInstanceByUUIDRequestBodyProp = "image"
	UpdateInstanceByUUIDRequestBodyPropArgs          UpdateInstanceByUUIDRequestBodyProp = "args"
	UpdateInstanceByUUIDRequestBodyPropEnv           UpdateInstanceByUUIDRequestBodyProp = "env"
	UpdateInstanceByUUIDRequestBodyPropMemory_mb     UpdateInstanceByUUIDRequestBodyProp = "memory_mb"
	UpdateInstanceByUUIDRequestBodyPropVcpus         UpdateInstanceByUUIDRequestBodyProp = "vcpus"
	UpdateInstanceByUUIDRequestBodyPropScale_to_zero UpdateInstanceByUUIDRequestBodyProp = "scale_to_zero"
	UpdateInstanceByUUIDRequestBodyPropTags          UpdateInstanceByUUIDRequestBodyProp = "tags"
	UpdateInstanceByUUIDRequestBodyPropDelete_lock   UpdateInstanceByUUIDRequestBodyProp = "delete_lock"
)

// The operation to perform on the property.
type UpdateInstanceByUUIDRequestBodyOp string

const (
	UpdateInstanceByUUIDRequestBodyOpSet UpdateInstanceByUUIDRequestBodyOp = "set"
	UpdateInstanceByUUIDRequestBodyOpAdd UpdateInstanceByUUIDRequestBodyOp = "add"
	UpdateInstanceByUUIDRequestBodyOpDel UpdateInstanceByUUIDRequestBodyOp = "del"
)

type UpdateInstanceByUUIDRequestBody struct {
	// (Optional).  A client-provided identifier for tracking this operation in
	// the response.
	Id *string `json:"id,omitempty"`
	// The property to modify.
	Prop UpdateInstanceByUUIDRequestBodyProp `json:"prop"`
	// The operation to perform on the property.
	Op UpdateInstanceByUUIDRequestBodyOp `json:"op"`
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
