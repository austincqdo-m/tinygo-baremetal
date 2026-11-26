package main

import "unsafe"

func main() {
	println("before fault")
	// The heap is non-executable. Copying a RET instruction there
	// and calling it must fault on instruction fetch.
	buf := make([]byte, 8)
	buf[0] = 0xc0
	buf[1] = 0x03
	buf[2] = 0x5f
	buf[3] = 0xd6
	// A TinyGo func value is a (context, pointer) pair.
	pair := [2]uintptr{0, uintptr(unsafe.Pointer(&buf[0]))}
	fn := *(*func())(unsafe.Pointer(&pair))
	fn()
	println("after fault")
}
