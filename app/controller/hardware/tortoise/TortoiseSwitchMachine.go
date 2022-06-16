package tortoise

import (
	"github.com/ZacharyDuve/SwitchMachineDriverServer/app/controller/switchmachine"
)

type tortoiseSwitchMachine struct {
	id         switchmachine.Id
	motorState switchmachine.MotorState
	position   switchmachine.Position
	gpio0      switchmachine.GPIOState
	gpio1      switchmachine.GPIOState
}

func (this *tortoiseSwitchMachine) getCurrentState() switchmachine.State {
	curState := &tortoiseState{}
	curState.id = this.id
	curState.motorState = this.motorState
	curState.position = this.position
	curState.gpio0 = this.gpio0
	curState.gpio1 = this.gpio1

	return curState
}

type tortoiseState struct {
	id         switchmachine.Id
	motorState switchmachine.MotorState
	position   switchmachine.Position
	gpio0      switchmachine.GPIOState
	gpio1      switchmachine.GPIOState
}

func (this *tortoiseState) Id() switchmachine.Id {
	return this.id
}

func (this *tortoiseState) Position() switchmachine.Position {
	return this.position
}

func (this *tortoiseState) MotorState() switchmachine.MotorState {
	return this.motorState
}

func (this *tortoiseState) GPIO0State() switchmachine.GPIOState {
	return this.gpio0
}

func (this *tortoiseState) GPIO1State() switchmachine.GPIOState {
	return this.gpio1
}
