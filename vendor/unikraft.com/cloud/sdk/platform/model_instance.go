// This file is auto-generated. DO NOT EDIT.
// SPDX-License-Identifier: BSD-3-Clause
// Copyright (c) 2025, Unikraft GmbH.
// Licensed under the BSD-3-Clause License (the "License").
// You may not use this file except in compliance with the License.

package platform

import "time"

// An instance is a unikernel virtual machine running an application.
// The state of the instance.  This indicates the current state of the
// instance, such as whether it is running, stopped, or in an error state.
type InstanceState string

const (
	InstanceStateStopped  InstanceState = "stopped"
	InstanceStateStarting InstanceState = "starting"
	InstanceStateRunning  InstanceState = "running"
	InstanceStateDraining InstanceState = "draining"
	InstanceStateStopping InstanceState = "stopping"
	InstanceStateTemplate InstanceState = "template"
	InstanceStateStandby  InstanceState = "standby"
)

// The restart configuration for the instance.
//
// When an instance stops either because the application exits or the instance
// crashes, Unikraft Cloud can auto-restart your instance.  Auto-restarts are
// performed according to the restart policy configured for a particular
// instance.
//
// The policy can have the following values:
//
// | Policy       | Description |
// |--------------|-------------|
// | `never`      | Never restart the instance (default). |
// | `always`     | Always restart the instance when the stop is initiated from within the instance (i.e., the application exits or the instance crashes). |
// | `on-failure` | Only restart the instance if it crashes. |
//
// When an instance stops, the stop reason and the configured restart policy
// are evaluated to decide if a restart should be performed.  Unikraft Cloud
// uses an exponential back-off delay (immediate, 5s, 10s, 20s, 40s, ..., 5m)
// to slow down restarts in tight crash loops.  If an instance runs without
// problems for 10s the back-off delay is reset and the restart sequence ends.
//
// The `restart.attempt` attribute reported in counts the number of restarts
// performed in the current sequence.  The `restart.next_at` field indicates
// when the next restart will take place if a back-off delay is in effect.
//
// A manual start or stop of the instance aborts the restart sequence and
// resets the back-off delay.
type InstanceRestartPolicy string

const (
	InstanceRestartPolicyNever      InstanceRestartPolicy = "never"
	InstanceRestartPolicyAlways     InstanceRestartPolicy = "always"
	InstanceRestartPolicyOn_failure InstanceRestartPolicy = "on_failure"
)

