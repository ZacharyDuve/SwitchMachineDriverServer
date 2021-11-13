package model

import "github.com/ZacharyDuve/SwitchMachineDriverServer/app/hardware"

type SwitchMachine struct {
	SMId int `json:"id,omitempty"`

	Pos SwitchMachinePosition `json:"position,omitempty"`

	Gpio0 GpioState `json:"gpio0,omitempty"`

	Gpio1 GpioState `json:"gpio1,omitempty"`
}

func (this *SwitchMachine) Id() hardware.SwitchMachineId {
	return hardware.SwitchMachineId(this.SMId)
}

func (this *SwitchMachine) Position() hardware.SwitchMachinePosition {
	if this.Pos == Position0 {
		return hardware.Position0
	} else {
		return hardware.Position1
	}
}

func (this *SwitchMachine) MotorState() hardware.SwitchMachineMotorState {
	return hardware.MotorStateIdle
}

func (this *SwitchMachine) GPIO0State() hardware.GPIOState {
	return mapAPIGPIOStateToHardwareState(this.Gpio0)
}

func (this *SwitchMachine) GPIO1State() hardware.GPIOState {
	return mapAPIGPIOStateToHardwareState(this.Gpio1)
}

func mapAPIGPIOStateToHardwareState(apiState GpioState) hardware.GPIOState {
	if apiState == ON {
		return hardware.GPIOOn
	} else {
		return hardware.GPIOOFF
	}
}
