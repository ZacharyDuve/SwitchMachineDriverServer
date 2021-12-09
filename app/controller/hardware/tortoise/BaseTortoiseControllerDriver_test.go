package tortoise

import (
	"testing"
	"time"

	"github.com/ZacharyDuve/SwitchMachineDriverServer/app/controller/hardware"
	"github.com/ZacharyDuve/SwitchMachineDriverServer/app/controller/model"
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
	_, err := newBaseTortiseControllerDriver(nil, noopCloseFunc, make(<-chan time.Time), &mockSMEventListener{})

	if err == nil {
		t.Fail()
	}
}

func TestNewBaseControllerDriverReturnsErrorIfNoCloseFuncProvided(t *testing.T) {
	_, err := newBaseTortiseControllerDriver(noopTRXFunc, nil, make(<-chan time.Time), &mockSMEventListener{})

	if err == nil {
		t.Fail()
	}
}

func TestNewBaseControllerDriverReturnsErrorIfNoEventTriggerProvided(t *testing.T) {
	_, err := newBaseTortiseControllerDriver(noopTRXFunc, noopCloseFunc, nil, &mockSMEventListener{})

	if err == nil {
		t.Fail()
	}
}

func TestNewBaseControllerDriverReturnsNoErrorIfTRXFuncAndCloseFuncProvided(t *testing.T) {
	_, err := newBaseTortiseControllerDriver(noopTRXFunc, noopCloseFunc, make(<-chan time.Time), &mockSMEventListener{})

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

	eventTrigger := make(chan time.Time)

	newBaseTortiseControllerDriver(trxFunc, noopCloseFunc, eventTrigger, &mockSMEventListener{})

	eventTrigger <- time.Now()

	<-calledTRXFuncChan
}

type mockSMEventListener struct {
}

func (this *mockSMEventListener) SwitchMachineAdded(model.SwitchMachineState) {

}

func (this *mockSMEventListener) SwitchMachineUpdated(model.SwitchMachineState) {

}
func (this *mockSMEventListener) SwitchMachineRemoved(model.SwitchMachineId) {

}

//---------------------------- isConnectedFromPositionBits ----------------------------
func TestIsConnectedFromPositionBitsReturnsFalseIfPositionIsDisconnected(t *testing.T) {
	if isConnectedFromPositionBits(positionDisconnected) {
		t.Fail()
	}
}

func TestIsConnectedFromPositionBitsReturnsTrueIfPositionIsPos0(t *testing.T) {
	if !isConnectedFromPositionBits(position0) {
		t.Fail()
	}
}

func TestIsConnectedFromPositionBitsReturnsTrueIfPositionIsPos1(t *testing.T) {
	if !isConnectedFromPositionBits(position1) {
		t.Fail()
	}
}

func TestIsConnectedFromPositionBitsReturnsTrueIfPositionIsPositionUnknown(t *testing.T) {
	if !isConnectedFromPositionBits(positionUnknown) {
		t.Fail()
	}
}

//----------------------------------------- getSMPositionFromRxBits -----------------------

func TestGetSMPositionFromRxBitsReturnsReturnsPosition0WhenBitsMapToPosition0(t *testing.T) {
	if getSMPositionFromRxBits(position0) != model.Position0 {
		t.Fail()
	}
}

func TestGetSMPositionFromRxBitsReturnsReturnsPosition1WhenBitsMapToPosition1(t *testing.T) {
	if getSMPositionFromRxBits(position1) != model.Position1 {
		t.Fail()
	}
}

func TestGetSMPositionFromRxBitsReturnsReturnsPositionUnknownWhenBitsMapToPositionUnknown(t *testing.T) {
	if getSMPositionFromRxBits(positionUnknown) != model.PositionUnknown {
		t.Fail()
	}
}

func TestGetSMPositionFromRxBitsPanicsWhenPositionBitsMapToDisconnected(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fail()
		}
	}()

	getSMPositionFromRxBits(positionDisconnected)
}

//------------------------------------ getRxBitsForPortNumber ------------------------------------------
func TestGetRxBitsForPortNumberReturnsByteWithValueInBits0And1ForPort0(t *testing.T) {
	if hasDataInBits2Through7(getRxBitsForPortNumber(0xFF, 0)) {
		t.Fail()
	}
}

func TestGetRxBitsForPortNumberReturnsByteWithValueInBits0And1ForPort1(t *testing.T) {
	if hasDataInBits2Through7(getRxBitsForPortNumber(0xFF, 1)) {
		t.Fail()
	}
}

func TestGetRxBitsForPortNumberReturnsByteWithValueInBits0And1ForPort2(t *testing.T) {
	if hasDataInBits2Through7(getRxBitsForPortNumber(0xFF, 2)) {
		t.Fail()
	}
}

func TestGetRxBitsForPortNumberReturnsByteWithValueInBits0And1ForPort3(t *testing.T) {
	if hasDataInBits2Through7(getRxBitsForPortNumber(0xFF, 3)) {
		t.Fail()
	}
}

func hasDataInBits2Through7(v byte) bool {
	return (v & 0b11111100) != 0
}

func TestGetRxBitsForPortNumberPanicsForPortNumbersGreaterThan3(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fail()
		}
	}()

	getRxBitsForPortNumber(0xFF, 4)
}

func TestGetRxBitsForPortNumberPanicsForPortNumbersLessThan0(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fail()
		}
	}()

	getRxBitsForPortNumber(0xFF, -1)
}
