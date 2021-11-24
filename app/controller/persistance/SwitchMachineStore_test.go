package persistance

import (
	"testing"

	"github.com/ZacharyDuve/SwitchMachineDriverServer/app/controller/model"
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
	idUnderTest := model.SwitchMachineId(1)
	smStore := NewSwitchMachineStore()
	if smStore.HasSwitchMachine(idUnderTest) {
		t.Fail()
	}
}

func TestHasSwitchMachineReturnsTrueWhenSwitchMachineWithIdDoesExistInStore(t *testing.T) {
	idUnderTest := model.SwitchMachineId(1)
	sm := model.NewSwitchMachineState(idUnderTest, model.Position0, model.MotorStateIdle, model.GPIOOFF, model.GPIOOFF)
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
	idUnderTest := model.SwitchMachineId(1)
	origState := model.NewSwitchMachineState(idUnderTest, model.Position0, model.MotorStateIdle, model.GPIOOFF, model.GPIOOFF)

	smStore := NewSwitchMachineStore()
	smStore.AddSwitchMachine(origState)
	nextState := model.NewSwitchMachineState(idUnderTest, model.Position1, model.MotorStateIdle, model.GPIOOFF, model.GPIOOFF)
	smStore.UpdateSwitchMachine(nextState)
	lookedUpState := smStore.GetSwitchMachineById(idUnderTest)
	if !model.SwitchMachineStatesEqual(nextState, lookedUpState) {
		t.Fail()
	}
}

func TestRemoveSwitchMachineReturnsErrorIfSwitchMachineWithIdDoesNotExist(t *testing.T) {
	idUnderTest := model.SwitchMachineId(1)
	smStore := NewSwitchMachineStore()

	if smStore.RemoveSwitchMachine(idUnderTest) == nil {
		t.Fail()
	}
}

func TestRemoveSwitchMachineDoesNotReturnErrorWhenSwitchMachineWithIdDoesExist(t *testing.T) {
	idUnderTest := model.SwitchMachineId(1)
	smStore := NewSwitchMachineStore()
	sm := model.NewSwitchMachineState(idUnderTest, model.Position1, model.MotorStateIdle, model.GPIOOFF, model.GPIOOFF)
	smStore.AddSwitchMachine(sm)
	if smStore.RemoveSwitchMachine(idUnderTest) != nil {
		t.Fail()
	}
}

func TestRemoveSwitchMachineRemovesSwitchMachineWithId(t *testing.T) {
	idUnderTest := model.SwitchMachineId(1)
	smStore := NewSwitchMachineStore()
	sm := model.NewSwitchMachineState(idUnderTest, model.Position1, model.MotorStateIdle, model.GPIOOFF, model.GPIOOFF)
	smStore.AddSwitchMachine(sm)
	smStore.RemoveSwitchMachine(sm.Id())
	if smStore.HasSwitchMachine(sm.Id()) {
		t.Fail()
	}

}

func getSampleSwitchMachineState() model.SwitchMachineState {
	return model.NewSwitchMachineState(model.SwitchMachineId(0), model.Position0, model.MotorStateIdle, model.GPIOOFF, model.GPIOOFF)
}
