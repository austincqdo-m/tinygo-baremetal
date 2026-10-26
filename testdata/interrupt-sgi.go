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

func main() {
	interrupt.New(0, handler).Enable()
	sgir.Set(0x00010000)
	for i := 0; i < 1000000; i++ {
		if counter.Get() == 1 {
			break
		}
	}
	println("sgi interrupt received")
	println("count:", counter.Get())
}
