package main

import (
	"runtime/interrupt"
	"runtime/volatile"
	"unsafe"
)

var counter volatile.Register32
var uartIMSC = (*volatile.Register32)(unsafe.Pointer(uintptr(0x09000038)))
var uartICR = (*volatile.Register32)(unsafe.Pointer(uintptr(0x09000044)))

func handler(irq interrupt.Interrupt) {
	counter.Set(counter.Get() + 1)
	uartICR.Set(1 << 5)
	uartIMSC.Set(0)
}

func main() {
	interrupt.New(33, handler).Enable()
	uartIMSC.Set(1 << 5)
	for i := 0; i < 1000000; i++ {
		if counter.Get() == 1 {
			break
		}
	}
	println("uart interrupt received")
	println("count:", counter.Get())
}
