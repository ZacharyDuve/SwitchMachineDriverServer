package model

import "github.com/ZacharyDuve/SwitchMachineDriverServer/app/controller/event"

type SwitchMachineEventType string

const (
	SMAdded   SwitchMachineEventType = "SwitchMachineAdded"
	SMRemoved SwitchMachineEventType = "SwitchMachineRemoved"
	SMUpdated SwitchMachineEventType = "SwitchMachineUpdated"
)

type SwitchMachineEvent struct {
	EventType          SwitchMachineEventType `json:"eventType"`
	SwitchMachineState *SwitchMachine         `json:"switchMachineState"`
}

func MapSMEventToAPISMEventType(e event.SwitchMachineEvent) SwitchMachineEventType {
	if e.Type() == event.SwitchMachineAdded {
		return SMAdded
	} else if e.Type() == event.SwitchMachineRemoved {
		return SMRemoved
	} else if e.Type() == event.SwitchMachineUpdated {
		return SMUpdated
	} else {
		panic("Invalid event.Type unable to map")
	}
}
