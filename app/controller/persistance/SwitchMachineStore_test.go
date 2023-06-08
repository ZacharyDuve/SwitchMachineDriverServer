package persistance

import (
	"testing"

	"github.com/ZacharyDuve/SwitchMachineDriverServer/app/controller/switchmachine"
)

func TestNewSwitchMachineStoreStartsNonNil(t *testing.T) {
	smStore := NewSwitchMachineStore()
	if smStore.GetAll() == nil {
		t.FailNow()
	}
}

func TestNewSwitchMachineStoreStartsEmpty(t *testing.T) {
	smStore := NewSwitchMachineStore()
	if len(smStore.GetAll()) != 0 {
		t.FailNow()
	}
}

func TestSwitchMachineStoreHasOneSwitchMachineStateAfterOneIsAdded(t *testing.T) {
	smStore := NewSwitchMachineStore()
	smStore.AddSwitchMachine(getSampleSwitchMachineState())

	if len(smStore.GetAll()) != 1 {
		t.Fail()
	}
}

func TestSwitchMachineStoreAddingUniqueSwitchMachineReturnsNoError(t *testing.T) {
	smStore := NewSwitchMachineStore()
	err := smStore.AddSwitchMachine(getSampleSwitchMachineState())

	if err != nil {
		t.Error()
	}
}

func TestSwitchMachineStoreReturnsErrorWhenAddingADuplicatSwitchMachine(t *testing.T) {
	smStore := NewSwitchMachineStore()
	newState := getSampleSwitchMachineState()
	smStore.AddSwitchMachine(newState)
	err := smStore.AddSwitchMachine(newState)

	if err == nil {
		t.Error()
	}
}

func TestHasSwitchMachineReturnsFalseWhenSwitchMachineWithIdDoesNotExistInStore(t *testing.T) {
	idUnderTest := switchmachine.Id(1)
	smStore := NewSwitchMachineStore()
	if smStore.HasSwitchMachine(idUnderTest) {
		t.Fail()
	}
}

func TestHasSwitchMachineReturnsTrueWhenSwitchMachineWithIdDoesExistInStore(t *testing.T) {
	idUnderTest := switchmachine.Id(1)
	sm := switchmachine.NewState(idUnderTest, switchmachine.Position0, switchmachine.MotorStateIdle, switchmachine.GPIOOFF, switchmachine.GPIOOFF)
	smStore := NewSwitchMachineStore()
	smStore.AddSwitchMachine(sm)
	if !smStore.HasSwitchMachine(idUnderTest) {
		t.Fail()
	}
}

func TestUpdateSwitchMachineReturnsErrorWhenSwitchMachineWithIdDoesNotExistInStore(t *testing.T) {
	sm := getSampleSwitchMachineState()
	smStore := NewSwitchMachineStore()
	if smStore.UpdateSwitchMachine(sm) == nil {
		t.Fail()
	}
}

func TestUpdateSwitchMachineDoesNotReturnErrorWhenSwitchMachineDoesExistInStore(t *testing.T) {
	sm := getSampleSwitchMachineState()
	smStore := NewSwitchMachineStore()
	smStore.AddSwitchMachine(sm)
	if smStore.UpdateSwitchMachine(sm) != nil {
		t.Fail()
	}
}

func TestGetSwitchMachineAfterUpdatingReturnsNewState(t *testing.T) {
	idUnderTest := switchmachine.Id(1)
	origState := switchmachine.NewState(idUnderTest, switchmachine.Position0, switchmachine.MotorStateIdle, switchmachine.GPIOOFF, switchmachine.GPIOOFF)

	smStore := NewSwitchMachineStore()
	smStore.AddSwitchMachine(origState)
	nextState := switchmachine.NewState(idUnderTest, switchmachine.Position1, switchmachine.MotorStateIdle, switchmachine.GPIOOFF, switchmachine.GPIOOFF)
	smStore.UpdateSwitchMachine(nextState)
	lookedUpState := smStore.GetSwitchMachineById(idUnderTest)
	if !switchmachine.StatesEqual(nextState, lookedUpState) {
		t.Fail()
	}
}

func TestRemoveSwitchMachineReturnsErrorIfSwitchMachineWithIdDoesNotExist(t *testing.T) {
	idUnderTest := switchmachine.Id(1)
	smStore := NewSwitchMachineStore()

	if _, err := smStore.RemoveSwitchMachine(idUnderTest); err == nil {
		t.Fail()
	}
}

func TestRemoveSwitchMachineDoesNotReturnErrorWhenSwitchMachineWithIdDoesExist(t *testing.T) {
	idUnderTest := switchmachine.Id(1)
	smStore := NewSwitchMachineStore()
	sm := switchmachine.NewState(idUnderTest, switchmachine.Position1, switchmachine.MotorStateIdle, switchmachine.GPIOOFF, switchmachine.GPIOOFF)
	smStore.AddSwitchMachine(sm)
	if _, err := smStore.RemoveSwitchMachine(idUnderTest); err != nil {
		t.Fail()
	}
}

func TestRemoveSwitchMachineRemovesSwitchMachineWithId(t *testing.T) {
	idUnderTest := switchmachine.Id(1)
	smStore := NewSwitchMachineStore()
	sm := switchmachine.NewState(idUnderTest, switchmachine.Position1, switchmachine.MotorStateIdle, switchmachine.GPIOOFF, switchmachine.GPIOOFF)
	smStore.AddSwitchMachine(sm)
	smStore.RemoveSwitchMachine(sm.Id())
	if smStore.HasSwitchMachine(sm.Id()) {
		t.Fail()
	}

}

func getSampleSwitchMachineState() switchmachine.State {
	return switchmachine.NewState(switchmachine.Id(0), switchmachine.Position0, switchmachine.MotorStateIdle, switchmachine.GPIOOFF, switchmachine.GPIOOFF)
}
