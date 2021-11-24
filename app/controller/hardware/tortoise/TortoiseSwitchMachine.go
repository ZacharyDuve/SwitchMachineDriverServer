package tortoise

import (
	"github.com/ZacharyDuve/SwitchMachineDriverServer/app/controller/model"
)

type tortoiseSwitchMachine struct {
	id         model.SwitchMachineId
	motorState model.SwitchMachineMotorState
	position   model.SwitchMachinePosition
	gpio0      model.GPIOState
	gpio1      model.GPIOState
}

func (this *tortoiseSwitchMachine) getCurrentState() model.SwitchMachineState {
	curState := &tortoiseState{}
	curState.id = this.id
	curState.motorState = this.motorState
	curState.position = this.position
	curState.gpio0 = this.gpio0
	curState.gpio1 = this.gpio1

	return curState
}

type tortoiseState struct {
	id         model.SwitchMachineId
	motorState model.SwitchMachineMotorState
	position   model.SwitchMachinePosition
	gpio0      model.GPIOState
	gpio1      model.GPIOState
}

func (this *tortoiseState) Id() model.SwitchMachineId {
	return this.id
}

func (this *tortoiseState) Position() model.SwitchMachinePosition {
	return this.position
}

func (this *tortoiseState) MotorState() model.SwitchMachineMotorState {
	return this.motorState
}

func (this *tortoiseState) GPIO0State() model.GPIOState {
	return this.gpio0
}

func (this *tortoiseState) GPIO1State() model.GPIOState {
	return this.gpio1
}
