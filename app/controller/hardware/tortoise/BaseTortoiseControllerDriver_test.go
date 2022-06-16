package tortoise

import (
	"testing"
	"time"

	"github.com/ZacharyDuve/SwitchMachineDriverServer/app/controller/hardware"
	"github.com/ZacharyDuve/SwitchMachineDriverServer/app/controller/switchmachine"
)

func noopTRXFunc(w, r []byte) error {
	return nil
}

func noopCloseFunc() error {
	return nil
}

//Test to see if baseTortoiseControllerDriver implements functions of TortoiseControllerDriver. Compile should fail if it doesn't
func TestBaseTortoiseControllerDriverImplementsTortoiseControllerDriver(t *testing.T) {
	var _ hardware.Driver = &baseTortoiseControllerDriver{}
}

//Test that writing to trxfunc works as intended
func TestThatTRXFuncIsCalled(t *testing.T) {
	calledTRXFuncChan := make(chan bool)
	trxFunc := func(w, r []byte) error {
		calledTRXFuncChan <- true
		return nil
	}

	eventTrigger := make(chan time.Time)
	driver := getBaseDriverWithAllNOOP()
	driver.rxFunc = trxFunc
	driver.rxTrigger = eventTrigger

	driver.Start(&mockDriverEventListener{})

	eventTrigger <- time.Now()

	<-calledTRXFuncChan
}

//---------------------------- isConnectedFromPositionBits ----------------------------
func TestIsConnectedFromPositionBitsReturnsFalseIfPositionIsDisconnected(t *testing.T) {
	if isConnectedFromPositionBits(positionDisconnected) {
		t.Fail()
	}
}

func TestIsConnectedFromPositionBitsReturnsTrueIfPositionIsPos0(t *testing.T) {
	if !isConnectedFromPositionBits(position0Port03) {
		t.Fail()
	}
}

func TestIsConnectedFromPositionBitsReturnsTrueIfPositionIsPos1(t *testing.T) {
	if !isConnectedFromPositionBits(position1Port03) {
		t.Fail()
	}
}

func TestIsConnectedFromPositionBitsReturnsTrueIfPositionIsPositionUnknown(t *testing.T) {
	if !isConnectedFromPositionBits(positionUnknown) {
		t.Fail()
	}
}

//----------------------------------------- getSMPositionFromRxBits -----------------------

func TestGetSMPositionFromRxBitsReturnsReturnsPosition0WhenBitsMapToPosition0ForPort0(t *testing.T) {
	if getSMPositionFromRxBits(position0Port03, 0) != switchmachine.Position0 {
		t.Fail()
	}
}

func TestGetSMPositionFromRxBitsReturnsReturnsPosition1WhenBitsMapToPosition1ForPort0(t *testing.T) {
	if getSMPositionFromRxBits(position1Port03, 0) != switchmachine.Position1 {
		t.Fail()
	}
}

func TestGetSMPositionFromRxBitsReturnsReturnsPosition0WhenBitsMapToPosition0ForPort1(t *testing.T) {
	if getSMPositionFromRxBits(position0Port12, 1) != switchmachine.Position0 {
		t.Fail()
	}
}

func TestGetSMPositionFromRxBitsReturnsReturnsPosition1WhenBitsMapToPosition1ForPort1(t *testing.T) {
	if getSMPositionFromRxBits(position1Port12, 1) != switchmachine.Position1 {
		t.Fail()
	}
}

func TestGetSMPositionFromRxBitsReturnsReturnsPosition0WhenBitsMapToPosition0ForPort2(t *testing.T) {
	if getSMPositionFromRxBits(position0Port12, 2) != switchmachine.Position0 {
		t.Fail()
	}
}

func TestGetSMPositionFromRxBitsReturnsReturnsPosition1WhenBitsMapToPosition1ForPort2(t *testing.T) {
	if getSMPositionFromRxBits(position1Port12, 2) != switchmachine.Position1 {
		t.Fail()
	}
}

func TestGetSMPositionFromRxBitsReturnsReturnsPosition0WhenBitsMapToPosition0ForPort3(t *testing.T) {
	if getSMPositionFromRxBits(position0Port03, 0) != switchmachine.Position0 {
		t.Fail()
	}
}

func TestGetSMPositionFromRxBitsReturnsReturnsPosition1WhenBitsMapToPosition1ForPort3(t *testing.T) {
	if getSMPositionFromRxBits(position1Port03, 0) != switchmachine.Position1 {
		t.Fail()
	}
}