type Instance struct {
	// The UUID of the instance.
	//
	// This is a unique identifier for the instance that is generated when the
	// instance is created.  The UUID is used to reference the instance in API
	// calls and can be used to identify the instance in all API calls that
	// require an instance identifier.
	Uuid *string `json:"uuid,omitempty"`
	// The name of the instance.
	//
	// This is a human-readable name that can be used to identify the instance.
	// The name must be unique within the context of your account.  The name can
	// also be used to identify the instance in API calls.
	Name *string `json:"name,omitempty"`
	// The time the instance was created.
	CreatedAt *time.Time `json:"created_at,omitempty"`
	// The state of the instance.  This indicates the current state of the
	// instance, such as whether it is running, stopped, or in an error state.
	State *InstanceState `json:"state,omitempty"`
	// The internal hostname of the instance.  This address can be used privately
	// within the Unikraft Cloud network to access the instance.  It is not
	// accessible from the public Internet.
	PrivateFqdn *string `json:"private_fqdn,omitempty"`
	// The image used to create the instance.  This is a reference to the
	// Unikraft image that was used to create the instance.
	Image *string `json:"image,omitempty"`
	// The amount of memory in megabytes allocated for the instance.  This is the
	// total amount of memory that is available to the instance for its
	// operations.
	MemoryMb *uint64 `json:"memory_mb,omitempty"`
	// The number of vCPUs allocated for the instance.  This is the total
	// number of virtual CPUs that are available to the instance for its
	// operations.
	Vcpus *uint32 `json:"vcpus,omitempty"`
	// The arguments passed to the instance when it was started.  This is a
	// list of command-line arguments that were provided to the instance at
	// startup.  These arguments can be used to configure the behavior of the
	// instance and its applications.
	Args []string `json:"args,omitempty"`
	// Environment variables set for the instance.
	Env map[string]string `json:"env,omitempty"`
	// The total number of times the instance has been started.  This is a counter
	// that increments each time the instance is started, regardless of whether it
	// was manually stopped or restarted.  This can be useful for tracking the
	// usage of the instance over time and/or for debugging purposes.
	//
	// Not used for template instances.
	StartCount *uint64 `json:"start_count,omitempty"`
	// The total number of times the instance has been restarted. This is a counter
	// that increments each time the instance has been restarted. This can be
	// useful for tracking the usage of the instance over time and/or for
	// debugging purposes.
	// Not used for template instances.
	RestartCount *uint64 `json:"restart_count,omitempty"`
	// The time the instance was started.  This is the timestamp when the
	// instance was last started.
	// Not used for template instances.
	StartedAt *time.Time `json:"started_at,omitempty"`
	// The time the instance was stopped.  This is the timestamp when the
	// instance was last stopped.  If the instance is currently running, this
	// field will be empty.
	// Not used for template instances.
	StoppedAt *time.Time `json:"stopped_at,omitempty"`
	// The total amount of time the instance has been running in milliseconds.
	// Not used for template instances.
	UptimeMs *uint64 `json:"uptime_ms,omitempty"`
	// (Developer-only).  The time taken between the main controller and the
	// beginning of execution of the VMM (Virtual Machine Monitor) measured in
	// microseconds.  This field is primarily used for debugging and performance
	// analysis purposes.
	// Not used for template instances.
	VmmStartTimeUs *uint64 `json:"vmm_start_time_us,omitempty"`
	// (Developer-only).  The time it took the VMM (Virtual Machine Monitor) to
	// load the instance's kernel and initramfs into VM memory measured in
	// microseconds.  This field is primarily used for debugging and performance
	// analysis purposes.
	// Not used for template instances.
	VmmLoadTimeUs *uint64 `json:"vmm_load_time_us,omitempty"`
	// (Developer-only).  The time taken for the VMM (Virtual Machine Monitor) to
	// become ready to execute the instance measured in microseconds.  This is the
	// time from when the VMM started until it was ready to execute the instance's
	// code.  This field is primarily used for debugging and performance analysis
	// purposes.
	// Not used for template instances.
	VmmReadyTimeUs *uint64 `json:"vmm_ready_time_us,omitempty"`
	// The boot time of the instance in microseconds.  We take a pragmatic
	// approach is to define the boot time.  We calculate this as the difference
	// in time between the moment the virtualization toolstack is invoked to
	// respond to a VM boot request and the moment the OS starts executing user
	// code (i.e., the end of the guest OS boot process).  This is essentially the
	// time that a user would experience in a deployment, minus the application
	// initialization time, which we leave out since it is independent from the
	// OS.
	// Not used for template instances.
	BootTimeUs *uint64 `json:"boot_time_us,omitempty"`
	// This is the time it took for the user-level application to start listening
	// on a non-localhost port measured in microseconds.  This is the time from
	// when the instance started until it reasonably ready to start responding to
	// network requests.  This is useful for measuring the time it takes for the
	// instance to become operationally ready.
	// Not used for template instances.
	NetTimeUs *uint64 `json:"net_time_us,omitempty"`
	// The instance stop reason.
	//
	// Provides reason as to why an instance is stopped or in the process of
	// shutting down.  The stop reason is a bitmask that tells you the origin of
	// the shutdown:
	//
	// | Bit     | 4          | 3          | 2          | 1          | 0 (LSB)      |
	// |---------|------------|------------|------------|------------|--------------|
	// | Purpose | [F]orced   | [U]ser     | [P]latform | [A]pp      | [K]ernel     |
	//
	// - **Forced**:   This was a force stop.  A forced stop does not give the
	//                 instance a chance to perform a clean shutdown.  Bits 0
	//                 (Kernel) and 1 (App) can thus never be set for forced
	//                 shutdowns.  Consequently, there won't be an `exit_code` or
	//                 `stop_code`.
	// - **User**:     Stop initiated by user, e.g. via an API call.
	// - **Platform**: Stop initiated by platform, e.g. an autoscale policy.
	// - **App**:      The Application exited.  The `exit_code` field will be set.
	// - **Kernel**:   The kernel exited.  The `stop_code` field will be set.
	//
	// For example, the stop reason will contain the following values in the given
	// scenarios:
	//
	// | Value | Bitmask | Aliases | Scenario |
	// |-------|---------|---------|----------|
	// | 28    | `11100` | `FUP--` | Forced user-initiated shutdown. |
	// | 15    | `01111` | `-UPAK` | Regular user-initiated shutdown. The application and kernel have exited. The exit_code and stop_code indicate if the application and kernel shut down cleanly. |
	// | 13    | `01101` | `-UP-K` | The user initiated a shutdown but the application was forcefully killed by the kernel during shutdown. This can be the case if the image does not support a clean application exit or the application crashed after receiving a termination signal. The exit_code won’t be present in this scenario. |
	// | 7     | `00111` | `--PAK` | Unikraft Cloud initiated the shutdown, for example, due to scale-to-zero. The application and kernel have exited. The exit_code and stop_code indicate if the application and kernel shut down cleanly. |
	// | 3     | `00011` | `---AK` | The application exited. The exit_code and stop_code indicate if the application and kernel shut down cleanly. |
	// | 1     | `00001` | `----K` | The instance likely expierenced a fatal crash and the stop_code contains more information about the cause of the crash. |
	// | 0     | `00000` | `-----` | The stop reason is unknown. |
	// Not used for template instances.
	StopReason *uint32 `json:"stop_reason,omitempty"`
	// The application exit code.
	//
	// This is the code which the application returns upon leaving its main entry
	// point.  The encoding of the exit code is application specific.  See the
	// documentation of the application for more details.  Usually, an exit code
	// of `0` indicates success / no failure.
	// Not used for template instances.
	ExitCode *uint32 `json:"exit_code,omitempty"`
	// The kernel stop code.
	//
	// This value encodes multiple details about the stop irrespective of the
	// application.
	//
	// ```
	// MSB                                                     LSB
	// ┌──────────────┬──────────┬──────────┬───────────┬────────┐
	// │ 31 ────── 24 │ 23 ── 16 │    15    │ 14 ──── 8 │ 7 ── 0 │
	// ├──────────────┼──────────┼──────────┼───────────┼────────┤
	// │ reserved[^1] │ errno    │ shutdown │ initlevel │ reason │
	// └──────────────┴──────────┴──────────┴───────────┴────────┘
	// ```
	//
	// - **errno**:     The application errno, using Linux's errno.h values.
	//                  (Optional, can be 0.)
	// - **shutdown**:  Whether the shutdown originated from the inittable (0) or
	//                  from the termtable (1).
	// - **initlevel**: The initlevel at the time of the stop.
	// - **reason**:    The reason for the stop.  See `StopCodeReason`.
	//
	// [^1]: Reserved for future use.
	// Not used for template instances.
	StopCode *uint32 `json:"stop_code,omitempty"`
	// The restart configuration for the instance.
	//
	// When an instance stops either because the application exits or the instance
	// crashes, Unikraft Cloud can auto-restart your instance.  Auto-restarts are
	// performed according to the restart policy configured for a particular
	// instance.
	//
	// The policy can have the following values:
	//
	// | Policy       | Description |
	// |--------------|-------------|
	// | `never`      | Never restart the instance (default). |
	// | `always`     | Always restart the instance when the stop is initiated from within the instance (i.e., the application exits or the instance crashes). |
	// | `on-failure` | Only restart the instance if it crashes. |
	//
	// When an instance stops, the stop reason and the configured restart policy
	// are evaluated to decide if a restart should be performed.  Unikraft Cloud
	// uses an exponential back-off delay (immediate, 5s, 10s, 20s, 40s, ..., 5m)
	// to slow down restarts in tight crash loops.  If an instance runs without
	// problems for 10s the back-off delay is reset and the restart sequence ends.
	//
	// The `restart.attempt` attribute reported in counts the number of restarts
	// performed in the current sequence.  The `restart.next_at` field indicates
	// when the next restart will take place if a back-off delay is in effect.
	//
	// A manual start or stop of the instance aborts the restart sequence and
	// resets the back-off delay.
	RestartPolicy *InstanceRestartPolicy `json:"restart_policy,omitempty"`
	ScaleToZero   *InstanceScaleToZero   `json:"scale_to_zero,omitempty"`
	// The list of volumes attached to the instance.
	Volumes      []InstanceVolume      `json:"volumes,omitempty"`
	ServiceGroup *InstanceServiceGroup `json:"service_group,omitempty"`
	// The network interfaces of the instance.
	// Not used for template instances.
	NetworkInterfaces []InstanceNetworkInterface `json:"network_interfaces,omitempty"`
	// The tags associated with the instance.
	Tags []string `json:"tags,omitempty"`
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
	Error    *int32            `json:"error,omitempty"`
	Snapshot *InstanceSnapshot `json:"snapshot,omitempty"`
	// If set to true, the instance cannot be deleted until the lock is removed.
	DeleteLock *bool            `json:"delete_lock,omitempty"`
	Restart    *InstanceRestart `json:"restart,omitempty"`
}
