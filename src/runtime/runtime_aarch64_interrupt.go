//go:build tinygo.aarch64 && baremetal

package runtime

import (
	"device/arm64"
	"runtime/interrupt"
)

//go:extern tinygo_exception_kind
var exceptionKind uintptr

const (
	exceptionKindSync   = 0
	exceptionKindIRQ    = 1
	exceptionKindFIQ    = 2
	exceptionKindSError = 3
)

//export handleAArch64Exception
func handleAArch64Exception() {
	switch exceptionKind {
	case exceptionKindIRQ, exceptionKindFIQ:
		interrupt.HandleInterrupt()
	default:
		handleSyncException()
	}
}

func handleSyncException() {
	esr := uintptr(arm64.AsmFull("mrs {}, esr_el1", nil))
	far := uintptr(arm64.AsmFull("mrs {}, far_el1", nil))
	elr := uintptr(arm64.AsmFull("mrs {}, elr_el1", nil))
	ec := (esr >> 26) & 0x3f
	printstring("panic: exception: EC=")
	printuint64(uint64(ec))
	printstring(" (")
	printstring(exceptionClassName(ec))
	printstring(") ESR=0x")
	printHex64(uint64(esr))
	printstring(" FAR=0x")
	printHex64(uint64(far))
	printstring(" ELR=0x")
	printHex64(uint64(elr))
	printstring("\r\n")
	abort()
}

func exceptionClassName(ec uintptr) string {
	switch ec {
	case 0x00:
		return "Unknown"
	case 0x01:
		return "Trapped WFI/WFE"
	case 0x02:
		return "Trapped FP/SIMD"
	case 0x15:
		return "SVC"
	case 0x20:
		return "Instruction Abort (lower)"
	case 0x21:
		return "Instruction Abort (current)"
	case 0x24:
		return "Data Abort (lower)"
	case 0x25:
		return "Data Abort (current)"
	case 0x2f:
		return "SError"
	case 0x3c:
		return "Breakpoint/BRK"
	default:
		return "Unknown"
	}
}

func printHex64(n uint64) {
	digits := "0123456789abcdef"
	started := false
	for shift := 60; shift >= 0; shift -= 4 {
		digit := byte((n >> uint(shift)) & 0xf)
		if digit != 0 || started || shift == 0 {
			started = true
			putchar(digits[digit])
		}
	}
}
