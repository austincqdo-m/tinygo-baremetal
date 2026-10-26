package main

import "device/arm64"

func main() {
	println("before fault")
	arm64.Asm("udf #0")
	println("after fault")
}
