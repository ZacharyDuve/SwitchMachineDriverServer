package controller

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/ZacharyDuve/SwitchMachineDriverServer/app/controller/event"
	"github.com/ZacharyDuve/SwitchMachineDriverServer/app/controller/hardware"
	"github.com/ZacharyDuve/SwitchMachineDriverServer/app/controller/persistance"
	"github.com/ZacharyDuve/SwitchMachineDriverServer/app/controller/switchmachine"
)

const (
	switchMachineNotExistErrorMessage string        = "Switch Machine with matching Id %d does not exist"
	defaultMotorRunTime               time.Duration = time.Second * 4
)

type TortoiseController interface {
	UpdateSwitchMachine(switchmachine.State) error
	GetSwitchMachines() []switchmachine.State
	GetSwitchMachineById(id switchmachine.Id) (switchmachine.State, error)
	SetSwitchMachineEventListenerFunc(func(event.SwitchMachineEvent))
	HandleDriverEvent(dE hardware.DriverEvent)
}

type tortoiseControllerImpl struct {
	driver              hardware.Driver
	existingSMStates    persistance.SwitchMachineStore
	smEventListenerFunc func(event.SwitchMachineEvent)
}

//Wrapping the internal testable call as an external facing interface to restrict functions
func NewTortoiseController(driver hardware.Driver) TortoiseController {
	if driver == nil {
		panic("driver is required for NewTortoiseController")
	}
	controller := newTortoiseController()

	controller.driver = driver
	driver.Start(controller)

	return controller
}

//Function that returns a *tortoiseControllerImpl so that we can use its functions in tests
func newTortoiseController() *tortoiseControllerImpl {
	controller := &tortoiseControllerImpl{}
	controller.existingSMStates = persistance.NewSwitchMachineStore()

	return controller
}

func (this *tortoiseControllerImpl) UpdateSwitchMachine(requestState switchmachine.State) error {
	log.Println("tortoiseControllerImpl-UpdateSwitchMachine called")
	var err error
	curState := this.existingSMStates.GetSwitchMachineById(requestState.Id())
	log.Println("requestState:", switchmachine.StateToString(requestState))
	log.Println("curState:", switchmachine.StateToString(curState))
	if curState == nil {
		//We don't have a switchmachine for this id
		err = newSwitchMachineNotExistError(requestState.Id())
	} else if !areUpdateableFieldsEqual(curState, requestState) || isMotorRunningToOppositePosition(requestState, curState) {
		//Figure out if we need to set a new motor state
		newMotorState := switchmachine.MotorStateIdle
		if curState.Position() != requestState.Position() {
			if requestState.Position() == switchmachine.Position0 {
				newMotorState = switchmachine.MotorStateToPos0
			} else if requestState.Position() == switchmachine.Position1 {
				newMotorState = switchmachine.MotorStateToPos1
			}
		}

		newState := switchmachine.NewState(requestState.Id(), curState.Position(), newMotorState, requestState.GPIO0State(), requestState.GPIO1State())
		log.Println("newState:", switchmachine.StateToString(newState))
		this.driver.UpdateSwitchMachine(newState)
		//If we just told it to change position then we need to stop it at some point
		if curState.MotorState() != newState.MotorState() {
			this.createStopMotorCallback(curState.Id())
		}
		if !areGPIOEqual(curState, newState) || curState.MotorState() != newState.MotorState() {
			this.existingSMStates.UpdateSwitchMachine(newState)
			this.sendSMEventToListener(event.NewSwitchMachineUpdatedEvent(newState))
		}
	}
	return err
}

//Trying to cover the case were we are where we want but we are moving away from it
func isMotorRunningToOppositePosition(newS, curS switchmachine.State) bool {
	return newS.Position() == switchmachine.Position0 && curS.MotorState() == switchmachine.MotorStateToPos1 ||
		newS.Position() == switchmachine.Position1 && curS.MotorState() == switchmachine.MotorStateToPos0
}

