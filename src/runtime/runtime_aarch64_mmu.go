//go:build tinygo.aarch64 && baremetal

package runtime

import (
	"device/arm64"
	"unsafe"
)

// Identity-mapped page tables with a 4K granule. RAM is normal
// write-back inner-shareable memory, everything below RAM is device
// memory (Device-nGnRE), and everything else is unmapped. Code is
// read-only and executable, data and the heap are non-executable.
// See ARM ARM, table D8-1 for the descriptor format.

//go:extern _pagetables
var pagetablesBase [0]byte

//go:extern _etext
var _etext [0]byte

//go:extern _ram_start
var _ram_start [0]byte

//go:extern _heap_end
var _heap_end [0]byte

const (
	pageTableEntries = 512
	pageSize         = 4096
	blockSize        = 2 * 1024 * 1024
	maxL3Tables      = 32
)

const (
	descTable = 0x3
	descBlock = 0x1
	descPage  = 0x3

	descAF         = 1 << 10
	descSHInner    = 0x3 << 8
	descAPRO       = 0x2 << 6
	descAttrDevice = 0x0 << 2
	descAttrNormal = 0x1 << 2

	descXN  = 1 << 54
	descPXN = 1 << 53
)

func pageTableAt(index int) *[pageTableEntries]uint64 {
	base := uintptr(unsafe.Pointer(&pagetablesBase))
	return (*[pageTableEntries]uint64)(unsafe.Pointer(base + uintptr(index)*pageSize))
}

func tableDescriptor(addr uintptr) uint64 {
	return uint64(addr) | descTable
}

func blockDescriptor(addr uintptr) uint64 {
	return uint64(addr) | descBlock | descAF | descSHInner | descAttrDevice
}

func pageDescriptor(addr uintptr, code bool) uint64 {
	desc := uint64(addr) | descPage | descAF | descSHInner | descAttrNormal
	if code {
		desc |= descAPRO
	} else {
		desc |= descXN | descPXN
	}
	return desc
}

// buildPageTables fills the identity-mapped tables. RAM starts at
// _ram_start and ends at _heap_end.
func buildPageTables() {
	l0 := pageTableAt(0)
	l1 := pageTableAt(1)
	l2low := pageTableAt(2)
	l2ram := pageTableAt(3)

	l1Addr := uintptr(unsafe.Pointer(l1))
	l2lowAddr := uintptr(unsafe.Pointer(l2low))
	l2ramAddr := uintptr(unsafe.Pointer(l2ram))

	l0[0] = tableDescriptor(l1Addr)
	l1[0] = tableDescriptor(l2lowAddr)
	l1[1] = tableDescriptor(l2ramAddr)

	ramBase := uintptr(unsafe.Pointer(&_ram_start))
	for i := 0; i < pageTableEntries; i++ {
		l2low[i] = blockDescriptor(uintptr(i) * blockSize)
	}

	ramEnd := uintptr(unsafe.Pointer(&_heap_end))
	codeEnd := uintptr(unsafe.Pointer(&_etext))
	nL3 := int((ramEnd - ramBase + blockSize - 1) / blockSize)
	if nL3 > maxL3Tables {
		nL3 = maxL3Tables
	}
	for i := 0; i < nL3; i++ {
		l3 := pageTableAt(4 + i)
		l2ram[i] = tableDescriptor(uintptr(unsafe.Pointer(l3)))
		for j := 0; j < pageTableEntries; j++ {
			addr := ramBase + uintptr(i)*blockSize + uintptr(j)*pageSize
			l3[j] = pageDescriptor(addr, addr < codeEnd)
		}
	}
}

// enableMMU programs MAIR_EL1, TCR_EL1, and TTBR0_EL1, invalidates
// the TLB and instruction cache, then enables the MMU and caches
// through SCTLR_EL1. The caller must have built the tables and run
// with interrupts masked. See ARM ARM, section D13.2.
func enableMMU(ttbr uintptr) {
	mair := uintptr(0xff00)
	// T0SZ 16 selects a 48 bit virtual address space, so the
	// walk starts at level 0, the first table TTBR0_EL1 names.
	tcr := uintptr(16) |
		(0x3 << 12) | // SH0: inner shareable
		(0x1 << 10) | // ORGN0: write-back, write-allocate
		(0x1 << 8) | // IRGN0: write-back, write-allocate
		// EPD0 is 0: walks enabled for TTBR0.
		(0x1 << 23) | // EPD1: walks disabled for TTBR1
		(0x2 << 32) // IPS: 40 bit physical addresses
	arm64.AsmFull(`
		dsb ishst
		msr mair_el1, {mair}
		msr tcr_el1, {tcr}
		msr ttbr0_el1, {ttbr}
		isb
		tlbi vmalle1
		dsb nsh
		isb
		ic iallu
		dsb nsh
		isb
		mrs {sctlr}, sctlr_el1
		orr {sctlr}, {sctlr}, {bits}
		msr sctlr_el1, {sctlr}
		dsb sy
		isb
	`, map[string]interface{}{
		"mair":  mair,
		"tcr":   tcr,
		"ttbr":  ttbr,
		"bits":  uintptr(1 | (1 << 2) | (1 << 12)),
		"sctlr": uintptr(0),
	})
}

// mmuInit builds the page tables and enables the MMU and caches.
// Called from preinit after BSS and data are initialized and before
// any heap allocation or interrupt unmasking.
func mmuInit() {
	buildPageTables()
	enableMMU(uintptr(unsafe.Pointer(&pagetablesBase)))
}
