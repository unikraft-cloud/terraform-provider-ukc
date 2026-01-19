// SPDX-License-Identifier: BSD-3-Clause
// Copyright (c) 2025, Unikraft GmbH.
// Licensed under the BSD-3-Clause License (the "License").
// You may not use this file except in compliance with the License.

package platform

import (
	"fmt"
	"math"
	"strings"
	"syscall"
	"time"
)

// Stop code of the kernel.  This value encodes multiple details about the stop
// irrespective of the application.
//
// MSB                                                     LSB
// ┌──────────────┬──────────┬──────────┬───────────┬────────┐
// │ 31 ────── 24 │ 23 ── 16 │    15    │ 14 ──── 8 │ 7 ── 0 │
// ├──────────────┼──────────┼──────────┼───────────┼────────┤
// │ reserved[^1] │ errno    │ shutdown │ initlevel │ reason │
// └──────────────┴──────────┴──────────┴───────────┴────────┘
//
// [^1]:      Reserved for future use.
// errno:     The application errno, using Linux's errno.h values.  (Optional,
//            can be 0.)
// shutdown:  Whether the shutdown originated from the inittable (0) or from the
//            termtable (1).
// initlevel: The initlevel at the time of the stop.
// reason:    The reason for the stop. See `StopCodeReason`.

const (
	StopCodeMaskErrno     uint32 = 0xFF0000
	StopCodeMaskShutdown  uint32 = 0x008000
	StopCodeMaskInitLevel uint32 = 0x007F00
	StopCodeMaskReason    uint32 = 0x0000FF
)

// StopCodeErrno returns the application errno, using Linux's errno.h values.
func (instance *Instance) StopCodeErrno() uint8 {
	if instance.StopCode == nil {
		return 0
	}

	return uint8((*instance.StopCode & StopCodeMaskErrno) >> 16)
}

// StopCodeShutdownTable returns whether the stop originated from the inittable
// (0) or from the termtable (1).
func (instance *Instance) StopCodeShutdownTable() uint8 {
	if instance.StopCode == nil {
		return 0
	}

	return uint8((*instance.StopCode & StopCodeMaskShutdown) >> 15)
}

// StopCodeInitLevel returns the initlevel at the time of the stop.
func (instance *Instance) StopCodeInitLevel() uint8 {
	if instance.StopCode == nil {
		return 0
	}

	return uint8((*instance.StopCode & StopCodeMaskInitLevel) >> 8)
}

const (
	// 0 - Success
	StopCodeReasonOK uint8 = iota

	// 1 - Explicit crash (bugon/assert/crash/unhandled breakpoint)
	StopCodeReasonEXP

	// 2 - Arithmetic error
	StopCodeReasonMATH

	// 3 - Invalid CPU instruction or instruction error (e.g., operand alignment)
	StopCodeReasonINVLOP

	// 4 - Page fault - see errno (out of mem, EFAULT)
	StopCodeReasonPGFAULT

	// 5 - Segmentation fault
	StopCodeReasonSEGFAULT

	// 6 - Hardware error, NMI, alignment checks
	StopCodeReasonHWERR

	// 7 - Security violation, control protection (MTE, shadow stacks, PKU?)
	StopCodeReasonSECERR
)

// StopCodeReason returns all the reasons for the stop.
func StopCodeReasons() []string {
	return []string{
		"OK",
		"EXP",
		"MATH",
		"INVLOP",
		"PGFAULT",
		"SEGFAULT",
		"HWERR",
		"SECERR",
	}
}

// StopCodeReason provides the identity value for the reason for the stop.
func (instance *Instance) StopCodeReason() uint8 {
	if instance.StopCode == nil {
		return StopCodeReasonOK
	}

	return uint8(*instance.StopCode & StopCodeMaskReason)
}

const (
	StopReasonKernel      uint32 = 1 << iota // 0b00001
	StopReasonApplication                    // 0b00010
	StopReasonPlatform                       // 0b00100
	StopReasonUser                           // 0b01000
	StopReasonForced                         // 0b10000

	// Common stop reason scenarios as constants
	StopReasonUnknown                uint32 = 0                                                                              // 00000: -----
	StopReasonKernelCrash            uint32 = StopReasonKernel                                                               // 00001: ----K
	StopReasonAppExit                uint32 = StopReasonApplication | StopReasonKernel                                       // 00011: ---AK
	StopReasonPlatformShutdown       uint32 = StopReasonPlatform | StopReasonApplication | StopReasonKernel                  // 00111: --PAK
	StopReasonUserShutdownIncomplete uint32 = StopReasonUser | StopReasonPlatform | StopReasonKernel                         // 01101: -UP-K
	StopReasonUserShutdownComplete   uint32 = StopReasonUser | StopReasonPlatform | StopReasonApplication | StopReasonKernel // 01111: -UPAK
	StopReasonForcedUserShutdown     uint32 = StopReasonForced | StopReasonUser | StopReasonPlatform                         // 11100: FUP--
)

// DescribeStopOrigin provides a human-readable interpretation of the stop
// reason.
func (instance *Instance) DescribeStopOrigin() string {
	if instance.StopReason == nil || *instance.StopReason == 0 {
		return "unknown"
	}

	var ret strings.Builder

	if *instance.StopReason&StopReasonForced != 0 {
		ret.WriteString("force ")
	}

	ret.WriteString("initiated by ")

	switch true {
	case *instance.StopReason&StopReasonPlatform == StopReasonPlatform && *instance.StopReason&StopReasonUser != StopReasonUser:
		ret.WriteString("platform")
	case *instance.StopReason&StopReasonUser == StopReasonUser:
		ret.WriteString("user")
	case *instance.StopReason&StopReasonApplication == StopReasonApplication:
		ret.WriteString("app")
	case *instance.StopReason&StopReasonKernel == StopReasonKernel:
		ret.WriteString("kernel")
	default:
		ret.WriteString("unknown")
	}

	return ret.String()
}

