package model

import "github.com/ZacharyDuve/SwitchMachineDriverServer/app/controller/model"

type SwitchMachinePositionUpdateRequest struct {
	Position SwitchMachinePosition `json:"position,omitempty"`
}

func (this *SwitchMachinePositionUpdateRequest) GetPosition() model.SwitchMachinePosition {
	return MapApiPosToModelPos(this.Position)
}
