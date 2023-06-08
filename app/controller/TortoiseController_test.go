package controller

import (
	"testing"

	"github.com/ZacharyDuve/SwitchMachineDriverServer/app/controller/hardware"
	"github.com/ZacharyDuve/SwitchMachineDriverServer/app/controller/switchmachine"
)

func TestUpdateSwitchMachineReturnsErrorIfSwitchMachineIsUnknown(t *testing.T) {
	c := newTortoiseController()
	idUnderTest := switchmachine.Id(0)
	sm := switchmachine.NewState(idUnderTest, switchmachine.Position0, switchmachine.MotorStateIdle, switchmachine.GPIOOFF, switchmachine.GPIOOFF)
	if c.UpdateSwitchMachine(sm) == nil {
		t.Fail()
	}
}

func TestUpdateSwitchMachineDoesNotReturnErrorIfSwitchMachineAttached(t *testing.T) {
	c := newTortoiseController()
	c.driver = &mockHardwareDriver{}
	idUnderTest := switchmachine.Id(1)
	sm := switchmachine.NewState(idUnderTest, switchmachine.Position0, switchmachine.MotorStateIdle, switchmachine.GPIOOFF, switchmachine.GPIOOFF)
	c.HandleDriverEvent(hardware.NewSwitchMachineAddedEvent(idUnderTest, sm))

	if c.UpdateSwitchMachine(sm) != nil {
		t.Fail()
	}
}

func TestUpdateSwitchMachineDoesNotCallDriverUpdateSwitchMachineIfSwitchMachineDoesNotExist(t *testing.T) {
	c := newTortoiseController()
	idUnderTest := switchmachine.Id(1)

	wasDriverUpdateCalled := false
	driver := &mockHardwareDriver{}
	driver.updateSwitchMachineFunc = func(sms switchmachine.State) {
		wasDriverUpdateCalled = true
	}
	c.driver = driver

	sm := switchmachine.NewState(idUnderTest, switchmachine.Position0, switchmachine.MotorStateIdle, switchmachine.GPIOOFF, switchmachine.GPIOOFF)
	c.UpdateSwitchMachine(sm)

	if wasDriverUpdateCalled {
		t.Fail()
	}
}

func TestUpdateSwitchMachineDoesNotCallDriverIfNewStateHasPositionAndGPIOEqualToExistingAndMotorIdle(t *testing.T) {
	c := newTortoiseController()

	wasDriverUpdateCalled := false
	driver := &mockHardwareDriver{}
	driver.updateSwitchMachineFunc = func(sms switchmachine.State) {
		wasDriverUpdateCalled = true
	}
	c.driver = driver
	idUnderTest := switchmachine.Id(1)
	positionUnderTest := switchmachine.Position0
	gpio0 := switchmachine.GPIOOFF
	gpio1 := switchmachine.GPIOOFF

	smOrig := switchmachine.NewState(idUnderTest, positionUnderTest, switchmachine.MotorStateIdle, gpio0, gpio1)

	smNext := switchmachine.NewState(idUnderTest, positionUnderTest, switchmachine.MotorStateIdle, gpio0, gpio1)
	c.UpdateSwitchMachine(smOrig)
	c.UpdateSwitchMachine(smNext)

	if wasDriverUpdateCalled {
		t.Fail()
	}
}

func TestUpdateSwitchMachineDoesNotCallDriverIfNewStateHasPositionAndGPIOEqualToExistingAndMotorBrake(t *testing.T) {
	c := newTortoiseController()

	wasDriverUpdateCalled := false
	driver := &mockHardwareDriver{}
	driver.updateSwitchMachineFunc = func(sms switchmachine.State) {
		wasDriverUpdateCalled = true
	}
	c.driver = driver
	idUnderTest := switchmachine.Id(1)
	positionUnderTest := switchmachine.Position0
	gpio0 := switchmachine.GPIOOFF
	gpio1 := switchmachine.GPIOOFF

	smOrig := switchmachine.NewState(idUnderTest, positionUnderTest, switchmachine.MotorStateBrake, gpio0, gpio1)
	c.HandleDriverEvent(hardware.NewSwitchMachineAddedEvent(smOrig.Id(), smOrig))
	//Here we are changing
	smNext := switchmachine.NewState(idUnderTest, positionUnderTest, switchmachine.MotorStateBrake, gpio0, gpio1)
	//c.SwitchMachineAdded(smOrig)
	c.UpdateSwitchMachine(smNext)

	if wasDriverUpdateCalled {
		t.Fail()
	}
}

func TestThatisMotorRunningToOppositePositionReturnsFalseIfMotorIsCurrentlyIdle(t *testing.T) {
	idUnderTest := switchmachine.Id(0)
	gpio0 := switchmachine.GPIOOFF
	gpio1 := switchmachine.GPIOOFF
	curS := switchmachine.NewState(idUnderTest, switchmachine.Position0, switchmachine.MotorStateIdle, gpio0, gpio1)

	newS := switchmachine.NewState(idUnderTest, switchmachine.Position0, switchmachine.MotorStateToPos0, gpio0, gpio1)

	if isMotorRunningToOppositePosition(newS, curS) {
		t.Fail()
	}
}

