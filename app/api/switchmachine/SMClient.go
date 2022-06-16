package switchmachine

import (
	"net"

	"github.com/ZacharyDuve/SwitchMachineDriverServer/app/api/model"
)

//Client written in go for current SwitchMachineHandler
type SwitchMachineClient interface {
	UpdateSwitchMachines([]*model.SwitchMachine) error
	GetSwitchMachines() ([]*model.SwitchMachine, error)
	RegisterSwitchMachineEventFunc(func(SMEvent))
}

func NewClient(a net.Addr) SwitchMachineClient {
	panic("UNIMPLEMENTED")
}