func (this *tortoiseControllerImpl) GetSwitchMachines() []switchmachine.State {
	return this.existingSMStates.GetAll()
}

func (this *tortoiseControllerImpl) GetSwitchMachineById(id switchmachine.Id) (switchmachine.State, error) {
	sm := this.existingSMStates.GetSwitchMachineById(id)
	var err error
	if sm == nil {
		err = newSwitchMachineNotExistError(id)
	}

	return sm, err
}

func (this *tortoiseControllerImpl) createStopMotorCallback(id switchmachine.Id) {
	go this.stopMotorCallbackFunc(id, defaultMotorRunTime)
}

func (this *tortoiseControllerImpl) stopMotorCallbackFunc(id switchmachine.Id, delay time.Duration) {
	if delay > 0 {
		time.Sleep(delay)
	}
	stateBeforeMotorStop := this.existingSMStates.GetSwitchMachineById(id)
	stoppedMotorState := switchmachine.NewState(id,
		stateBeforeMotorStop.Position(),
		switchmachine.MotorStateIdle,
		stateBeforeMotorStop.GPIO0State(),
		stateBeforeMotorStop.GPIO1State())
	if this.existingSMStates.HasSwitchMachine(id) {
		this.driver.UpdateSwitchMachine(stoppedMotorState)
		this.existingSMStates.UpdateSwitchMachine(stoppedMotorState)
		this.sendSMEventToListener(event.NewSwitchMachineUpdatedEvent(stoppedMotorState))
	}
}

func (this *tortoiseControllerImpl) SetSwitchMachineEventListenerFunc(smEventListenFunc func(event.SwitchMachineEvent)) {
	this.smEventListenerFunc = smEventListenFunc
}

func (this *tortoiseControllerImpl) HandleDriverEvent(dE hardware.DriverEvent) {
	var err error
	var e event.SwitchMachineEvent
	if dE.Type() == hardware.SwitchMachineAdded {
		err = this.existingSMStates.AddSwitchMachine(dE.State())
		if err == nil {
			e = event.NewSwitchMachineAddedEvent(dE.State())
		}
	} else if dE.Type() == hardware.SwitchMachinePositionChanged {
		err = this.existingSMStates.UpdateSwitchMachine(dE.State())
		if err == nil {
			e = event.NewSwitchMachineUpdatedEvent(dE.State())
		}
	} else if dE.Type() == hardware.SwitchMachineRemoved {
		var lastState switchmachine.State
		lastState, err = this.existingSMStates.RemoveSwitchMachine(dE.Id())
		if err == nil {
			log.Println("Reseting output for switchmachine with id:", dE.Id())
			this.driver.UpdateSwitchMachine(switchmachine.NewState(dE.Id(), lastState.Position(), switchmachine.MotorStateIdle, switchmachine.GPIOOFF, switchmachine.GPIOOFF))
			e = event.NewSwitchMachineRemovedEvent(lastState)
		}
	}

	if err != nil {
		panic(err)
	}

	this.sendSMEventToListener(e)
}

func (this *tortoiseControllerImpl) sendSMEventToListener(sme event.SwitchMachineEvent) {
	if this.smEventListenerFunc != nil {
		this.smEventListenerFunc(sme)
	}
}

func newSwitchMachineNotExistError(id switchmachine.Id) error {
	return errors.New(fmt.Sprintf(switchMachineNotExistErrorMessage, id))
}

func IsSwitchMachineNotExistError(err error) bool {
	return err.Error() == switchMachineNotExistErrorMessage
}

func areUpdateableFieldsEqual(s0, s1 switchmachine.State) bool {
	return areGPIOEqual(s0, s1) &&
		s0.Position() == s1.Position()
}

func areGPIOEqual(s0, s1 switchmachine.State) bool {
	return s0.GPIO0State() == s1.GPIO0State() &&
		s0.GPIO1State() == s1.GPIO1State()
}
