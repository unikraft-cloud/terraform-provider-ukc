// This file is auto-generated. DO NOT EDIT.
// SPDX-License-Identifier: BSD-3-Clause
// Copyright (c) 2025, Unikraft GmbH.
// Licensed under the BSD-3-Clause License (the "License").
// You may not use this file except in compliance with the License.

package platform

// The request message for creating a new instance.
// Restart policy for the instance.  This defines how the instance should
// behave when it stops or crashes.
type CreateInstanceRequestRestartPolicy string

const (
	CreateInstanceRequestRestartPolicyNever      CreateInstanceRequestRestartPolicy = "never"
	CreateInstanceRequestRestartPolicyAlways     CreateInstanceRequestRestartPolicy = "always"
	CreateInstanceRequestRestartPolicyOn_failure CreateInstanceRequestRestartPolicy = "on_failure"
)

// Features to enable for the instance.  Features are specific
// configurations or capabilities that can be enabled for the instance.
type CreateInstanceRequestFeatures string

const (
	CreateInstanceRequestFeaturesDelete_on_stop CreateInstanceRequestFeatures = "delete_on_stop"
)

type CreateInstanceRequest struct {
	// (Optional).  The name of the instance.
	//
	// If not provided, a random name will be generated.  The name must be unique.
	Name *string `json:"name,omitempty"`
	// The image to use for the instance.
	Image string `json:"image"`
	// (Optional).  The arguments to pass to the instance when it starts.
	Args []string `json:"args,omitempty"`
	// (Optional).  Environment variables to set for the instance.
	Env map[string]string `json:"env,omitempty"`
	// (Optional).  Memory in MB to allocate for the instance.  Default is 128.
	MemoryMb     *int64                             `json:"memory_mb,omitempty"`
	ServiceGroup *CreateInstanceRequestServiceGroup `json:"service_group,omitempty"`
	// Volumes to attach to the instance.
	//
	// This list can contain both existing and new volumes to create as part of
	// the instance creation.  Existing volumes can be referenced by their name or
	// UUID.  New volumes can be created by specifying a name, size in MiB, and
	// mount point in the instance.  The mount point is the directory in the
	// instance where the volume will be mounted.
	Volumes []CreateInstanceRequestVolume `json:"volumes,omitempty"`
	// Whether the instance should start automatically on creation.
	Autostart *bool `json:"autostart,omitempty"`
	// Number of replicas for the instance.
	Replicas *int64 `json:"replicas,omitempty"`
	// Restart policy for the instance.  This defines how the instance should
	// behave when it stops or crashes.
	RestartPolicy *CreateInstanceRequestRestartPolicy `json:"restart_policy,omitempty"`
	ScaleToZero   *CreateInstanceRequestScaleToZero   `json:"scale_to_zero,omitempty"`
	// Number of vCPUs to allocate for the instance.
	Vcpus *int32 `json:"vcpus,omitempty"`
	// Timeout to wait for all new instances to reach running state in
	// milliseconds.  If you autostart your new instance, you can wait for it to
	// finish starting with a blocking API call if you specify a wait timeout
	// greater than zero.  No wait performed for a value of 0.
	WaitTimeoutMs *int64 `json:"wait_timeout_ms,omitempty"`
	// Features to enable for the instance.  Features are specific
	// configurations or capabilities that can be enabled for the instance.
	Features []CreateInstanceRequestFeatures `json:"features,omitempty"`
}
