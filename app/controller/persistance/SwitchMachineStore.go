package persistance

import (
	"errors"
	"sync"

	"github.com/ZacharyDuve/SwitchMachineDriverServer/app/controller/switchmachine"
)

const (
	dupSMMessage     string = "Switch Machine with same id already exists in store"
	missingSMMessage string = "Switch Machine with id is not in store"
)

type SwitchMachineStore interface {
	AddSwitchMachine(switchmachine.State) error
	HasSwitchMachine(switchmachine.Id) bool
	GetSwitchMachineById(switchmachine.Id) switchmachine.State
	GetAll() []switchmachine.State
	RemoveSwitchMachine(switchmachine.Id) (switchmachine.State, error)
	UpdateSwitchMachine(switchmachine.State) error
}

type switchMachineStoreImpl struct {
	switchMachines map[switchmachine.Id]switchmachine.State
	rwLock         *sync.RWMutex
}

func NewSwitchMachineStore() SwitchMachineStore {
	smStore := &switchMachineStoreImpl{}
	smStore.rwLock = &sync.RWMutex{}
	smStore.switchMachines = make(map[switchmachine.Id]switchmachine.State)
	return smStore
}

func (this *switchMachineStoreImpl) AddSwitchMachine(newSwitchMachine switchmachine.State) error {
	var err error
	if this.HasSwitchMachine(newSwitchMachine.Id()) {
		err = newAlreadyHaveSwitchMachineError()
	} else {
		this.rwLock.Lock()
		this.switchMachines[newSwitchMachine.Id()] = newSwitchMachine
		this.rwLock.Unlock()
	}

	return err
}

func (this *switchMachineStoreImpl) HasSwitchMachine(sMachineId switchmachine.Id) bool {
	this.rwLock.RLock()
	_, ok := this.switchMachines[sMachineId]
	this.rwLock.RUnlock()
	return ok
}

func (this *switchMachineStoreImpl) GetSwitchMachineById(sMachineId switchmachine.Id) switchmachine.State {
	this.rwLock.RLock()
	sm, _ := this.switchMachines[sMachineId]
	this.rwLock.RUnlock()
	return sm
}

func (this *switchMachineStoreImpl) RemoveSwitchMachine(smId switchmachine.Id) (switchmachine.State, error) {
	var err error
	var lastState switchmachine.State

	if !this.HasSwitchMachine(smId) {
		err = newDoesNotContainSwitchMacineWithIdError()
	} else {
		this.rwLock.Lock()
		lastState, _ = this.switchMachines[smId]
		delete(this.switchMachines, smId)
		this.rwLock.Unlock()
	}

	return lastState, err
}

func (this *switchMachineStoreImpl) UpdateSwitchMachine(sm switchmachine.State) error {
	var err error
	this.rwLock.RLock()
	contains := this.HasSwitchMachine(sm.Id())
	this.rwLock.RUnlock()
	if !contains {
		err = newDoesNotContainSwitchMacineWithIdError()
	} else {
		this.rwLock.Lock()
		this.switchMachines[sm.Id()] = sm
		this.rwLock.Unlock()
	}
	return err
}

func (this *switchMachineStoreImpl) GetAll() []switchmachine.State {
	this.rwLock.RLock()
	allSM := make([]switchmachine.State, 0, len(this.switchMachines))
	for _, curState := range this.switchMachines {
		allSM = append(allSM, curState)
	}
	this.rwLock.RUnlock()
	return allSM
}

func newAlreadyHaveSwitchMachineError() error {
	return errors.New(dupSMMessage)
}

func newDoesNotContainSwitchMacineWithIdError() error {
	return errors.New(missingSMMessage)
}
