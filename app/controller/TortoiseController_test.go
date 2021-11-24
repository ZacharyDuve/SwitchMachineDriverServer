package controller

import (
	"testing"

	"github.com/ZacharyDuve/SwitchMachineDriverServer/app/controller/model"
)

func TestUpdateSwitchMachineReturnsErrorIfSwitchMachineIsUnknown(t *testing.T) {
	c := newTortoiseController()
	idUnderTest := model.SwitchMachineId(1)
	sm := NewSwitchMachineUpdateRequest(idUnderTest, model.Position0, model.GPIOOFF, model.GPIOOFF)
	if c.UpdateSwitchMachine(sm) == nil {
		t.Fail()
	}
}

func TestUpdateSwitchMachineDoesNotReturnErrorIfSwitchMachineAttached(t *testing.T) {
	c := newTortoiseController()
	idUnderTest := model.SwitchMachineId(1)
	sm := model.NewSwitchMachineState(idUnderTest, model.Position0, model.MotorStateIdle, model.GPIOOFF, model.GPIOOFF)
	//SwitchMachines are added via the driver calling this function when the switchmachine is connected
	c.SwitchMachineAdded(sm)

	if c.UpdateSwitchMachine(sm) != nil {
		t.Fail()
	}
}

func TestUpdateSwitchMachineDoesNotCallDriverUpdateSwitchMachineIfSwitchMachineDoesNotExist(t *testing.T) {
	c := newTortoiseController()
	idUnderTest := model.SwitchMachineId(1)

	wasDriverUpdateCalled := false
	driver := &mockHardwareDriver{}
	driver.updateSwitchMachineFunc = func(sms model.SwitchMachineState) error {
		wasDriverUpdateCalled = true
		return nil
	}
	c.driver = driver

	sm := NewSwitchMachineUpdateRequest(idUnderTest, model.Position0, model.GPIOOFF, model.GPIOOFF)
	c.UpdateSwitchMachine(sm)

	if wasDriverUpdateCalled {
		t.Fail()
	}
}

func TestUpdateSwitchMachineDoesNotCallDriverIfNewStateHasPositionAndGPIOEqualToExistingAndMotorIdle(t *testing.T) {
	c := newTortoiseController()

	wasDriverUpdateCalled := false
	driver := &mockHardwareDriver{}
	driver.updateSwitchMachineFunc = func(sms model.SwitchMachineState) error {
		wasDriverUpdateCalled = true
		return nil
	}
	c.driver = driver
	idUnderTest := model.SwitchMachineId(1)
	positionUnderTest := model.Position0
	gpio0 := model.GPIOOFF
	gpio1 := model.GPIOOFF

	smOrig := model.NewSwitchMachineState(idUnderTest, positionUnderTest, model.MotorStateIdle, gpio0, gpio1)

	smNext := NewSwitchMachineUpdateRequest(idUnderTest, positionUnderTest, gpio0, gpio1)
	c.SwitchMachineAdded(smOrig)
	c.UpdateSwitchMachine(smNext)

	if wasDriverUpdateCalled {
		t.Fail()
	}
}

func TestUpdateSwitchMachineDoesNotCallDriverIfNewStateHasPositionAndGPIOEqualToExistingAndMotorBrake(t *testing.T) {
	c := newTortoiseController()

	wasDriverUpdateCalled := false
	driver := &mockHardwareDriver{}
	driver.updateSwitchMachineFunc = func(sms model.SwitchMachineState) error {
		wasDriverUpdateCalled = true
		return nil
	}
	c.driver = driver
	idUnderTest := model.SwitchMachineId(1)
	positionUnderTest := model.Position0
	gpio0 := model.GPIOOFF
	gpio1 := model.GPIOOFF

	smOrig := model.NewSwitchMachineState(idUnderTest, positionUnderTest, model.MotorStateBrake, gpio0, gpio1)
	//Here we are changing
	smNext := NewSwitchMachineUpdateRequest(idUnderTest, positionUnderTest, gpio0, gpio1)
	c.SwitchMachineAdded(smOrig)
	c.UpdateSwitchMachine(smNext)

	if wasDriverUpdateCalled {
		t.Fail()
	}
}

func TestUpdateSwitchMachineCallsDriverIfNewStateHasPositionAndGPIOEqualToExistingButMotorGoingToOpposite(t *testing.T) {
	c := newTortoiseController()

	wasDriverUpdateCalled := false
	driver := &mockHardwareDriver{}
	driver.updateSwitchMachineFunc = func(sms model.SwitchMachineState) error {
		wasDriverUpdateCalled = true
		return nil
	}
	c.driver = driver
	idUnderTest := model.SwitchMachineId(1)
	gpio0 := model.GPIOOFF
	gpio1 := model.GPIOOFF

	smOrig := model.NewSwitchMachineState(idUnderTest, model.Position0, model.MotorStateToPos1, gpio0, gpio1)
	//Here we are changing
	smNext := NewSwitchMachineUpdateRequest(idUnderTest, model.Position0, gpio0, gpio1)
	c.SwitchMachineAdded(smOrig)
	c.UpdateSwitchMachine(smNext)

	if !wasDriverUpdateCalled {
		t.Fail()
	}
}

func TestIsMotorRunningToOpposingPositionReturnsTrueIfExistingHasMotorRunningToPosition1AndNewMovesToPosition0(t *testing.T) {
	idUnderTest := model.SwitchMachineId(1)
	existing := model.NewSwitchMachineState(idUnderTest, model.Position0, model.MotorStateToPos1, model.GPIOOFF, model.GPIOOFF)
	req := NewSwitchMachineUpdateRequest(idUnderTest, model.Position0, model.GPIOOFF, model.GPIOOFF)
	if !isMotorRunningToOpposingPosition(existing, req) {
		t.Fail()
	}
}

func TestIsMotorRunningToOpposingPositionReturnsTrueIfExistingHasMotorRunningToPosition0AndNewMovesToPosition1(t *testing.T) {
	idUnderTest := model.SwitchMachineId(1)
	existing := model.NewSwitchMachineState(idUnderTest, model.Position0, model.MotorStateToPos0, model.GPIOOFF, model.GPIOOFF)
	req := NewSwitchMachineUpdateRequest(idUnderTest, model.Position1, model.GPIOOFF, model.GPIOOFF)
	if !isMotorRunningToOpposingPosition(existing, req) {
		t.Fail()
	}
}

func TestIsMotorRunningToOpposingPositionReturnsFalseIfExistingHasMotorRunningToPositionSameAsNewMovesToPosition(t *testing.T) {
	idUnderTest := model.SwitchMachineId(1)
	existing := model.NewSwitchMachineState(idUnderTest, model.Position0, model.MotorStateToPos0, model.GPIOOFF, model.GPIOOFF)
	req := NewSwitchMachineUpdateRequest(idUnderTest, model.Position0, model.GPIOOFF, model.GPIOOFF)
	if isMotorRunningToOpposingPosition(existing, req) {
		t.Fail()
	}
}

type mockHardwareDriver struct {
	updateSwitchMachineFunc func(model.SwitchMachineState) error
	closeFunc               func() error
}

func (this *mockHardwareDriver) UpdateSwitchMachine(sm model.SwitchMachineState) error {
	if this.updateSwitchMachineFunc != nil {
		return this.updateSwitchMachineFunc(sm)
	}
	return nil
}

func (this *mockHardwareDriver) Close() error {
	if this.closeFunc != nil {
		return this.closeFunc()
	}
	return nil
}
