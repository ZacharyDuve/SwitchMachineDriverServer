package model

import "github.com/ZacharyDuve/SwitchMachineDriverServer/app/controller/model"

type SwitchMachinePosition string

// List of SwitchMachinePosition
const (
	Position0 SwitchMachinePosition = "position 0"
	Position1 SwitchMachinePosition = "position 1"
	Unknown   SwitchMachinePosition = "unknown"
)

func mapApiPosToModelPos(apiPos SwitchMachinePosition) model.SwitchMachinePosition {
	if apiPos == Position0 {
		return model.Position0
	} else if apiPos == Position1 {
		return model.Position1
	} else {
		return model.PositionUnknown
	}
}

func mapModelPosToApiPos(modelPos model.SwitchMachinePosition) SwitchMachinePosition {
	if modelPos == model.Position0 {
		return Position0
	} else if modelPos == model.Position1 {
		return Position1
	} else {
		return Unknown
	}
}
