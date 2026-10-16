.section .init
.global _start
.type _start,@function

_start:
    // QEMU and firmware may enter at EL1, EL2, or EL3. Drop to EL1, the
    // level a kernel normally runs at. SPSR value 0x3c5 selects EL1h with
    // all interrupts masked. See ARM ARM, CurrentEL and SPSR_ELx.
    mrs x0, CurrentEL
    cmp x0, #(3 << 2)
    b.eq .Lfrom_el3
    cmp x0, #(2 << 2)
    b.eq .Lfrom_el2
    b .Lat_el1

.Lfrom_el3:
    mrs x0, scr_el3
    orr x0, x0, #(1 << 10) // RW: EL2 and below are AArch64
    orr x0, x0, #1         // NS: non-secure
    msr scr_el3, x0
    mov x0, #0x3c5
    msr spsr_el3, x0
    adr x0, .Lat_el1
    msr elr_el3, x0
    eret

.Lfrom_el2:
    mrs x0, hcr_el2
    orr x0, x0, #(1 << 31) // RW: EL1 is AArch64
    msr hcr_el2, x0
    mov x0, #0x3c5
    msr spsr_el2, x0
    adr x0, .Lat_el1
    msr elr_el2, x0
    eret

.Lat_el1:
    // Mask interrupts until a vector table is installed (task 1.2).
    msr daifset, #0xf

    // Load the stack pointer.
    ldr x0, =_stack_top
    mov sp, x0

    // Enable floating point and SIMD. Without this, any FP instruction
    // traps. CPACR_EL1 FPEN bits, ARM ARM section D1.2.22.
    mrs x0, cpacr_el1
    orr x0, x0, #(3 << 20)
    msr cpacr_el1, x0
    isb

    // Allow unaligned memory accesses. The compiler emits unaligned stack
    // stores, so SCTLR_EL1 A bit must be clear. ARM ARM section D1.2.22.
    mrs x0, sctlr_el1
    bic x0, x0, #(1 << 1)
    msr sctlr_el1, x0
    isb

    // Jump to runtime.main, which zeroes BSS, copies data, and runs Go.
    bl main

    // main should never return. Park the core if it does.
.Lpark:
    wfe
    b .Lpark

// Semihosting call. The operation number is in x0 and the argument in x1,
// following the normal calling convention. See the ARM semihosting spec.
.global tinygo_aarch64_semihosting_call
.type tinygo_aarch64_semihosting_call,@function
tinygo_aarch64_semihosting_call:
    hlt #0xf000
    ret
