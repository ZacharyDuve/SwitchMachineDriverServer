package switchmachine

import (
	"time"
)

type SwitchMachineId uint16

type SwitchPosition uint8

const (
	PositionUnknown SwitchPosition = iota
	PositionA
	PositionB
)

type SwitchMotorState uint8

const (
	Off SwitchMotorState = iota
	ToPositionA
	ToPositionB
)

type SwitchMachine struct {
	// Unique identifier for a switch machine on this server
	id SwitchMachineId
	// Is the switch machine currently connected
	connected bool
	// Current position of the switch machine
	position SwitchPosition
	// Current state of the motor of the switch machine
	motorState SwitchMotorState
	// Is the motor reversed. False No and True Yes
	motorReversed bool
	// Throw Time for motor ToPositionA
	throwTimePositionA time.Duration
	// Throw Time for motor ToPositionB
	throwTimePositionB time.Duration
}
