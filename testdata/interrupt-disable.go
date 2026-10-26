package main

import (
	"runtime/interrupt"
	"runtime/volatile"
	"unsafe"
)

var counter volatile.Register32
var sgir = (*volatile.Register32)(unsafe.Pointer(uintptr(0x08000f00)))

func handler(irq interrupt.Interrupt) {
	counter.Set(counter.Get() + 1)
}

func wait() {
	for i := 0; i < 100000; i++ {
	}
}

func main() {
	interrupt.New(0, handler).Enable()
	outer := interrupt.Disable()
	sgir.Set(0x00010000)
	wait()
	println("masked count:", counter.Get())
	inner := interrupt.Disable()
	interrupt.Restore(inner)
	sgir.Set(0x00010000)
	wait()
	println("still masked count:", counter.Get())
	interrupt.Restore(outer)
	for i := 0; i < 1000000; i++ {
		if counter.Get() == 1 {
			break
		}
	}
	println("delivered count:", counter.Get())
}
