package switchmachine

import (
	"time"

	"github.com/ZacharyDuve/SwitchMachineDriverServer/app/controller/model"
	"github.com/ZacharyDuve/eventsocket"
	"github.com/google/uuid"
)

const (
	SMAdded   string = "SwitchMachineAdded"
	SMRemoved string = "SwitchMachineRemoved"
	SMUpdated string = "SwitchMachineUpdated"
)

type SMEvent interface {
	Name() string
	OriginTime() time.Time
	OriginId() uuid.UUID
	SwitchMachineId() model.SwitchMachineId
	SwitchMachinePosition() model.SwitchMachinePosition
	SwitchMachineMotorState() model.SwitchMachineMotorState
	SwitchMachineGPIO0() model.GPIOState
	SwitchMachineGPIO1() model.GPIOState
}

type smEvent struct {
	e     eventsocket.Event
	id    model.SwitchMachineId
	pos   model.SwitchMachinePosition
	motor model.SwitchMachineMotorState
	gpio0 model.GPIOState
	gpio1 model.GPIOState
}

func (this *smEvent) Name() string {
	return this.e.Name()
}
func (this *smEvent) OriginTime() time.Time {
	return this.e.OriginTime()
}
func (this *smEvent) OriginId() uuid.UUID {
	return this.e.OriginID()
}

func (this *smEvent) SwitchMachineId() model.SwitchMachineId {
	return this.id
}
func (this *smEvent) SwitchMachinePosition() model.SwitchMachinePosition {
	return this.pos
}

func (this *smEvent) SwitchMachineMotorState() model.SwitchMachineMotorState {
	return this.motor
}

func (this *smEvent) SwitchMachineGPIO0() model.GPIOState {
	return this.gpio0
}

func (this *smEvent) SwitchMachineGPIO1() model.GPIOState {
	return this.gpio1
}

type jsonSMEventData struct {
	Id       model.SwitchMachineId         `json:"id"`
	Position model.SwitchMachinePosition   `json:"position"`
	Motor    model.SwitchMachineMotorState `json:"motor"`
	GPIO0    model.GPIOState               `json:"gpio0"`
	GPIO1    model.GPIOState               `json:"gpio1"`
}
