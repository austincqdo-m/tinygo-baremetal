//go:build tinygo.aarch64 && baremetal

package runtime

import (
	"device/arm64"
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

func sleepTicks(d timeUnit) {
	// No timer interrupt is configured, so wait by polling the counter.
	// WFI would never wake up without an interrupt source.
	deadline := ticks() + d
	for ticks() < deadline {
	}
}
