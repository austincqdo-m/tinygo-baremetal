//go:build tinygo.aarch64 && baremetal

package interrupt

import (
	"device/arm64"
)

// Enable is a no-op until an interrupt controller is configured.
func (irq Interrupt) Enable() {
}

// Disable is a no-op until an interrupt controller is configured.
func (irq Interrupt) Disable() {
}

// State represents the previous global interrupt state. It holds the DAIF
// register, so nesting works by saving and restoring the mask bits.
type State uintptr

// Disable disables all interrupts and returns the previous interrupt state.
func Disable() (state State) {
	state = State(arm64.AsmFull("mrs {}, daif", nil))
	arm64.Asm("msr daifset, #0xf")
	return state
}

// Restore restores interrupts to what they were before Disable.
func Restore(state State) {
	arm64.AsmFull("msr daif, {state}", map[string]interface{}{
		"state": uintptr(state),
	})
}

// In returns whether the system is currently in an interrupt. It remains false
// until a vector table and interrupt controller are configured.
func In() bool {
	return false
}
