package model

import "github.com/ZacharyDuve/SwitchMachineDriverServer/app/controller/model"

type SwitchMachineId int

type SwitchMachine struct {
	SMId SwitchMachineId `json:"id"`

	Pos SwitchMachinePosition `json:"position"`

	Gpio0 GPIOState `json:"gpio0"`

	Gpio1 GPIOState `json:"gpio1"`
}

func NewAPISwitchMachineFromModel(modelSM model.SwitchMachineState) *SwitchMachine {
	apiSM := &SwitchMachine{}
	apiSM.SMId = SwitchMachineId(modelSM.Id())
	apiSM.Pos = MapModelPosToApiPos(modelSM.Position())
	apiSM.Gpio0 = MapModelGPIOToAPI(modelSM.GPIO0State())
	apiSM.Gpio1 = MapModelGPIOToAPI(modelSM.GPIO1State())
	return apiSM
}

func (this *SwitchMachine) Id() model.SwitchMachineId {
	return model.SwitchMachineId(this.SMId)
}

func (this *SwitchMachine) Position() model.SwitchMachinePosition {
	if this.Pos == Position0 {
		return model.Position0
	} else {
		return model.Position1
	}
}

func (this *SwitchMachine) GPIO0State() model.GPIOState {
	return mapAPIGPIOStateToHardwareState(this.Gpio0)
}

func (this *SwitchMachine) GPIO1State() model.GPIOState {
	return mapAPIGPIOStateToHardwareState(this.Gpio1)
}

func mapAPIGPIOStateToHardwareState(apiState GPIOState) model.GPIOState {
	if apiState == ON {
		return model.GPIOOn
	} else {
		return model.GPIOOFF
	}
}
