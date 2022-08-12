package hardware

import (
	"io"

	"git.zmanhobbies.com/software/SwitchMachineDriverServer/app/controller/switchmachine"
)

type Driver interface {
	//Start checking for updates
	Start(DriverEventListener)
	UpdateSwitchMachine(switchmachine.State)
	io.Closer
}
