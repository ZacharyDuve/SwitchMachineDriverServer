package tortoise

import (
	"encoding/hex"
	"fmt"
	"time"

	"github.com/ZacharyDuve/SwitchMachineDriverServer/app/controller/event"
	"github.com/ZacharyDuve/SwitchMachineDriverServer/app/controller/hardware"
)

func NewMockTortoiseControllerDriver(smEventListener event.SwitchMachineEventListener) (hardware.SwitchMachineDriver, error) {
	ticker := time.NewTicker(time.Second * 2)
	clsFunc := func() error {
		ticker.Stop()
		return nil
	}
	return newBaseTortiseControllerDriver(printTRXFunc, clsFunc, ticker.C, smEventListener)
}

func printTRXFunc(w, r []byte) error {
	fmt.Println(hex.EncodeToString(w))

	return nil
}
