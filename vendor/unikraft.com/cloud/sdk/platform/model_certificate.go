// This file is auto-generated. DO NOT EDIT.
// SPDX-License-Identifier: BSD-3-Clause
// Copyright (c) 2025, Unikraft GmbH.
// Licensed under the BSD-3-Clause License (the "License").
// You may not use this file except in compliance with the License.

package platform

import "time"

// The current state of the certificate.
//
// This indicates whether the certificate is pending issuance, valid and
// ready for use, or in an error state. See CertificateState enum for
// detailed state descriptions.
type CertificateState string

const (
	CertificateStatePending CertificateState = "pending"
	CertificateStateValid   CertificateState = "valid"
	CertificateStateError   CertificateState = "error"
)

type Certificate struct {
	// The UUID of the certificate.
	//
	// This is a unique identifier for the certificate that is generated when the
	// certificate is created.  The UUID is used to reference the certificate in
	// API calls and can be used to identify the certificate in all API calls that
	// require an identifier.
	Uuid *string `json:"uuid,omitempty"`
	// The name of the certificate.
	//
	// This is a human-readable name that can be used to identify the certificate.
	// The name must be unique within the context of your account.  The name can
	// also be used to identify the certificate in API calls.
	Name *string `json:"name,omitempty"`
	// The time the certificate was created.
	CreatedAt *time.Time `json:"created_at,omitempty"`
	// The common name (CN) field from the certificate's subject.
	//
	// This is typically the primary domain name that the certificate is issued
	// for. It represents the main identity that the certificate validates.
	CommonName *string `json:"common_name,omitempty"`
	// The complete subject distinguished name (DN) of the certificate.
	//
	// This contains the full subject information from the certificate, including
	// the common name, organization, organizational unit, locality, state, and
	// country. The subject identifies the entity that the certificate is issued to.
	Subject *string `json:"subject,omitempty"`
	// The complete issuer distinguished name (DN) of the certificate.
	//
	// This identifies the Certificate Authority (CA) that issued the certificate.
	// It contains information about the CA including its common name, organization,
	// and country.
	Issuer *string `json:"issuer,omitempty"`
	// The unique serial number assigned to the certificate by the issuing CA.
	//
	// This is a unique identifier within the scope of the issuing CA that can be
	// used to identify and track the certificate. Serial numbers are typically
	// represented as hexadecimal strings.
	SerialNumber *string `json:"serial_number,omitempty"`
	// The date and time when the certificate becomes valid.
	//
	// The certificate should not be trusted before this date. This timestamp
	// marks the beginning of the certificate's validity period.
	NotBefore *time.Time `json:"not_before,omitempty"`
	// The date and time when the certificate expires.
	//
	// The certificate should not be trusted after this date. This timestamp
	// marks the end of the certificate's validity period. Certificates should
	// be renewed before this date to maintain service availability.
	NotAfter *time.Time `json:"not_after,omitempty"`
	// The current state of the certificate.
	//
	// This indicates whether the certificate is pending issuance, valid and
	// ready for use, or in an error state. See CertificateState enum for
	// detailed state descriptions.
	State *CertificateState `json:"state,omitempty"`
	// An optional field representing the status of the request.  This field is
	// only set when this message object is used as a response message.
	Status *ResponseStatus `json:"status,omitempty"`
	// An optional message providing additional information about the status.
	// This field is only set when this message object is used as a response
	// message, and is useful when the status is not `success`.
	Message *string `json:"message,omitempty"`
	// An optional error code providing additional information about the status.
	// This field is only set when this message object is used as a response
	// message, and is useful when the status is not `success`.
	Error *int32 `json:"error,omitempty"`
}
