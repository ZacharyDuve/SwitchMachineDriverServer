package tortoise

import (
	"testing"

	"github.com/ZacharyDuve/SwitchMachineDriverServer/app/hardware"
)

func noopTRXFunc(w, r []byte) error {
	return nil
}

func noopCloseFunc() error {
	return nil
}

//Test to see if baseTortoiseControllerDriver implements functions of TortoiseControllerDriver. Compile should fail if it doesn't
func TestBaseTortoiseControllerDriverImplementsTortoiseControllerDriver(t *testing.T) {
	var _ hardware.SwitchMachineDriver = &baseTortoiseControllerDriver{}
}

//------------------------------------- newBaseTortoiseControllerDriver tests ---------------------------------
func TestNewBaseControllerDriverReturnsErrorIfNoTRXFuncProvided(t *testing.T) {
	_, err := newBaseTortiseControllerDriver(nil, noopCloseFunc)

	if err == nil {
		t.Fail()
	}
}

func TestNewBaseControllerDriverReturnsErrorIfNoCloseFuncProvided(t *testing.T) {
	_, err := newBaseTortiseControllerDriver(noopTRXFunc, nil)

	if err == nil {
		t.Fail()
	}
}

func TestNewBaseControllerDriverReturnsNoErrorIfTRXFuncAndCloseFuncProvided(t *testing.T) {
	_, err := newBaseTortiseControllerDriver(noopTRXFunc, noopCloseFunc)

	if err != nil {
		t.Fail()
	}
}

//Test that writing to trxfunc works as intended
func TestThatTRXFuncIsCalled(t *testing.T) {
	calledTRXFuncChan := make(chan bool)
	trxFunc := func(w, r []byte) error {
		calledTRXFuncChan <- true
		return nil
	}
	newBaseTortiseControllerDriver(trxFunc, noopCloseFunc)

	<-calledTRXFuncChan
}
