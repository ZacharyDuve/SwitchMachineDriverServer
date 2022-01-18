package hardware

import (
	"errors"

	"github.com/ZacharyDuve/SwitchMachineDriverServer/app/controller/model"
)

//SwitchMachineStateBroker manages reporting states of switch machines.
type SwitchMachineStateBroker struct {
}

func NewSwitchMachineStateBroker() *SwitchMachineStateBroker {
	broker := &SwitchMachineStateBroker{}

	return broker
}

//UpdateStates takes in states and updates internally. If are new then they are added. If existing then checks to see if new state is different from old
func (this *SwitchMachineStateBroker) UpdateStates([]model.SwitchMachineState) {

}

func (this *SwitchMachineStateBroker) GetStateBySwitchMachineIds([]model.SwitchMachineId) model.SwitchMachineState {
	panic(errors.New("Unimplemented"))
}

func (this *SwitchMachineStateBroker) GetAllStates() []model.SwitchMachineState {
	resultsChan := make(chan []model.SwitchMachineState)

	return <-resultsChan
}

func (this *SwitchMachineStateBroker) RemoveStateForSwitchMachineIds([]model.SwitchMachineId) {

}

type smStateBrokerRequestType int

type smStateBrokerRequest struct {
	rType     smStateBrokerRequestType
	rChan     chan []model.SwitchMachineState
	ids       []model.SwitchMachineId
	newStates []model.SwitchMachineState
}
