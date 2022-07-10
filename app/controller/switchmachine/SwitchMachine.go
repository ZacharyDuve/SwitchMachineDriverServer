package switchmachine

import (
	"fmt"
	"log"
	"time"
)

type Id uint16

type Position uint8

type MotorState uint8

type GPIOState bool

const (
	MotorStateIdle   MotorState = 0
	MotorStateToPos0 MotorState = 1
	MotorStateToPos1 MotorState = 2
	MotorStateBrake  MotorState = 3

	Position0       Position = 0
	Position1       Position = 1
	PositionUnknown Position = 2

	GPIOOn  GPIOState = true
	GPIOOFF GPIOState = false
)

type State interface {
	Id() Id
	Position() Position
	MotorState() MotorState
	GPIO0State() GPIOState
	GPIO1State() GPIOState
	UpdateTime() time.Time
}

type switchMachineStateImpl struct {
	id           Id
	pos          Position
	motor        MotorState
	gpio0, gpio1 GPIOState
	updateTime   time.Time
}

func NewState(id Id, pos Position, m MotorState, g0, g1 GPIOState) State {
	s := &switchMachineStateImpl{}
	s.id = id
	s.pos = pos
	s.motor = m
	s.gpio0 = g0
	s.gpio1 = g1
	s.updateTime = time.Now()
	return s
}

func (this *switchMachineStateImpl) Id() Id {
	return this.id
}

func (this *switchMachineStateImpl) Position() Position {
	return this.pos
}

func (this *switchMachineStateImpl) MotorState() MotorState {
	return this.motor
}

func (this *switchMachineStateImpl) GPIO0State() GPIOState {
	return this.gpio0
}

func (this *switchMachineStateImpl) GPIO1State() GPIOState {
	return this.gpio1
}

func (this *switchMachineStateImpl) UpdateTime() time.Time {
	return this.updateTime
}

func StatesEqual(sm1, sm2 State) bool {
	return sm1.Id() == sm2.Id() &&
		sm1.GPIO0State() == sm2.GPIO0State() &&
		sm1.GPIO1State() == sm2.GPIO1State() &&
		sm1.MotorState() == sm2.MotorState() &&
		sm1.Position() == sm2.Position()
}

//----------------------------------- Printing functions for convience
func StateToString(state State) string {
	log.Println("state", state)
	return fmt.Sprintf("Id: %d, Position: %s, Motor: %s, GPIO0: %s, GPIO1: %s", state.Id(),
		positionToString(state.Position()), motorToString(state.MotorState()), gpioToString(state.GPIO0State()), gpioToString(state.GPIO1State()))
}

func gpioToString(g GPIOState) string {
	switch g {
	case GPIOOFF:
		return "OFF"
	default:
		return "ON"
	}
}

func positionToString(p Position) string {
	switch p {
	case Position0:
		return "0"
	case Position1:
		return "1"
	default:
		return "Unknown"
	}
}

func motorToString(m MotorState) string {
	switch m {
	case MotorStateBrake:
		return "Brake"
	case MotorStateToPos0:
		return "To 0"
	case MotorStateToPos1:
		return "To 1"
	default:
		return "Idle"
	}
}
