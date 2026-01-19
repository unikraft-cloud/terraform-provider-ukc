// This file is auto-generated. DO NOT EDIT.
// SPDX-License-Identifier: BSD-3-Clause
// Copyright (c) 2025, Unikraft GmbH.
// Licensed under the BSD-3-Clause License (the "License").
// You may not use this file except in compliance with the License.

package platform

import "time"

type Image struct {
	// The canonical name of the image is known as the "tag".
	Tag *string `json:"tag,omitempty"`
	// The digest of the image is a unique identifier of the image manifest which
	// is a string representation including the hashing algorithm and the hash
	// value separated by a colon.
	Digest *string `json:"digest,omitempty"`
	// A description of the image.
	Description *string `json:"description,omitempty"`
	// When the image was created.
	CreatedAt *time.Time `json:"created_at,omitempty"`
	// The architecture of the image.
	Arch *string `json:"arch,omitempty"`
	// The entrypoint of the image is the command that is run when the image is
	// started.
	Entrypoint []string `json:"entrypoint,omitempty"`
	// The command to run when the image is started.
	Cmd []string `json:"cmd,omitempty"`
	// The environment variables to set when the image is started.
	Env []string `json:"env,omitempty"`
	// Documented port mappings for the image.
	Ports []string `json:"ports,omitempty"`
	// Documented volumes for the image.
	Volumes []string `json:"volumes,omitempty"`
	// Labels are key-value pairs.
	Labels map[string]string `json:"labels,omitempty"`
	// The working directory for the image is the directory that is set as the
	// current working directory when the image is started.
	Workdir *string      `json:"workdir,omitempty"`
	Kernel  *ImageKernel `json:"kernel,omitempty"`
	// List of auxiliary ROMs that are used by the image.
	AuxiliaryRoms []Object `json:"auxiliary_roms,omitempty"`
}
