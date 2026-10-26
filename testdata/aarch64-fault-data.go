package main

import "unsafe"

func main() {
	println("before fault")
	*(*byte)(unsafe.Pointer(uintptr(0x80000000))) = 1
	println("after fault")
}
