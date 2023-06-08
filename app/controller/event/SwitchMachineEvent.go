package event

import (
	"time"

	"github.com/ZacharyDuve/SwitchMachineDriverServer/app/controller/switchmachine"
)

const (
	SwitchMachineAdded   EventType = "Switch-Machine-Added"
	SwitchMachineRemoved EventType = "Switch-Machine-Removed"
	//Update of Position, Motor, or GPIO
	SwitchMachineUpdated EventType = "Switch-Machine-Updated"
)

type SwitchMachineEvent interface {
	Event
	State() switchmachine.State
}

type smEvent struct {
	eventType  EventType
	state      switchmachine.State
	originTime time.Time
}

func (this *smEvent) Type() EventType {
	return this.eventType
}

func (this *smEvent) OriginTime() time.Time {
	return this.originTime
}

func (this *smEvent) State() switchmachine.State {
	return this.state
}

func NewSwitchMachineAddedEvent(newState switchmachine.State) SwitchMachineEvent {
	return &smEvent{eventType: SwitchMachineAdded, state: newState, originTime: time.Now()}
}

func NewSwitchMachineRemovedEvent(lastState switchmachine.State) SwitchMachineEvent {
	return &smEvent{eventType: SwitchMachineRemoved, state: lastState, originTime: time.Now()}
}

func NewSwitchMachineUpdatedEvent(newState switchmachine.State) SwitchMachineEvent {
	return &smEvent{eventType: SwitchMachineUpdated, state: newState, originTime: time.Now()}
}
