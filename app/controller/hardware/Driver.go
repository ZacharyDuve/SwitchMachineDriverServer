package hardware

import (
	"io"

	"github.com/ZacharyDuve/SwitchMachineDriverServer/app/controller/switchmachine"
)

type Driver interface {
	//Start checking for updates
	Start(DriverEventListener)
	UpdateSwitchMachine(switchmachine.State)
	io.Closer
}
