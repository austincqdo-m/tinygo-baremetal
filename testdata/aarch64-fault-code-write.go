package main

import "unsafe"

func main() {
	println("before fault")
	// The vector table is at the start of RAM, in the read-only
	// code region. Writing to it must fault.
	*(*byte)(unsafe.Pointer(uintptr(0x40000000))) = 1
	println("after fault")
}