// StopOriginCode provides a human-readable interpretation of the stop reason in
// the form of a short-code.
func (instance *Instance) StopOriginCode() string {
	if instance.StopCode == nil || *instance.StopCode == 0 {
		return "-----"
	}

	var ret strings.Builder

	if *instance.StopReason&StopReasonForced == StopReasonForced {
		ret.WriteString("f")
	} else {
		ret.WriteString("-")
	}
	if *instance.StopReason&StopReasonUser == StopReasonUser {
		ret.WriteString("u")
	} else {
		ret.WriteString("-")
	}
	if *instance.StopReason&StopReasonPlatform == StopReasonPlatform {
		ret.WriteString("p")
	} else {
		ret.WriteString("-")
	}
	if *instance.StopReason&StopReasonApplication == StopReasonApplication {
		ret.WriteString("a")
	} else {
		ret.WriteString("-")
	}
	if *instance.StopReason&StopReasonKernel == StopReasonKernel {
		ret.WriteString("k")
	} else {
		ret.WriteString("-")
	}

	return ret.String()
}

// DescribeStopReason provides a human-readable description of the stop reason.
func (instance *Instance) DescribeStopReason() string {
	if instance.StopCode == nil || *instance.StopCode == 0 {
		return ""
	}

	var ret strings.Builder

	switch true {
	case instance.StopCodeReason() == StopCodeReasonOK:
		ret.WriteString("shutdown")
	case instance.StopCodeReason() == StopCodeReasonEXP:
		ret.WriteString("assertion error")
	case instance.StopCodeReason() == StopCodeReasonPGFAULT && instance.StopCodeErrno() == 0xc:
		ret.WriteString("out of memory")
	case instance.StopCodeReason() == StopCodeReasonPGFAULT && (instance.StopCodeErrno() == 0xe || instance.StopCodeErrno() == 0x1):
		ret.WriteString("illegal memory access")
	case instance.StopCodeReason() == StopCodeReasonSEGFAULT:
		ret.WriteString("segmentation fault")
	case instance.StopCodeReason() == StopCodeReasonPGFAULT:
		ret.WriteString("page fault")
	case instance.StopCodeReason() == StopCodeReasonMATH:
		ret.WriteString("arithmetic error")
	case instance.StopCodeReason() == StopCodeReasonINVLOP:
		ret.WriteString("instruction error")
	case instance.StopCodeReason() == StopCodeReasonHWERR:
		ret.WriteString("hardware error")
	case instance.StopCodeReason() == StopCodeReasonSECERR:
		ret.WriteString("security violation")
	default:
		ret.WriteString("unexpected error")
	}

	return ret.String()
}

// StopReasonCode returns a human-readable short-code representation of the stop
// reason.
func (instance *Instance) StopReasonCode() string {
	if instance.StopCode == nil || *instance.StopCode == 0 {
		return ""
	}

	var ret strings.Builder

	if instance.StopCodeShutdownTable() == 0 {
		ret.WriteString("i")
	} else {
		ret.WriteString("t")
	}

	ret.WriteString(fmt.Sprintf("%d", instance.StopCodeInitLevel()))

	ret.WriteString(" ")

	ret.WriteString(StopCodeReasons()[instance.StopCodeReason()])

	if instance.StopCodeErrno() != 0 {
		ret.WriteString(" ")
		errno, ok := ErrnoNames()[syscall.Errno(instance.StopCodeErrno())]
		if ok {
			ret.WriteString(errno)
		} else {
			ret.WriteString(fmt.Sprintf("%d", instance.StopCodeErrno()))
		}
	}

	return ret.String()
}

// DescribeStatus returns a human-readable description of the instance's status.
func (instance *Instance) DescribeStatus() string {
	if instance.State == nil {
		return ""
	}

	switch *instance.State {
	case InstanceStateRunning:
		dur := time.Since(*instance.StartedAt)
		days := int64(dur.Hours() / 24)
		hours := int64(math.Mod(dur.Hours(), 24))
		minutes := int64(math.Mod(dur.Minutes(), 60))
		seconds := int64(math.Mod(dur.Seconds(), 60))

		chunks := []struct {
			singularName string
			amount       int64
		}{
			{"day", days},
			{"hr", hours},
			{"min", minutes},
			{"sec", seconds},
		}

		parts := []string{}

		for i, chunk := range chunks {
			if len(parts) > 0 && i+1 == len(chunks) { // Skip secs if greater than 1m
				continue
			}
			switch chunk.amount {
			case 0:
				continue
			case 1:
				parts = append(parts, fmt.Sprintf("%d%s", chunk.amount, chunk.singularName))
			default:
				parts = append(parts, fmt.Sprintf("%d%ss", chunk.amount, chunk.singularName))
			}
		}

		return fmt.Sprintf("since %s", strings.Join(parts, " "))
	case InstanceStateStopped:
		reason := instance.DescribeStopReason()
		if reason == "shutdown" {
			return ""
		}

		return reason
	default:
		return string(*instance.State)
	}
}
