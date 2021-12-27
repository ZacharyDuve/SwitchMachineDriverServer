package tortoise

import (
	"time"

	"github.com/ZacharyDuve/SwitchMachineDriverServer/app/controller/event"
	"github.com/ZacharyDuve/SwitchMachineDriverServer/app/controller/hardware"
)

func NewMockTortoiseControllerDriver(smEventListener event.SwitchMachineEventListener, txDataOut, rxDataIn []byte) (hardware.SwitchMachineDriver, error) {
	ticker := time.NewTicker(time.Second * 2)
	clsFunc := func() error {
		ticker.Stop()
		return nil
	}

	txFunc := func(w, r []byte) error {
		//fmt.Println(hex.EncodeToString(w))
		copy(txDataOut, w)
		return nil
	}

	rxFunc := func(w, r []byte) error {
		//fmt.Println(hex.EncodeToString(w))
		copy(r, rxDataIn)
		return nil
	}
	driver, err := newBaseTortiseControllerDriver(txFunc, rxFunc, clsFunc, ticker.C, smEventListener)

	return driver, err
}

func printTRXFunc(w, r []byte) error {

	return nil
}
