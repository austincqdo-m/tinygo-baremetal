//go:build tinygo.aarch64 && qemu && virt

package runtime

import (
	"device/arm64"
	"runtime/volatile"
	"unsafe"
)

// PL011 UART on QEMU virt. See the PL011 technical reference manual:
// DR at offset 0x00, FR at offset 0x18, TXFF is FR bit 5.
const pl011Base = 0x09000000

var (
	uartData = (*volatile.Register32)(unsafe.Pointer(uintptr(pl011Base)))
	uartFlag = (*volatile.Register32)(unsafe.Pointer(uintptr(pl011Base + 0x18)))
)

func putchar(c byte) {
	for uartFlag.Get()&(1<<5) != 0 {
	}
	uartData.Set(uint32(c))
}

func getchar() byte {
	return 0
}

func buffered() int {
	return 0
}

func abort() {
	exit(1)
}

// Semihosting call numbers. See the ARM semihosting specification.
const (
	semihostingExit       = 0x18
	semihostingAppExit    = 0x20026
	semihostingExitArgsSz = 2
)

var semihostingExitArgs [semihostingExitArgsSz]uintptr

func exit(code int) {
	// Exit through semihosting so QEMU reports the program exit code.
	// The argument block is [reason, exit code].
	semihostingExitArgs[0] = semihostingAppExit
	semihostingExitArgs[1] = uintptr(code)
	semihostingCall(semihostingExit, uintptr(unsafe.Pointer(&semihostingExitArgs[0])))

	// Semihosting is not expected to return. Park as a fallback.
	for {
		arm64.Asm("wfe")
	}
}

// semihostingCall runs one semihosting call with the operation in x0 and
// the argument in x1, as the semihosting specification requires.
//
//go:export tinygo_aarch64_semihosting_call
func semihostingCall(op uintptr, arg uintptr)
