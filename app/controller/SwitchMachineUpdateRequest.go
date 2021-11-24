package controller

import "github.com/ZacharyDuve/SwitchMachineDriverServer/app/controller/model"

type SwitchMachineUpdateRequest interface {
	Id() model.SwitchMachineId
	Position() model.SwitchMachinePosition
	GPIO0State() model.GPIOState
	GPIO1State() model.GPIOState
}

func NewSwitchMachineUpdateRequest(id model.SwitchMachineId, pos model.SwitchMachinePosition, g0, g1 model.GPIOState) SwitchMachineUpdateRequest {
	return &switchMachineRequestImpl{id: id, pos: pos, gpio0: g0, gpio1: g1}
}

type switchMachineRequestImpl struct {
	id    model.SwitchMachineId
	pos   model.SwitchMachinePosition
	gpio0 model.GPIOState
	gpio1 model.GPIOState
}

func (this *switchMachineRequestImpl) Id() model.SwitchMachineId {
	return this.id
}

func (this *switchMachineRequestImpl) Position() model.SwitchMachinePosition {
	return this.pos
}

func (this *switchMachineRequestImpl) GPIO0State() model.GPIOState {
	return this.gpio0
}

func (this *switchMachineRequestImpl) GPIO1State() model.GPIOState {
	return this.gpio1
}
