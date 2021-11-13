package tortoise

import "github.com/ZacharyDuve/SwitchMachineDriverServer/app/hardware"

type tortoiseSwitchMachine struct {
	id         hardware.SwitchMachineId
	motorState hardware.SwitchMachineMotorState
	position   hardware.SwitchMachinePosition
	gpio0      hardware.GPIOState
	gpio1      hardware.GPIOState
}

func (this *tortoiseSwitchMachine) getCurrentState() hardware.SwitchMachineState {
	curState := &tortoiseState{}
	curState.id = this.id
	curState.motorState = this.motorState
	curState.position = this.position
	curState.gpio0 = this.gpio0
	curState.gpio1 = this.gpio1

	return curState
}

type tortoiseState struct {
	id         hardware.SwitchMachineId
	motorState hardware.SwitchMachineMotorState
	position   hardware.SwitchMachinePosition
	gpio0      hardware.GPIOState
	gpio1      hardware.GPIOState
}

func (this *tortoiseState) Id() hardware.SwitchMachineId {
	return this.id
}

func (this *tortoiseState) Position() hardware.SwitchMachinePosition {
	return this.position
}

func (this *tortoiseState) MotorState() hardware.SwitchMachineMotorState {
	return this.motorState
}

func (this *tortoiseState) GPIO0State() hardware.GPIOState {
	return this.gpio0
}

func (this *tortoiseState) GPIO1State() hardware.GPIOState {
	return this.gpio1
}
