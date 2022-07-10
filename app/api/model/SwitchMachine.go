package model

import (
	"time"

	"github.com/ZacharyDuve/SwitchMachineDriverServer/app/controller/switchmachine"
	"github.com/ZacharyDuve/serverid"
	"github.com/google/uuid"
)

var serveridSrv serverid.ServerIdService

func init() {
	var err error
	serveridSrv, err = serverid.NewFileServerIdService("")
	if err != nil {
		panic(err)
	}
}

type SwitchMachineId int

type SwitchMachine struct {
	SMId SwitchMachineId `json:"id"`

	Pos SwitchMachinePosition `json:"position"`

	Motor SwitchMachineMotorState `json:"motorState"`

	Gpio0 GPIOState `json:"gpio0"`

	Gpio1 GPIOState `json:"gpio1"`

	UpdTime time.Time `json:"updateTime"`

	OriginServerId uuid.UUID `json:"originServerId"`
}

func NewAPISwitchMachineFromModel(modelSM switchmachine.State) *SwitchMachine {
	apiSM := &SwitchMachine{}
	apiSM.SMId = SwitchMachineId(modelSM.Id())
	apiSM.Pos = MapModelPosToApiPos(modelSM.Position())
	apiSM.Gpio0 = MapModelGPIOToAPI(modelSM.GPIO0State())
	apiSM.Gpio1 = MapModelGPIOToAPI(modelSM.GPIO1State())
	apiSM.Motor = MapModelMStateToAPIMState(modelSM.MotorState())
	apiSM.UpdTime = modelSM.UpdateTime()
	apiSM.OriginServerId = serveridSrv.GetServerId()
	return apiSM
}

func (this *SwitchMachine) Id() switchmachine.Id {
	return switchmachine.Id(this.SMId)
}

func (this *SwitchMachine) Position() switchmachine.Position {
	if this.Pos == Position0 {
		return switchmachine.Position0
	} else {
		return switchmachine.Position1
	}
}

func (this *SwitchMachine) MotorState() switchmachine.MotorState {
	switch this.Motor {
	case BRAKE:
		return switchmachine.MotorStateBrake
	case TO_POSITION_0:
		return switchmachine.MotorStateToPos0
	case TO_POSITION_1:
		return switchmachine.MotorStateToPos1
	default:
		return switchmachine.MotorStateIdle
	}
}

func (this *SwitchMachine) GPIO0State() switchmachine.GPIOState {
	return mapAPIGPIOStateToHardwareState(this.Gpio0)
}

func (this *SwitchMachine) GPIO1State() switchmachine.GPIOState {
	return mapAPIGPIOStateToHardwareState(this.Gpio1)
}

func (this *SwitchMachine) UpdateTime() time.Time {
	return this.UpdTime
}

func mapAPIGPIOStateToHardwareState(apiState GPIOState) switchmachine.GPIOState {
	if apiState == ON {
		return switchmachine.GPIOOn
	} else {
		return switchmachine.GPIOOFF
	}
}
