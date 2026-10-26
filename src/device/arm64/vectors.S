// AArch64 exception vector table and entry code.
// VBAR_EL1 points at _vectors, which must be 2048 byte aligned.
// Each table entry is 128 bytes apart. Entries branch to a stub that
// pushes x0/x1, sets x0 to the exception kind, and jumps to the common
// entry. The common entry saves every register a Go handler can clobber,
// calls handleAArch64Exception, restores the registers, and returns with eret.
// Frame layout, offsets in bytes:
//   0..248: x0..x30, 248: saved SP, 256: ELR, 264: SPSR, 272: ESR, 280: FAR,
//   288: FPCR, 296: FPSR, 304..816: q0..q31.

.equ FRAME_SIZE, 816
.equ SYNC_KIND, 0
.equ IRQ_KIND, 1
.equ FIQ_KIND, 2
.equ SERROR_KIND, 3

.section .vectors, "ax"
.balign 2048
.global _vectors
.type _vectors, @function
_vectors:
    .balign 0x80
    stp x0, x1, [sp, #-16]!
    mov x0, #SYNC_KIND
    b exceptionEntry
    .balign 0x80
    stp x0, x1, [sp, #-16]!
    mov x0, #IRQ_KIND
    b exceptionEntry
    .balign 0x80
    stp x0, x1, [sp, #-16]!
    mov x0, #FIQ_KIND
    b exceptionEntry
    .balign 0x80
    stp x0, x1, [sp, #-16]!
    mov x0, #SERROR_KIND
    b exceptionEntry
    .balign 0x80
    stp x0, x1, [sp, #-16]!
    mov x0, #SYNC_KIND
    b exceptionEntry
    .balign 0x80
    stp x0, x1, [sp, #-16]!
    mov x0, #IRQ_KIND
    b exceptionEntry
    .balign 0x80
    stp x0, x1, [sp, #-16]!
    mov x0, #FIQ_KIND
    b exceptionEntry
    .balign 0x80
    stp x0, x1, [sp, #-16]!
    mov x0, #SERROR_KIND
    b exceptionEntry
    .balign 0x80
    stp x0, x1, [sp, #-16]!
    mov x0, #SYNC_KIND
    b exceptionEntry
    .balign 0x80
    stp x0, x1, [sp, #-16]!
    mov x0, #IRQ_KIND
    b exceptionEntry
    .balign 0x80
    stp x0, x1, [sp, #-16]!
    mov x0, #FIQ_KIND
    b exceptionEntry
    .balign 0x80
    stp x0, x1, [sp, #-16]!
    mov x0, #SERROR_KIND
    b exceptionEntry
    .balign 0x80
    stp x0, x1, [sp, #-16]!
    mov x0, #SYNC_KIND
    b exceptionEntry
    .balign 0x80
    stp x0, x1, [sp, #-16]!
    mov x0, #IRQ_KIND
    b exceptionEntry
    .balign 0x80
    stp x0, x1, [sp, #-16]!
    mov x0, #FIQ_KIND
    b exceptionEntry
    .balign 0x80
    stp x0, x1, [sp, #-16]!
    mov x0, #SERROR_KIND
    b exceptionEntry

.section .text
.global exceptionEntry
.type exceptionEntry, @function
exceptionEntry:
    // x0 holds the kind, [sp] holds original x0, [sp, #8] holds original x1.
    // Save the kind where Go can read it. x1 is scratch here, its original
    // value is on the stack.
    adr x1, tinygo_exception_kind
    str x0, [x1]

    sub sp, sp, #FRAME_SIZE
    // Save original x2/x3 first, they become scratch below.
    str x2, [sp, #16]
    str x3, [sp, #24]
    // Load original x0/x1 from the entry push at sp + FRAME_SIZE.
    // ldp cannot encode an 816 byte offset, so form the address first.
    add x2, sp, #FRAME_SIZE
    ldp x2, x3, [x2]
    str x2, [sp, #0]
    str x3, [sp, #8]
    // Save the remaining general registers.
    stp x4, x5, [sp, #32]
    stp x6, x7, [sp, #48]
    stp x8, x9, [sp, #64]
    stp x10, x11, [sp, #80]
    stp x12, x13, [sp, #96]
    stp x14, x15, [sp, #112]
    stp x16, x17, [sp, #128]
    stp x18, x19, [sp, #144]
    stp x20, x21, [sp, #160]
    stp x22, x23, [sp, #176]
    stp x24, x25, [sp, #192]
    stp x26, x27, [sp, #208]
    stp x28, x29, [sp, #224]
    str x30, [sp, #240]
    // Saved SP is the frame base plus FRAME_SIZE plus the 16 byte entry push.
    add x2, sp, #FRAME_SIZE
    add x2, x2, #16
    str x2, [sp, #248]
    mrs x2, elr_el1
    str x2, [sp, #256]
    mrs x2, spsr_el1
    str x2, [sp, #264]
    mrs x2, esr_el1
    str x2, [sp, #272]
    mrs x2, far_el1
    str x2, [sp, #280]
    mrs x2, fpcr
    str x2, [sp, #288]
    mrs x2, fpsr
    str x2, [sp, #296]
    stp q0, q1, [sp, #304]
    stp q2, q3, [sp, #336]
    stp q4, q5, [sp, #368]
    stp q6, q7, [sp, #400]
    stp q8, q9, [sp, #432]
    stp q10, q11, [sp, #464]
    stp q12, q13, [sp, #496]
    stp q14, q15, [sp, #528]
    stp q16, q17, [sp, #560]
    stp q18, q19, [sp, #592]
    stp q20, q21, [sp, #624]
    stp q22, q23, [sp, #656]
    stp q24, q25, [sp, #688]
    stp q26, q27, [sp, #720]
    stp q28, q29, [sp, #752]
    stp q30, q31, [sp, #784]

    bl handleAArch64Exception

    // Restore in reverse order. x0 is scratch until the final restore.
    ldp q30, q31, [sp, #784]
    ldp q28, q29, [sp, #752]
    ldp q26, q27, [sp, #720]
    ldp q24, q25, [sp, #688]
    ldp q22, q23, [sp, #656]
    ldp q20, q21, [sp, #624]
    ldp q18, q19, [sp, #592]
    ldp q16, q17, [sp, #560]
    ldp q14, q15, [sp, #528]
    ldp q12, q13, [sp, #496]
    ldp q10, q11, [sp, #464]
    ldp q8, q9, [sp, #432]
    ldp q6, q7, [sp, #400]
    ldp q4, q5, [sp, #368]
    ldp q2, q3, [sp, #336]
    ldp q0, q1, [sp, #304]
    ldr x0, [sp, #288]
    msr fpcr, x0
    ldr x0, [sp, #296]
    msr fpsr, x0
    ldr x0, [sp, #256]
    msr elr_el1, x0
    ldr x0, [sp, #264]
    msr spsr_el1, x0
    ldr x0, [sp, #272]
    msr esr_el1, x0
    ldr x0, [sp, #280]
    msr far_el1, x0
    ldr x30, [sp, #240]
    ldp x28, x29, [sp, #224]
    ldp x26, x27, [sp, #208]
    ldp x24, x25, [sp, #192]
    ldp x22, x23, [sp, #176]
    ldp x20, x21, [sp, #160]
    ldp x18, x19, [sp, #144]
    ldp x16, x17, [sp, #128]
    ldp x14, x15, [sp, #112]
    ldp x12, x13, [sp, #96]
    ldp x10, x11, [sp, #80]
    ldp x8, x9, [sp, #64]
    ldp x6, x7, [sp, #48]
    ldp x4, x5, [sp, #32]
    ldp x2, x3, [sp, #16]
    ldp x0, x1, [sp, #0]
    add sp, sp, #FRAME_SIZE
    add sp, sp, #16
    eret

.section .bss
.global tinygo_exception_kind
.type tinygo_exception_kind, @object
tinygo_exception_kind:
    .space 8
