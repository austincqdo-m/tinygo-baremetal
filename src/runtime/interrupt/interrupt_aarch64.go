//go:build tinygo.aarch64 && baremetal

package interrupt

import (
	"device/arm64"
)

// Controller is the board interrupt controller. The runtime registers one
// during init. The GIC on QEMU virt implements it. Task 3.1 registers a
// Broadcom controller behind the same interface.
type Controller interface {
	EnableInterrupt(id uint32)
	DisableInterrupt(id uint32)
	SetInterruptPriority(id uint32, priority uint8)
	AcknowledgeInterrupt() uint32
	EndInterrupt(id uint32)
}

var controller Controller

// SetController registers the board interrupt controller.
func SetController(c Controller) {
	controller = c
}

// Enable enables this interrupt. Right after calling this function, the
// interrupt may be invoked if it was already pending.
func (irq Interrupt) Enable() {
	if controller != nil {
		controller.EnableInterrupt(uint32(irq.num))
	}
}

// Disable disables this interrupt.
func (irq Interrupt) Disable() {
	if controller != nil {
		controller.DisableInterrupt(uint32(irq.num))
	}
}

// SetPriority sets the interrupt priority. A lower number means a higher
// priority. The GIC delivers a pending interrupt only when its priority is
// higher than the priority of the interrupt currently being handled, so
// SetPriority is what makes nested handling possible.
func (irq Interrupt) SetPriority(priority uint8) {
	if controller != nil {
		controller.SetInterruptPriority(uint32(irq.num), priority)
	}
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

var interruptDepth int

// In reports whether execution is currently in an interrupt.
func In() bool {
	return interruptDepth > 0
}

// HandleInterrupt is called from the exception entry for IRQ and FIQ. It
// acknowledges the interrupt, unmasks interrupts so a higher priority
// interrupt can nest, calls the registered handler, then ends the interrupt.
func HandleInterrupt() {
	interruptDepth++
	if controller != nil {
		id := controller.AcknowledgeInterrupt()
		if id < 1020 {
			arm64.Asm("msr daifclr, #0xf")
			callInterruptHandler(id)
			arm64.Asm("msr daifset, #0xf")
			controller.EndInterrupt(id)
		}
	}
	interruptDepth--
}

// callHandlers dispatches to registered interrupt handlers for a given
// interrupt number. The compiler replaces each call with the handlers
// registered through interrupt.New, or removes it when there are none.
//
//go:linkname callHandlers runtime/interrupt.callHandlers
func callHandlers(num int)

// callInterruptHandler dispatches one acknowledged interrupt ID. The IDs
// covered are SGIs 0-15, PPIs 16-31, and SPIs 32-47, which cover the
// generic timer, software interrupts, and the QEMU virt devices. IDs
// outside this range are acknowledged and ended with no handler called.
// Later tasks extend this range when they need more interrupt IDs.
func callInterruptHandler(id uint32) {
	switch id {
	case 0:
		callHandlers(0)
	case 1:
		callHandlers(1)
	case 2:
		callHandlers(2)
	case 3:
		callHandlers(3)
	case 4:
		callHandlers(4)
	case 5:
		callHandlers(5)
	case 6:
		callHandlers(6)
	case 7:
		callHandlers(7)
	case 8:
		callHandlers(8)
	case 9:
		callHandlers(9)
	case 10:
		callHandlers(10)
	case 11:
		callHandlers(11)
	case 12:
		callHandlers(12)
	case 13:
		callHandlers(13)
	case 14:
		callHandlers(14)
	case 15:
		callHandlers(15)
	case 16:
		callHandlers(16)
	case 17:
		callHandlers(17)
	case 18:
		callHandlers(18)
	case 19:
		callHandlers(19)
	case 20:
		callHandlers(20)
	case 21:
		callHandlers(21)
	case 22:
		callHandlers(22)
	case 23:
		callHandlers(23)
	case 24:
		callHandlers(24)
	case 25:
		callHandlers(25)
	case 26:
		callHandlers(26)
	case 27:
		callHandlers(27)
	case 28:
		callHandlers(28)
	case 29:
		callHandlers(29)
	case 30:
		callHandlers(30)
	case 31:
		callHandlers(31)
	case 32:
		callHandlers(32)
	case 33:
		callHandlers(33)
	case 34:
		callHandlers(34)
	case 35:
		callHandlers(35)
	case 36:
		callHandlers(36)
	case 37:
		callHandlers(37)
	case 38:
		callHandlers(38)
	case 39:
		callHandlers(39)
	case 40:
		callHandlers(40)
	case 41:
		callHandlers(41)
	case 42:
		callHandlers(42)
	case 43:
		callHandlers(43)
	case 44:
		callHandlers(44)
	case 45:
		callHandlers(45)
	case 46:
		callHandlers(46)
	case 47:
		callHandlers(47)
	}
}
