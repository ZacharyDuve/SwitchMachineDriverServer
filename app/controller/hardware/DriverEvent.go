package hardware

import "github.com/ZacharyDuve/SwitchMachineDriverServer/app/controller/switchmachine"

type DriverEventType uint8

const (
	SwitchMachineAdded           = 0
	SwitchMachineRemoved         = 1
	SwitchMachinePositionChanged = 2
)

type DriverEvent interface {
	Type() DriverEventType
	Id() switchmachine.Id
	//Could be nil
	State() switchmachine.State
}

type driverEvent struct {
	eventType DriverEventType
	id        switchmachine.Id
	state     switchmachine.State
}

func (this *driverEvent) Type() DriverEventType {
	return this.eventType
}

func (this *driverEvent) Id() switchmachine.Id {
	return this.id
}

func (this *driverEvent) State() switchmachine.State {
	return this.state
}

func NewSwitchMachineAddedEvent(id switchmachine.Id, state switchmachine.State) DriverEvent {
	return &driverEvent{eventType: SwitchMachineAdded, id: id, state: state}
}

func NewSwitchMachineRemovedEvent(id switchmachine.Id) DriverEvent {
	return &driverEvent{eventType: SwitchMachineRemoved, id: id}
}

func NewSwitchMachinePositionChangedEvent(id switchmachine.Id, state switchmachine.State) DriverEvent {
	return &driverEvent{eventType: SwitchMachinePositionChanged, id: id, state: state}
}

type DriverEventListener interface {
	HandleDriverEvent(DriverEvent)
}
