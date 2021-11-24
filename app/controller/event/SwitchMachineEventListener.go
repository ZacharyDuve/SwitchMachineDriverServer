package event

import "github.com/ZacharyDuve/SwitchMachineDriverServer/app/controller/model"

type SwitchMachineEventListener interface {
	SwitchMachineAdded(model.SwitchMachineState)
	SwitchMachineUpdated(model.SwitchMachineState)
	SwitchMachineRemoved(model.SwitchMachineId)
}
