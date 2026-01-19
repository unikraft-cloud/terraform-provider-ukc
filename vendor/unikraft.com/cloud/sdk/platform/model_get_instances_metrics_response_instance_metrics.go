// This file is auto-generated. DO NOT EDIT.
// SPDX-License-Identifier: BSD-3-Clause
// Copyright (c) 2025, Unikraft GmbH.
// Licensed under the BSD-3-Clause License (the "License").
// You may not use this file except in compliance with the License.

package platform

type GetInstancesMetricsResponseInstanceMetrics struct {
	// Resident set size of the VMM in bytes.
	//
	// The resident set size (RSS) specifies the amount of physical memory that
	// has been touched by the instance and is currently reserved for the
	// instance on the Unikraft Cloud server.  The RSS grows until the instance
	// has touched all memory assigned to it via the memory_mb setting and may
	// also exceed this value as supporting services running outside the
	// instance acquire memory.  The RSS is different from the current amount of
	// memory allocated by the application, which is likely to fluctuate over
	// the lifetime of the application.  The RSS is not a cumulative metric.
	// When the instance is stopped rss goes down to 0.
	RssBytes *uint64 `json:"rss_bytes,omitempty"`
	// Consumed CPU time in milliseconds.
	CpuTimeMs *uint64 `json:"cpu_time_ms,omitempty"`
	// The boot time of the instance in microseconds.  We take a pragmatic
	// approach is to define the boot time.  We calculate this as the difference
	// in time between the moment the virtualization toolstack is invoked to
	// respond to a VM boot request and the moment the OS starts executing user
	// code (i.e., the end of the guest OS boot process).  This is essentially the
	// time that a user would experience in a deployment, minus the application
	// initialization time, which we leave out since it is independent from the
	// OS.
	BootTimeUs *uint64 `json:"boot_time_us,omitempty"`
	// This is the time it took for the user-level application to start listening
	// on a non-localhost port measured in microseconds.  This is the time from
	// when the instance started until it reasonably ready to start responding to
	// network requests.  This is useful for measuring the time it takes for the
	// instance to become operationally ready.
	NetTimeUs *uint64 `json:"net_time_us,omitempty"`
	// Total amount of bytes received from network.
	RxBytes *uint64 `json:"rx_bytes,omitempty"`
	// Total count of packets received from network.
	RxPackets *uint64 `json:"rx_packets,omitempty"`
	// Total amount of bytes transmitted over network.
	TxBytes *uint64 `json:"tx_bytes,omitempty"`
	// Total count of packets transmitted over network.
	TxPackets *uint64 `json:"tx_packets,omitempty"`
	// Number of currently established inbound connections (non-HTTP).
	Nconns *uint64 `json:"nconns,omitempty"`
	// Number of in-flight HTTP requests.
	Nreqs *uint64 `json:"nreqs,omitempty"`
	// Number of queued inbound connections and HTTP requests.
	Nqueued *uint64 `json:"nqueued,omitempty"`
	// Total number of inbound connections and HTTP requests handled.
	Ntotal *uint64 `json:"ntotal,omitempty"`
	// An optional message providing additional information about the status.
	// This field is useful when the status is not `success`.
	Message *string `json:"message,omitempty"`
	// An optional error code providing additional information about the status.
	// This field is useful when the status is not `success`.
	Error *int32 `json:"error,omitempty"`
}
