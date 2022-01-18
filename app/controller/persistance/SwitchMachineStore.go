package persistance

import (
	"errors"
	"sync"

	"github.com/ZacharyDuve/SwitchMachineDriverServer/app/controller/model"
)

const (
	dupSMMessage     string = "Switch Machine with same id already exists in store"
	missingSMMessage string = "Switch Machine with id is not in store"
)

type SwitchMachineStore interface {
	AddSwitchMachine(model.SwitchMachineState) error
	HasSwitchMachine(model.SwitchMachineId) bool
	GetSwitchMachineById(model.SwitchMachineId) model.SwitchMachineState
	GetAll() []model.SwitchMachineState
	RemoveSwitchMachine(model.SwitchMachineId) error
	UpdateSwitchMachine(model.SwitchMachineState) error
}

type switchMachineStoreImpl struct {
	switchMachines map[model.SwitchMachineId]model.SwitchMachineState
	rwLock         *sync.RWMutex
}

func NewSwitchMachineStore() SwitchMachineStore {
	smStore := &switchMachineStoreImpl{}
	smStore.rwLock = &sync.RWMutex{}
	smStore.switchMachines = make(map[model.SwitchMachineId]model.SwitchMachineState)
	return smStore
}

func (this *switchMachineStoreImpl) AddSwitchMachine(newSwitchMachine model.SwitchMachineState) error {
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

func (this *switchMachineStoreImpl) HasSwitchMachine(sMachineId model.SwitchMachineId) bool {
	this.rwLock.RLock()
	_, ok := this.switchMachines[sMachineId]
	this.rwLock.RUnlock()
	return ok
}

func (this *switchMachineStoreImpl) GetSwitchMachineById(sMachineId model.SwitchMachineId) model.SwitchMachineState {
	this.rwLock.RLock()
	sm, _ := this.switchMachines[sMachineId]
	this.rwLock.RUnlock()
	return sm
}

func (this *switchMachineStoreImpl) RemoveSwitchMachine(smId model.SwitchMachineId) error {
	var err error

	if !this.HasSwitchMachine(smId) {
		err = newDoesNotContainSwitchMacineWithIdError()
	} else {
		this.rwLock.Lock()
		delete(this.switchMachines, smId)
		this.rwLock.Unlock()
	}

	return err
}

func (this *switchMachineStoreImpl) UpdateSwitchMachine(sm model.SwitchMachineState) error {
	var err error
	if !this.HasSwitchMachine(sm.Id()) {
		err = newDoesNotContainSwitchMacineWithIdError()
	} else {
		this.rwLock.Lock()
		this.switchMachines[sm.Id()] = sm
		this.rwLock.Unlock()
	}
	return err
}

func (this *switchMachineStoreImpl) GetAll() []model.SwitchMachineState {
	this.rwLock.RLock()
	allSM := make([]model.SwitchMachineState, 0, len(this.switchMachines))
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
