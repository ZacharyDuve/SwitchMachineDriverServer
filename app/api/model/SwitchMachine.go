package model

import "github.com/ZacharyDuve/SwitchMachineDriverServer/app/controller/model"

type SwitchMachine struct {
	SMId int `json:"id,omitempty"`

	Pos SwitchMachinePosition `json:"position,omitempty"`

	Gpio0 GpioState `json:"gpio0,omitempty"`

	Gpio1 GpioState `json:"gpio1,omitempty"`
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

func mapAPIGPIOStateToHardwareState(apiState GpioState) model.GPIOState {
	if apiState == ON {
		return model.GPIOOn
	} else {
		return model.GPIOOFF
	}
}
