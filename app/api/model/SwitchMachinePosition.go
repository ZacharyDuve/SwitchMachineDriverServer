package model

import "github.com/ZacharyDuve/SwitchMachineDriverServer/app/controller/switchmachine"

type SwitchMachinePosition string

// List of SwitchMachinePosition
const (
	Position0 SwitchMachinePosition = "position 0"
	Position1 SwitchMachinePosition = "position 1"
	Unknown   SwitchMachinePosition = "unknown"
)

func MapApiPosToModelPos(apiPos SwitchMachinePosition) switchmachine.Position {
	if apiPos == Position0 {
		return switchmachine.Position0
	} else if apiPos == Position1 {
		return switchmachine.Position1
	} else {
		return switchmachine.PositionUnknown
	}
}

func MapModelPosToApiPos(modelPos switchmachine.Position) SwitchMachinePosition {
	if modelPos == switchmachine.Position0 {
		return Position0
	} else if modelPos == switchmachine.Position1 {
		return Position1
	} else {
		return Unknown
	}
}
