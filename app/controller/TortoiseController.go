package controller

import (
	"errors"
	"os"

	"github.com/ZacharyDuve/SwitchMachineDriverServer/app/controller/hardware"
	"github.com/ZacharyDuve/SwitchMachineDriverServer/app/controller/hardware/tortoise"
	"github.com/ZacharyDuve/SwitchMachineDriverServer/app/controller/model"
	"github.com/ZacharyDuve/SwitchMachineDriverServer/app/controller/persistance"
)

const (
	switchMachineNotExistErrorMessage string = "Switch Machine with matching Id does not exist"
)

type TortoiseController interface {
	UpdateSwitchMachine(SwitchMachineUpdateRequest) error
}

type tortoiseControllerImpl struct {
	driver           hardware.SwitchMachineDriver
	existingSMStates persistance.SwitchMachineStore
}

//Wrapping the internal testable call as an external facing interface to restrict functions
func NewTortoiseController() TortoiseController {
	return newTortoiseController()
}

//Function that returns a *tortoiseControllerImpl so that we can use its functions in tests
func newTortoiseController() *tortoiseControllerImpl {
	controller := &tortoiseControllerImpl{}
	controller.existingSMStates = persistance.NewSwitchMachineStore()
	var err error
	if os.Getenv("environment") == "production" {
		controller.driver, err = tortoise.NewPiTortoiseControllerDriver(controller)
	} else {
		controller.driver, err = tortoise.NewMockTortoiseControllerDriver(controller)
	}

	if err != nil {
		panic(err)
	}

	return controller
}

func (this *tortoiseControllerImpl) UpdateSwitchMachine(req SwitchMachineUpdateRequest) error {
	var err error

	if !this.existingSMStates.HasSwitchMachine(req.Id()) {
		err = errors.New(switchMachineNotExistErrorMessage)
	} else {
		curSMState := this.existingSMStates.GetSwitchMachineById(req.Id())
		if isMotorRunningToOpposingPosition(curSMState, req) {
			newState := createNewStateFromRequest(req)
			this.driver.UpdateSwitchMachine(newState)
		}
	}

	return err
}

func isMotorRunningToOpposingPosition(existingState model.SwitchMachineState, req SwitchMachineUpdateRequest) bool {
	return existingState.MotorState() == model.MotorStateToPos0 && req.Position() == model.Position1 ||
		existingState.MotorState() == model.MotorStateToPos1 && req.Position() == model.Position0
}

func createNewStateFromRequest(req SwitchMachineUpdateRequest) model.SwitchMachineState {
	var motorState model.SwitchMachineMotorState

	if req.Position() == model.Position0 {
		motorState = model.MotorStateToPos0
	} else {
		motorState = model.MotorStateToPos1
	}
	//Position doesn't matter as it shouldn't be used in requesting to the driver
	return model.NewSwitchMachineState(req.Id(), req.Position(), motorState, req.GPIO0State(), req.GPIO1State())
}

func (this *tortoiseControllerImpl) SwitchMachineAdded(sm model.SwitchMachineState) {
	err := this.existingSMStates.AddSwitchMachine(sm)

	if err != nil {
		panic(err)
	}
}

func (this *tortoiseControllerImpl) SwitchMachineUpdated(sm model.SwitchMachineState) {
	err := this.existingSMStates.UpdateSwitchMachine(sm)

	if err != nil {
		panic(err)
	}
}

func (this *tortoiseControllerImpl) SwitchMachineRemoved(smId model.SwitchMachineId) {
	err := this.existingSMStates.RemoveSwitchMachine(smId)

	if err != nil {
		panic(err)
	}
}
