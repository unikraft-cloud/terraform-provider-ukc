// This file is auto-generated. DO NOT EDIT.
// SPDX-License-Identifier: BSD-3-Clause
// Copyright (c) 2025, Unikraft GmbH.
// Licensed under the BSD-3-Clause License (the "License").
// You may not use this file except in compliance with the License.

package platform

type GetInstanceLogsByUUIDRequestBody struct {
	// The byte offset of the log output to receive.  A negative sign makes the
	// offset relative to the end of the log.
	Offset *uint64 `json:"offset,omitempty"`
	// The amount of bytes to return at most.
	Limit *int64 `json:"limit,omitempty"`
}
