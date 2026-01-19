// This file is auto-generated. DO NOT EDIT.
// SPDX-License-Identifier: BSD-3-Clause
// Copyright (c) 2025, Unikraft GmbH.
// Licensed under the BSD-3-Clause License (the "License").
// You may not use this file except in compliance with the License.

package platform

type ServiceGroupInstance struct {
	// The UUID of the instance.  This is a unique identifier for the instance
	// that is generated when the instance is created.  The UUID is used to
	// reference the instance in API calls and can be used to identify the
	// instance in all API calls that require an instance identifier.
	Uuid *string `json:"uuid,omitempty"`
	// The name of the instance.  This is a human-readable name that can be used
	// to identify the instance.  The name must be unique within the context of
	// your account.  If no name is specified, a random name is generated for
	// you.  The name can also be used to identify the instance in API calls.
	Name *string `json:"name,omitempty"`
}
