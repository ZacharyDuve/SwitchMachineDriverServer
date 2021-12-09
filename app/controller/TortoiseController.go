package controller

import (
	"errors"
	"log"
	"time"

	"github.com/ZacharyDuve/SwitchMachineDriverServer/app/controller/hardware"
	"github.com/ZacharyDuve/SwitchMachineDriverServer/app/controller/hardware/tortoise"
	"github.com/ZacharyDuve/SwitchMachineDriverServer/app/controller/model"
	"github.com/ZacharyDuve/SwitchMachineDriverServer/app/controller/persistance"
)

const (
	switchMachineNotExistErrorMessage string        = "Switch Machine with matching Id does not exist"
	defaultMotorRunTime               time.Duration = time.Second * 4
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
	controller := newTortoiseController()
	var err error
	controller.driver, err = tortoise.NewPiTortoiseControllerDriver(controller)

	if err != nil {
		panic(err)
	}

	return controller
}

func NewTortoiseControllerWithMockDriver(txDataOut, rxDataIn []byte) TortoiseController {
	controller := newTortoiseController()
	var err error
	controller.driver, err = tortoise.NewMockTortoiseControllerDriver(controller, txDataOut, rxDataIn)

	if err != nil {
		panic(err)
	}

	return controller
}

//Function that returns a *tortoiseControllerImpl so that we can use its functions in tests
func newTortoiseController() *tortoiseControllerImpl {
	controller := &tortoiseControllerImpl{}
	controller.existingSMStates = persistance.NewSwitchMachineStore()

	return controller
}

func (this *tortoiseControllerImpl) UpdateSwitchMachine(req SwitchMachineUpdateRequest) error {
	var err error

	if !this.existingSMStates.HasSwitchMachine(req.Id()) {
		err = errors.New(switchMachineNotExistErrorMessage)
	} else {
		//curSMState := this.existingSMStates.GetSwitchMachineById(req.Id())
		//if isMotorRunningToOpposingPosition(curSMState, req) {
		if true {
			newState := createNewStateFromRequest(req)
			this.driver.UpdateSwitchMachine(newState)
			this.createStopMotorCallback(newState.Id())
		}
	}

	return err
}

func (this *tortoiseControllerImpl) createStopMotorCallback(id model.SwitchMachineId) {
	go func() {
		time.Sleep(defaultMotorRunTime)
		stateBeforeMotorStop := this.existingSMStates.GetSwitchMachineById(id)
		stoppedMotorState := model.NewSwitchMachineState(id,
			stateBeforeMotorStop.Position(),
			model.MotorStateIdle,
			stateBeforeMotorStop.GPIO0State(),
			stateBeforeMotorStop.GPIO1State())
		log.Println("Stopping motor with id:", id)
		this.driver.UpdateSwitchMachine(stoppedMotorState)
	}()
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

func IsSwitchMachineNotExistError(err error) bool {
	return err.Error() == switchMachineNotExistErrorMessage
}
