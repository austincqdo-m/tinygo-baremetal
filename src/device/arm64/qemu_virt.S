// Semihosting call. The operation number is in x0 and the argument in x1,
// following the normal calling convention. See the ARM semihosting spec.
.section .text
.global tinygo_aarch64_semihosting_call
.type tinygo_aarch64_semihosting_call,@function
tinygo_aarch64_semihosting_call:
    hlt #0xf000
    ret
