package main

import "device/arm64"

func main() {
	sctlr := arm64.AsmFull("mrs {}, sctlr_el1", nil)
	println("sctlr M:", int(sctlr&1), "C:", int((sctlr>>2)&1), "I:", int((sctlr>>12)&1))
	mair := arm64.AsmFull("mrs {}, mair_el1", nil)
	println("mair device:", int(mair&0xff), "normal:", int((mair>>8)&0xff))
	tcr := arm64.AsmFull("mrs {}, tcr_el1", nil)
	println("tcr t0sz:", int(tcr&0x3f), "tg0:", int((tcr>>14)&0x3), "epd1:", int((tcr>>23)&1))
	ttbr := arm64.AsmFull("mrs {}, ttbr0_el1", nil)
	println("ttbr0 aligned:", ttbr&0xfff == 0)
}
