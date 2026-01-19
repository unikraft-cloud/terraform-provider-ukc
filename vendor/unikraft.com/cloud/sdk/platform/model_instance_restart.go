// This file is auto-generated. DO NOT EDIT.
// SPDX-License-Identifier: BSD-3-Clause
// Copyright (c) 2025, Unikraft GmbH.
// Licensed under the BSD-3-Clause License (the "License").
// You may not use this file except in compliance with the License.

package platform

import "time"

// The current restart attempt for the instance.
// Not used for template instances.

type InstanceRestart struct {
	// Current restart attempt number. This is incremented each time the instance
	// is restarted automatically by the platform.
	Attempt *uint32 `json:"attempt,omitempty"`
	// Timestamp of the next scheduled restart attempt.
	NextAt *time.Time `json:"next_at,omitempty"`
}