func TestGetSMPositionFromRxBitsReturnsReturnsPositionUnknownWhenBitsMapToPositionUnknown(t *testing.T) {
	if getSMPositionFromRxBits(positionUnknown, 0) != switchmachine.PositionUnknown {
		t.Fail()
	}
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
	return v&0b11111100 != 0
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

//--------------------------------------------- calcTxByteOffsetFromId -----------------------
func TestWhencalcTxByteOffsetFromIdWithId0Returns0(t *testing.T) {
	if calcTxByteOffsetFromId(switchmachine.Id(0)) != 0 {
		t.Fail()
	}
}

func TestWhencalcTxByteOffsetFromIdWithId1Returns0(t *testing.T) {
	if calcTxByteOffsetFromId(switchmachine.Id(1)) != 0 {
		t.Fail()
	}
}

func TestWhencalcTxByteOffsetFromIdWithId2Returns1(t *testing.T) {
	if calcTxByteOffsetFromId(switchmachine.Id(2)) != 1 {
		t.Fail()
	}
}

func TestWhencalcTxByteOffsetFromIdWithId3Returns1(t *testing.T) {
	if calcTxByteOffsetFromId(switchmachine.Id(3)) != 1 {
		t.Fail()
	}
}

func TestWhencalcTxByteOffsetFromIdWithId4Returns2(t *testing.T) {
	if calcTxByteOffsetFromId(switchmachine.Id(4)) != 2 {
		t.Fail()
	}
}

//-------------------------------Close--------------------------------

func TestCallingCloseCallsConfiguredCloseFunc(t *testing.T) {
	wasCloseCalled := false
	eventTrigger := make(chan time.Time)
	driver := getBaseDriverWithAllNOOP()
	driver.closeFunc = func() error {
		wasCloseCalled = true
		return nil
	}
	driver.rxTrigger = eventTrigger

	driver.Start(&mockDriverEventListener{})

	driver.Close()
	if !wasCloseCalled {
		t.Fail()
	}
}

//---------------------------------RX-----------------------------------

func TestThatHavingSwitchMachineConnectOnId0InUnknownPositionCausesSwitchMachineAddedEventToBeFired(t *testing.T) {
	wasExpectedEventFired := false
	eventTrigger := make(chan time.Time)
	waitChan := make(chan bool)
	driver := getBaseDriverWithAllNOOP()
	driver.rxFunc = func(w, r []byte) error {
		copy(r, []byte{port0RxBitMask})
		return nil
	}
	driver.rxTrigger = eventTrigger
	driver.Start(&mockDriverEventListener{eventHandlerFunc: func(de hardware.DriverEvent) {
		wasExpectedEventFired = de.Type() == hardware.SwitchMachineAdded && de.State().Position() == switchmachine.PositionUnknown && de.Id() == switchmachine.Id(0)
		waitChan <- true
	}})
	eventTrigger <- time.Now()
	<-waitChan
	if !wasExpectedEventFired {
		t.Fail()
	}
}

func TestThatHavingSwitchMachineConnectOnId0InPosition0CausesSwitchMachineAddedEventToBeFiredWithCorrectPosition(t *testing.T) {
	wasExpectedEventFired := false
	eventTrigger := make(chan time.Time)
	waitChan := make(chan bool)
	driver := getBaseDriverWithAllNOOP()
	driver.rxFunc = func(w, r []byte) error {
		copy(r, []byte{position0Port03 << (numBitsPerPort * port0RxBitIndex)})
		return nil
	}
	driver.rxTrigger = eventTrigger
	driver.Start(&mockDriverEventListener{eventHandlerFunc: func(de hardware.DriverEvent) {
		wasExpectedEventFired = de.Type() == hardware.SwitchMachineAdded && de.State().Position() == switchmachine.Position0 && de.Id() == switchmachine.Id(0)
		waitChan <- true
	}})
	eventTrigger <- time.Now()
	<-waitChan
	if !wasExpectedEventFired {
		t.Fail()
	}
}

func TestThatUpdatingSwitchMachineConnectOnId0FromPosition0To1CausesUpdateEventToBeFiredWithPosition1(t *testing.T) {
	wasExpectedEventFired := false
	eventTrigger := make(chan time.Time)
	waitChan := make(chan bool)
	driver := getBaseDriverWithAllNOOP()
	driver.rxFunc = func(w, r []byte) error {
		copy(r, []byte{position0Port03 << (numBitsPerPort * port0RxBitIndex)})
		return nil
	}
	driver.rxTrigger = eventTrigger
	driver.Start(&mockDriverEventListener{eventHandlerFunc: func(de hardware.DriverEvent) {
		wasExpectedEventFired = de.Type() == hardware.SwitchMachinePositionChanged && de.State().Position() == switchmachine.Position1 && de.Id() == switchmachine.Id(0)
		if wasExpectedEventFired {
			waitChan <- true
		}
	}})
	eventTrigger <- time.Now()

	driver.rxFunc = func(w, r []byte) error {
		copy(r, []byte{position1Port03 << (numBitsPerPort * port0RxBitIndex)})
		return nil
	}

	eventTrigger <- time.Now()

	<-waitChan
	if !wasExpectedEventFired {
		t.Fail()
	}
}

func TestThatUpdatingSwitchMachineConnectOnId0FromPosition0To1CausesRemovedEventToBeFired(t *testing.T) {
	wasExpectedEventFired := false
	eventTrigger := make(chan time.Time)
	waitChan := make(chan bool)
	driver := getBaseDriverWithAllNOOP()
	driver.rxFunc = func(w, r []byte) error {
		copy(r, []byte{position0Port03 << (numBitsPerPort * port0RxBitIndex)})
		return nil
	}
	driver.rxTrigger = eventTrigger
	driver.Start(&mockDriverEventListener{eventHandlerFunc: func(de hardware.DriverEvent) {
		wasExpectedEventFired = de.Type() == hardware.SwitchMachineRemoved && de.Id() == switchmachine.Id(0)
		if wasExpectedEventFired {
			waitChan <- true
		}
	}})
	eventTrigger <- time.Now()

	driver.rxFunc = func(w, r []byte) error {
		copy(r, []byte{positionDisconnected << (numBitsPerPort * port0RxBitIndex)})
		return nil
	}

	eventTrigger <- time.Now()

	<-waitChan
	if !wasExpectedEventFired {
		t.Fail()
	}
}

//------------------------------------UpdateSwitchMachine----------------------------------
func TestThatUpdatingGPIO0OnPort0Causes0x10ToBeWrittenForCorrectByte(t *testing.T) {
	wasTxWrittenAsExpected := false
	waitChan := make(chan bool)
	driver := getBaseDriverWithAllNOOP()
	driver.Start(&mockDriverEventListener{})
	idUnderTest := switchmachine.Id(0)
	driver.txFunc = func(w, r []byte) error {
		byteIndex := getTxIndexFromBufferLengthAndId(len(driver.txBuffer), idUnderTest)
		if w[byteIndex] == 0x10 {
			wasTxWrittenAsExpected = true
			waitChan <- true
		}
		return nil
	}
	driver.UpdateSwitchMachine(switchmachine.NewState(idUnderTest, switchmachine.PositionUnknown, switchmachine.MotorStateIdle, switchmachine.GPIOOn, switchmachine.GPIOOFF))

	<-waitChan
	if !wasTxWrittenAsExpected {
		t.Fail()
	}
}

func TestThatUpdatingGPIO1OnPort0Causes0x20ToBeWrittenForCorrectByte(t *testing.T) {
	wasTxWrittenAsExpected := false
	waitChan := make(chan bool)
	driver := getBaseDriverWithAllNOOP()
	driver.Start(&mockDriverEventListener{})
	idUnderTest := switchmachine.Id(0)
	driver.txFunc = func(w, r []byte) error {
		byteIndex := getTxIndexFromBufferLengthAndId(len(driver.txBuffer), idUnderTest)
		if w[byteIndex] == 0x20 {
			wasTxWrittenAsExpected = true
			waitChan <- true
		}
		return nil
	}
	driver.UpdateSwitchMachine(switchmachine.NewState(idUnderTest, switchmachine.PositionUnknown, switchmachine.MotorStateIdle, switchmachine.GPIOOFF, switchmachine.GPIOOn))

	<-waitChan
	if !wasTxWrittenAsExpected {
		t.Fail()
	}
}

func TestThatUpdatingBothGPIOOnPort0Causes0x20ToBeWrittenForCorrectByte(t *testing.T) {
	wasTxWrittenAsExpected := false
	waitChan := make(chan bool)
	driver := getBaseDriverWithAllNOOP()
	driver.Start(&mockDriverEventListener{})
	idUnderTest := switchmachine.Id(0)
	driver.txFunc = func(w, r []byte) error {
		byteIndex := getTxIndexFromBufferLengthAndId(len(driver.txBuffer), idUnderTest)
		if w[byteIndex] == 0x30 {
			wasTxWrittenAsExpected = true
			waitChan <- true
		}
		return nil
	}
	driver.UpdateSwitchMachine(switchmachine.NewState(idUnderTest, switchmachine.PositionUnknown, switchmachine.MotorStateIdle, switchmachine.GPIOOn, switchmachine.GPIOOn))

	<-waitChan
	if !wasTxWrittenAsExpected {
		t.Fail()
	}
}

func getBaseDriverWithAllNOOP() *baseTortoiseControllerDriver {
	driver := &baseTortoiseControllerDriver{}
	driver.closeFunc = noopCloseFunc
	driver.rxFunc = noopTRXFunc
	driver.txFunc = noopTRXFunc

	return driver
}

//--------------------------------mockDriverEventListener-----------------------

type mockDriverEventListener struct {
	eventHandlerFunc func(hardware.DriverEvent)
}

func (this *mockDriverEventListener) HandleDriverEvent(e hardware.DriverEvent) {
	if this.eventHandlerFunc != nil {
		this.eventHandlerFunc(e)
	}
}