func TestThatisMotorRunningToOppositePositionReturnsFalseIfMotorIsCurrentlyBrake(t *testing.T) {
	idUnderTest := switchmachine.Id(0)
	gpio0 := switchmachine.GPIOOFF
	gpio1 := switchmachine.GPIOOFF
	curS := switchmachine.NewState(idUnderTest, switchmachine.Position0, switchmachine.MotorStateBrake, gpio0, gpio1)

	newS := switchmachine.NewState(idUnderTest, switchmachine.Position0, switchmachine.MotorStateToPos0, gpio0, gpio1)

	if isMotorRunningToOppositePosition(newS, curS) {
		t.Fail()
	}
}

func TestThatisMotorRunningToOppositePositionReturnsTrueIfMotorIsCurrentlyRunningTo1IfDesire0(t *testing.T) {
	idUnderTest := switchmachine.Id(0)
	gpio0 := switchmachine.GPIOOFF
	gpio1 := switchmachine.GPIOOFF
	curS := switchmachine.NewState(idUnderTest, switchmachine.Position0, switchmachine.MotorStateToPos1, gpio0, gpio1)

	newS := switchmachine.NewState(idUnderTest, switchmachine.Position0, switchmachine.MotorStateToPos0, gpio0, gpio1)

	if !isMotorRunningToOppositePosition(newS, curS) {
		t.Fail()
	}
}

func TestThatisMotorRunningToOppositePositionReturnsTrueIfMotorIsCurrentlyRunningTo0IfDesire1(t *testing.T) {
	idUnderTest := switchmachine.Id(0)
	gpio0 := switchmachine.GPIOOFF
	gpio1 := switchmachine.GPIOOFF
	curS := switchmachine.NewState(idUnderTest, switchmachine.Position0, switchmachine.MotorStateToPos0, gpio0, gpio1)

	newS := switchmachine.NewState(idUnderTest, switchmachine.Position1, switchmachine.MotorStateToPos0, gpio0, gpio1)

	if !isMotorRunningToOppositePosition(newS, curS) {
		t.Fail()
	}
}

func TestUpdateSwitchMachineCallsDriverIfNewStateHasPositionAndGPIOEqualToExistingButMotorGoingToOpposite(t *testing.T) {
	c := newTortoiseController()

	wasDriverUpdateCalled := false
	driver := &mockHardwareDriver{}
	driver.updateSwitchMachineFunc = func(sms switchmachine.State) {
		wasDriverUpdateCalled = true
	}
	c.driver = driver
	idUnderTest := switchmachine.Id(1)
	gpio0 := switchmachine.GPIOOFF
	gpio1 := switchmachine.GPIOOFF

	smOrig := switchmachine.NewState(idUnderTest, switchmachine.Position0, switchmachine.MotorStateToPos1, gpio0, gpio1)
	c.HandleDriverEvent(hardware.NewSwitchMachineAddedEvent(smOrig.Id(), smOrig))
	smNext := switchmachine.NewState(idUnderTest, switchmachine.Position0, switchmachine.MotorStateToPos0, gpio0, gpio1)
	c.UpdateSwitchMachine(smNext)

	if !wasDriverUpdateCalled {
		t.Fail()
	}
}

func TestThatCallbackSetsMotorToIdleWhilePreservingCurrentState(t *testing.T) {
	driver := &mockHardwareDriver{}
	c := newTortoiseController()
	c.driver = driver

	idUnderTest := switchmachine.Id(12)
	gpio0 := switchmachine.GPIOOn
	gpio1 := switchmachine.GPIOOn
	smOrig := switchmachine.NewState(idUnderTest, switchmachine.Position0, switchmachine.MotorStateToPos1, gpio0, gpio1)
	c.HandleDriverEvent(hardware.NewSwitchMachineAddedEvent(smOrig.Id(), smOrig))
	//smNext := switchmachine.NewState(idUnderTest, switchmachine.Position0, switchmachine.MotorStateToPos0, gpio0, gpio1)
	c.stopMotorCallbackFunc(smOrig.Id(), 0)

	if curS, _ := c.GetSwitchMachineById(idUnderTest); !areUpdateableFieldsEqual(curS, smOrig) || curS.MotorState() != switchmachine.MotorStateIdle {
		t.Fail()
	}
}

type mockHardwareDriver struct {
	updateSwitchMachineFunc func(switchmachine.State)
	closeFunc               func() error
}

func (this *mockHardwareDriver) UpdateSwitchMachine(sm switchmachine.State) {
	if this.updateSwitchMachineFunc != nil {
		this.updateSwitchMachineFunc(sm)
	}
}

func (this *mockHardwareDriver) Close() error {
	if this.closeFunc != nil {
		return this.closeFunc()
	}
	return nil
}

func (this *mockHardwareDriver) Start(hardware.DriverEventListener) {

}
