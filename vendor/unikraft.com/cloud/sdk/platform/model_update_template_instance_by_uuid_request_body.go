// This file is auto-generated. DO NOT EDIT.
// SPDX-License-Identifier: BSD-3-Clause
// Copyright (c) 2025, Unikraft GmbH.
// Licensed under the BSD-3-Clause License (the "License").
// You may not use this file except in compliance with the License.

package platform

// The property to modify.
type UpdateTemplateInstanceByUUIDRequestBodyProp string

const (
	UpdateTemplateInstanceByUUIDRequestBodyPropTags        UpdateTemplateInstanceByUUIDRequestBodyProp = "tags"
	UpdateTemplateInstanceByUUIDRequestBodyPropDelete_lock UpdateTemplateInstanceByUUIDRequestBodyProp = "delete_lock"
)

// The operation to perform on the property.
type UpdateTemplateInstanceByUUIDRequestBodyOp string

const (
	UpdateTemplateInstanceByUUIDRequestBodyOpSet UpdateTemplateInstanceByUUIDRequestBodyOp = "set"
	UpdateTemplateInstanceByUUIDRequestBodyOpAdd UpdateTemplateInstanceByUUIDRequestBodyOp = "add"
	UpdateTemplateInstanceByUUIDRequestBodyOpDel UpdateTemplateInstanceByUUIDRequestBodyOp = "del"
)

type UpdateTemplateInstanceByUUIDRequestBody struct {
	// (Optional).  A client-provided identifier for tracking this operation in
	// the response.
	Id *string `json:"id,omitempty"`
	// The property to modify.
	Prop UpdateTemplateInstanceByUUIDRequestBodyProp `json:"prop"`
	// The operation to perform on the property.
	Op UpdateTemplateInstanceByUUIDRequestBodyOp `json:"op"`
	// The value for the update operation. The type depends on the property and operation:
	// - For "tags": array of strings
	// - For "delete_lock": boolean
	Value *interface{} `json:"value,omitempty"`
}
