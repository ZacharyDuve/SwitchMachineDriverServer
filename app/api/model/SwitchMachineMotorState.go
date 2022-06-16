package model

import (
	"github.com/ZacharyDuve/SwitchMachineDriverServer/app/controller/switchmachine"
)

type SwitchMachineMotorState string

// List of SwitchMachineMotorState
const (
	IDLE          SwitchMachineMotorState = "idle"
	TO_POSITION_0 SwitchMachineMotorState = "to position 0"
	TO_POSITION_1 SwitchMachineMotorState = "to position 1"
	BRAKE         SwitchMachineMotorState = "brake"
)

func MapModelMStateToAPIMState(mS switchmachine.MotorState) SwitchMachineMotorState {
	if mS == switchmachine.MotorStateBrake {
		return BRAKE
	} else if mS == switchmachine.MotorStateToPos0 {
		return TO_POSITION_0
	} else if mS == switchmachine.MotorStateToPos1 {
		return TO_POSITION_1
	} else {
		return IDLE
	}
}
