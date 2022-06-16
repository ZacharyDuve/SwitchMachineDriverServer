package event

type SwitchMachineEventListener interface {
	HandleSwitchMachineEvent(SwitchMachineEvent)
}
