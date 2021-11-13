package hardware

import "errors"

//SwitchMachineStateBroker manages reporting states of switch machines.
type SwitchMachineStateBroker struct {
}

func NewSwitchMachineStateBroker() *SwitchMachineStateBroker {
	broker := &SwitchMachineStateBroker{}

	return broker
}

//UpdateStates takes in states and updates internally. If are new then they are added. If existing then checks to see if new state is different from old
func (this *SwitchMachineStateBroker) UpdateStates([]SwitchMachineState) {

}

func (this *SwitchMachineStateBroker) GetStateBySwitchMachineIds([]SwitchMachineId) SwitchMachineState {
	panic(errors.New("Unimplemented"))
}

func (this *SwitchMachineStateBroker) GetAllStates() []SwitchMachineState {
	resultsChan := make(chan []SwitchMachineState)

	return <-resultsChan
}

func (this *SwitchMachineStateBroker) RemoveStateForSwitchMachineIds([]SwitchMachineId) {

}

type smStateBrokerRequestType int

type smStateBrokerRequest struct {
	rType     smStateBrokerRequestType
	rChan     chan []SwitchMachineState
	ids       []SwitchMachineId
	newStates []SwitchMachineState
}
