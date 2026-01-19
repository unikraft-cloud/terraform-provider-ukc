// This file is auto-generated. DO NOT EDIT.
// SPDX-License-Identifier: BSD-3-Clause
// Copyright (c) 2025, Unikraft GmbH.
// Licensed under the BSD-3-Clause License (the "License").
// You may not use this file except in compliance with the License.

package platform

// A single update operation to be applied to a template instance.
// The property to modify.
type UpdateTemplateInstancesRequestItemProp string

const (
	UpdateTemplateInstancesRequestItemPropTags        UpdateTemplateInstancesRequestItemProp = "tags"
	UpdateTemplateInstancesRequestItemPropDelete_lock UpdateTemplateInstancesRequestItemProp = "delete_lock"
)

// The operation to perform on the property.
type UpdateTemplateInstancesRequestItemOp string

const (
	UpdateTemplateInstancesRequestItemOpSet UpdateTemplateInstancesRequestItemOp = "set"
	UpdateTemplateInstancesRequestItemOpAdd UpdateTemplateInstancesRequestItemOp = "add"
	UpdateTemplateInstancesRequestItemOpDel UpdateTemplateInstancesRequestItemOp = "del"
)

type UpdateTemplateInstancesRequestItem struct {
	// (Optional).  A client-provided identifier for tracking this operation in
	// the response.
	Id *string `json:"id,omitempty"`
	// The UUID of the template instance to update. Mutually exclusive with name.
	Uuid *string `json:"uuid,omitempty"`
	// The name of the template instance to update. Mutually exclusive with UUID.
	Name *string `json:"name,omitempty"`
	// The property to modify.
	Prop UpdateTemplateInstancesRequestItemProp `json:"prop"`
	// The operation to perform on the property.
	Op UpdateTemplateInstancesRequestItemOp `json:"op"`
	// The value for the update operation. The type depends on the property and operation:
	// - For "tags": array of strings
	// - For "delete_lock": boolean
	Value *interface{} `json:"value,omitempty"`
}
