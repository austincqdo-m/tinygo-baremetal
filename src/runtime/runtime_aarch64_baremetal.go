//go:build tinygo.aarch64 && baremetal

package runtime

import (
	"device/arm64"
	"runtime/interrupt"
	"runtime/volatile"
	"unsafe"
)

// Address 0 is never mapped on these targets, so it can stand in for
// zero-sized allocations.
const zeroSizeAllocPtr uintptr = 16

//go:extern _sbss
var _sbss [0]byte

//go:extern _ebss
var _ebss [0]byte

//go:extern _sdata
var _sdata [0]byte

//go:extern _sidata
var _sidata [0]byte

//go:extern _edata
var _edata [0]byte

var procPinnedMask interrupt.State

//go:linkname procPin sync/atomic.runtime_procPin
func procPin() {
	procPinnedMask = interrupt.Disable()
}

//go:linkname procUnpin sync/atomic.runtime_procUnpin
func procUnpin() {
	interrupt.Restore(procPinnedMask)
}

func preinit() {
	// Initialize .bss: zero-initialized global variables.
	ptr := unsafe.Pointer(&_sbss)
	for ptr != unsafe.Pointer(&_ebss) {
		*(*uint8)(ptr) = 0
		ptr = unsafe.Add(ptr, 1)
	}

	// Initialize .data: global variables initialized from flash. On the
	// targets so far the image is loaded into RAM, so this copies the
	// section onto itself, which is harmless.
	src := unsafe.Pointer(&_sidata)
	dst := unsafe.Pointer(&_sdata)
	for dst != unsafe.Pointer(&_edata) {
		*(*uint8)(dst) = *(*uint8)(src)
		dst = unsafe.Add(dst, 1)
		src = unsafe.Add(src, 1)
	}
}

//export main
func main() {
	preinit()
	run()
	exit(0)
}

func ticks() timeUnit {
	return timeUnit(arm64.AsmFull("mrs {}, cntpct_el0", nil))
}

func ticksFrequency() int64 {
	return int64(arm64.AsmFull("mrs {}, cntfrq_el0", nil))
}

func ticksToNanoseconds(ticks timeUnit) int64 {
	freq := ticksFrequency()
	if freq <= 0 {
		return int64(ticks)
	}
	return int64(ticks) * (1000 * 1000 * 1000) / freq
}

func nanosecondsToTicks(ns int64) timeUnit {
	freq := ticksFrequency()
	if freq <= 0 {
		return timeUnit(ns)
	}
	return timeUnit(ns * freq / (1000 * 1000 * 1000))
}

// Generic timer interrupt ID. The EL1 physical timer (CNTP)
// signals PPI 30 on the GIC.
const genericTimerInterruptID = 30

var timerFired volatile.Register32

func timerSetCompare(value uint64) {
	arm64.AsmFull("msr cntp_cval_el0, {v}", map[string]interface{}{
		"v": uintptr(value),
	})
}

func timerEnable() {
	arm64.AsmFull("msr cntp_ctl_el0, {v}", map[string]interface{}{
		"v": uintptr(1),
	})
}

func timerDisable() {
	arm64.AsmFull("msr cntp_ctl_el0, {v}", map[string]interface{}{
		"v": uintptr(0),
	})
}

// timerInterruptHandler runs when the generic timer fires. It
// disarms the timer so the level-triggered interrupt does not
// immediately re-enter, and wakes sleepTicks.
func timerInterruptHandler(irq interrupt.Interrupt) {
	timerDisable()
	timerFired.Set(1)
}

// initTimerInterrupt registers the generic timer handler with the
// interrupt controller. Called from board init, after the GIC is
// initialized and before interrupts are unmasked.
func initTimerInterrupt() {
	irq := interrupt.New(genericTimerInterruptID, timerInterruptHandler)
	irq.Enable()
}

func sleepTicks(d timeUnit) {
	if d <= 0 {
		return
	}
	// Clear the wake flag and program the timer with interrupts
	// disabled, so a wakeup cannot be lost between the two steps.
	// An interrupt that becomes pending before WFI still wakes
	// it, so the loop below cannot miss the wakeup.
	state := interrupt.Disable()
	timerFired.Set(0)
	timerSetCompare(uint64(ticks() + d))
	timerEnable()
	interrupt.Restore(state)
	for timerFired.Get() == 0 {
		arm64.Asm("wfi")
	}
}
