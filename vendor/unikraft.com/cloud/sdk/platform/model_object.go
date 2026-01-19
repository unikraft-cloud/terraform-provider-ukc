// This file is auto-generated. DO NOT EDIT.
// SPDX-License-Identifier: BSD-3-Clause
// Copyright (c) 2025, Unikraft GmbH.
// Licensed under the BSD-3-Clause License (the "License").
// You may not use this file except in compliance with the License.

package platform

// An object is a single component of an image which is external and can be
// uniquely identified by its digest.

type Object struct {
	// The digest is a string representation including the hashing
	// algorithm and the hash value separated by a colon.
	Digest *string `json:"digest,omitempty"`
	// The media type of the layer is a string that identifies the type of
	// content that the layer contains.
	MediaType *string `json:"media_type,omitempty"`
	// The size of the layer in bytes.
	Size *int64 `json:"size,omitempty"`
}
