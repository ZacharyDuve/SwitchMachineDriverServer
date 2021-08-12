package model

type SwitchMachineMotorState string

// List of SwitchMachineMotorState
const (
	IDLE          SwitchMachineMotorState = "idle"
	TO_POSITION_0 SwitchMachineMotorState = "to position 0"
	TO_POSITION_1 SwitchMachineMotorState = "to position 1"
	BRAKE         SwitchMachineMotorState = "brake"
)
