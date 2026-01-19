// SPDX-License-Identifier: BSD-3-Clause
// Copyright (c) 2023, Unikraft GmbH.
// Licensed under the BSD-3-Clause License (the "License").
// You may not use this file except in compliance with the License.

package platform

type APIHTTPError int

const (
	APIHTTPErrorUnknownError APIHTTPError = iota
	APIHTTPErrorNotSupported
	APIHTTPErrorWrongMethod
	APIHTTPErrorNoAPIEndpoint
	APIHTTPErrorFailNotAll
	APIHTTPErrorTooMany
	APIHTTPErrorUnknown
	APIHTTPErrorInvalid
	APIHTTPErrorNotFound
	APIHTTPErrorNoFree
	APIHTTPErrorFailedOperation
	APIHTTPErrorFailedWrongVMState
	APIHTTPErrorTimedOut
	APIHTTPErrorMissingID
	APIHTTPErrorNotAllowed
	APIHTTPErrorAttached
	APIHTTPErrorMalformedRequest
	APIHTTPErrorNotBoth
	APIHTTPErrorAddService
	APIHTTPErrorQuota
	APIHTTPErrorGuestMemoryTooSmall
	APIHTTPErrorNotAttached
	APIHTTPErrorTooLong
	APIHTTPErrorAlreadyExists
	APIHTTPErrorCannotBeUsed
	APIHTTPErrorAutoscaleInstance
	APIHTTPErrorAutoscaleConfigured
	APIHTTPErrorAutoscaleNotConfigured
	APIHTTPErrorAutoscaleDisavled
	APIHTTPErrorAutoscaleSizeOOR
	APIHTTPErrorCertCNMisMatch
)
