package main

import (
	"runtime/interrupt"
	"runtime/volatile"
	"unsafe"
)

var outerCount volatile.Register32
var innerCount volatile.Register32
var sgir = (*volatile.Register32)(unsafe.Pointer(uintptr(0x08000f00)))

func outerHandler(irq interrupt.Interrupt) {
	outerCount.Set(1)
	sgir.Set(0x00010001)
}

func innerHandler(irq interrupt.Interrupt) {
	innerCount.Set(1)
}

func main() {
	outer := interrupt.New(0, outerHandler)
	outer.SetPriority(0xa0)
	outer.Enable()
	inner := interrupt.New(1, innerHandler)
	inner.SetPriority(0x80)
	inner.Enable()
	sgir.Set(0x00010000)
	for i := 0; i < 1000000; i++ {
		if outerCount.Get() == 1 && innerCount.Get() == 1 {
			break
		}
	}
	println("outer:", outerCount.Get())
	println("inner:", innerCount.Get())
}
