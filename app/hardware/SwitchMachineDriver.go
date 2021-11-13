package hardware

import "io"

type SwitchMachineDriver interface {
	GetNumberSwitchMachinesConnected() uint
	// GetSwitchMachines() []SwitchMachineState
	// GetStateForSwitchMachine(SwitchMachineId) SwitchMachineState
	UpdateSwitchMachine(SwitchMachineState) error
	io.Closer
}
