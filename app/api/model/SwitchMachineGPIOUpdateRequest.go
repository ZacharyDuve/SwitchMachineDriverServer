package model

import "github.com/ZacharyDuve/SwitchMachineDriverServer/app/controller/model"

type SwitchMachineGPIOUpdateRequest struct {
	Gpio0 GPIOState `json:"gpio0,omitempty"`

	Gpio1 GPIOState `json:"gpio1,omitempty"`
}

func (this *SwitchMachineGPIOUpdateRequest) GPIO0() model.GPIOState {
	if this.Gpio0 == OFF {
		return model.GPIOOFF
	} else {
		return model.GPIOOn
	}
}

func (this *SwitchMachineGPIOUpdateRequest) GPIO1() model.GPIOState {
	if this.Gpio1 == OFF {
		return model.GPIOOFF
	} else {
		return model.GPIOOn
	}
}
