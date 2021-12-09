package hardware

import (
	"io"

	"github.com/ZacharyDuve/SwitchMachineDriverServer/app/controller/model"
)

type SwitchMachineDriver interface {
	UpdateSwitchMachine(model.SwitchMachineState)
	io.Closer
}
