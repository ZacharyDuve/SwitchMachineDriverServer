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
	SetSwitchMachinePosition(id model.SwitchMachineId, pos model.SwitchMachinePosition) error
	SetSwitchMachineGPIO(id model.SwitchMachineId, gpio0, gpio1 model.GPIOState) error
	GetSwitchMachines() []model.SwitchMachineState
	GetSwitchMachineById(id model.SwitchMachineId) (model.SwitchMachineState, error)
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

func (this *tortoiseControllerImpl) SetSwitchMachinePosition(id model.SwitchMachineId, pos model.SwitchMachinePosition) error {
	var err error

	if !this.existingSMStates.HasSwitchMachine(id) {
		err = newSwitchMachineNotExistError()
	} else {
		log.Println("Updating position for switch machine", id)
		curState := this.existingSMStates.GetSwitchMachineById(id)
		log.Println("curState:", curState)
		newState := model.NewSwitchMachineState(id, pos, getMotorStateForToPosition(pos), curState.GPIO0State(), curState.GPIO1State())
		log.Println("newState:", newState)
		this.driver.UpdateSwitchMachine(newState)
		this.createStopMotorCallback(id)
		this.existingSMStates.UpdateSwitchMachine(newState)
	}

	return err
}

func getMotorStateForToPosition(pos model.SwitchMachinePosition) model.SwitchMachineMotorState {
	if pos == model.Position0 {
		return model.MotorStateToPos0
	} else if pos == model.Position1 {
		return model.MotorStateToPos1
	} else {
		return model.MotorStateIdle
	}
}

func (this *tortoiseControllerImpl) SetSwitchMachineGPIO(id model.SwitchMachineId, gpio0, gpio1 model.GPIOState) error {
	var err error
	if !this.existingSMStates.HasSwitchMachine(id) {
		err = newSwitchMachineNotExistError()
	} else {
		log.Println("Updating GPIO for switch machine", id)
		curState := this.existingSMStates.GetSwitchMachineById(id)
		log.Println("curState:", curState)
		newState := model.NewSwitchMachineState(id, curState.Position(), curState.MotorState(), gpio0, gpio1)

		log.Println("newState:", newState)

		this.driver.UpdateSwitchMachine(newState)

		this.existingSMStates.UpdateSwitchMachine(newState)
	}

	return err
}

func (this *tortoiseControllerImpl) GetSwitchMachines() []model.SwitchMachineState {
	return this.existingSMStates.GetAll()
}

func (this *tortoiseControllerImpl) GetSwitchMachineById(id model.SwitchMachineId) (model.SwitchMachineState, error) {
	sm := this.existingSMStates.GetSwitchMachineById(id)
	var err error
	if sm == nil {
		err = newSwitchMachineNotExistError()
	}

	return sm, err
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
		this.existingSMStates.UpdateSwitchMachine(stoppedMotorState)
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

func newSwitchMachineNotExistError() error {
	return errors.New(switchMachineNotExistErrorMessage)
}

func IsSwitchMachineNotExistError(err error) bool {
	return err.Error() == switchMachineNotExistErrorMessage
}
