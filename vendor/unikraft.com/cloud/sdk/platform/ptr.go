// SPDX-License-Identifier: BSD-3-Clause
// Copyright (c) 2025, Unikraft GmbH.
// Licensed under the BSD-3-Clause License (the "License").
// You may not use this file except in compliance with the License.

package platform

// Ptr is a utility function that returns a pointer to the given value.
func Ptr[T any](v T) *T {
	return &v
}
