//go:build tinygo.aarch64 && qemu && virt

package runtime

import (
	"device/arm64"
	"runtime/interrupt"
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

// GICv2 on QEMU virt. Distributor at 0x08000000, CPU interface at
// 0x08010000. See the GICv2 architecture specification.
const (
	gicDistributorBase = 0x08000000
	gicCPUBase         = 0x08010000
)

func gicdReg(offset uintptr) *volatile.Register32 {
	return (*volatile.Register32)(unsafe.Pointer(uintptr(gicDistributorBase + offset)))
}

func giccReg(offset uintptr) *volatile.Register32 {
	return (*volatile.Register32)(unsafe.Pointer(uintptr(gicCPUBase + offset)))
}

type gicController struct{}

func (gicController) EnableInterrupt(id uint32) {
	gicdReg(0x100 + uintptr(id/32)*4).SetBits(1 << (id % 32))
}

func (gicController) DisableInterrupt(id uint32) {
	gicdReg(0x180 + uintptr(id/32)*4).SetBits(1 << (id % 32))
}

func (gicController) SetInterruptPriority(id uint32, priority uint8) {
	(*volatile.Register8)(unsafe.Pointer(uintptr(gicDistributorBase + 0x400 + uintptr(id)))).Set(priority)
}

func (gicController) AcknowledgeInterrupt() uint32 {
	return giccReg(0x00c).Get() & 0x3ff
}

func (gicController) EndInterrupt(id uint32) {
	giccReg(0x010).Set(id)
}

//go:extern _vectors
var vectorsSymbol [0]byte

func gicInit() {
	gicdReg(0x000).Set(0)
	giccReg(0x000).Set(0)
	lines := (gicdReg(0x004).Get() & 0x1f) + 1
	for i := uint32(0); i < lines; i++ {
		gicdReg(0x180 + uintptr(i)*4).Set(0xffffffff)
		gicdReg(0x280 + uintptr(i)*4).Set(0xffffffff)
		gicdReg(0x380 + uintptr(i)*4).Set(0xffffffff)
	}
	for id := uint32(0); id < lines*32; id++ {
		(*volatile.Register8)(unsafe.Pointer(uintptr(gicDistributorBase + 0x400 + uintptr(id)))).Set(0xa0)
		if id >= 32 {
			(*volatile.Register8)(unsafe.Pointer(uintptr(gicDistributorBase + 0x800 + uintptr(id)))).Set(0x01)
		}
	}
	gicdReg(0x000).Set(0x01)
	giccReg(0x004).Set(0xff)
	giccReg(0x008).Set(0)
	giccReg(0x000).Set(0x01)
}

func init() {
	interrupt.SetController(gicController{})
	arm64.AsmFull("msr vbar_el1, {addr}", map[string]interface{}{
		"addr": uintptr(unsafe.Pointer(&vectorsSymbol)),
	})
	arm64.Asm("isb")
	gicInit()
	initTimerInterrupt()
	arm64.Asm("msr daifclr, #0xf")
}
