package model

import "github.com/ZacharyDuve/SwitchMachineDriverServer/app/controller/model"

type SwitchMachineMotorState string

// List of SwitchMachineMotorState
const (
	IDLE          SwitchMachineMotorState = "idle"
	TO_POSITION_0 SwitchMachineMotorState = "to position 0"
	TO_POSITION_1 SwitchMachineMotorState = "to position 1"
	BRAKE         SwitchMachineMotorState = "brake"
)

func MapModelMStateToAPIMState(mS model.SwitchMachineMotorState) SwitchMachineMotorState {
	if mS == model.MotorStateBrake {
		return BRAKE
	} else if mS == model.MotorStateToPos0 {
		return TO_POSITION_0
	} else if mS == model.MotorStateToPos1 {
		return TO_POSITION_1
	} else {
		return IDLE
	}
}
