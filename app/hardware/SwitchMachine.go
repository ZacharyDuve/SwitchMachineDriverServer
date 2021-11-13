package hardware

type SwitchMachineId uint16

type SwitchMachinePosition uint8

type SwitchMachineMotorState uint8

type GPIOState bool

const (
	MotorStateIdle   SwitchMachineMotorState = 0
	MotorStateToPos0 SwitchMachineMotorState = 1
	MotorStateToPos1 SwitchMachineMotorState = 2
	MotorStateBrake  SwitchMachineMotorState = 3

	Position0       SwitchMachinePosition = 0
	Position1       SwitchMachinePosition = 1
	PositionUnknown SwitchMachinePosition = 2

	GPIOOn  GPIOState = true
	GPIOOFF GPIOState = false
)

type SwitchMachineState interface {
	Id() SwitchMachineId
	Position() SwitchMachinePosition
	MotorState() SwitchMachineMotorState
	GPIO0State() GPIOState
	GPIO1State() GPIOState
}

type switchMachineStateImpl struct {
	id           SwitchMachineId
	pos          SwitchMachinePosition
	motor        SwitchMachineMotorState
	gpio0, gpio1 GPIOState
}

func NewSwitchMachineState(id SwitchMachineId, pos SwitchMachinePosition, m SwitchMachineMotorState, g0, g1 GPIOState) *switchMachineStateImpl {
	s := &switchMachineStateImpl{}
	s.id = id
	s.pos = pos
	s.motor = m
	s.gpio0 = g0
	s.gpio1 = g1

	return s
}

func (this *switchMachineStateImpl) Id() SwitchMachineId {
	return this.id
}

func (this *switchMachineStateImpl) Position() SwitchMachinePosition {
	return this.pos
}

func (this *switchMachineStateImpl) MotorState() SwitchMachineMotorState {
	return this.motor
}

func (this *switchMachineStateImpl) GPIO0State() GPIOState {
	return this.gpio0
}

func (this *switchMachineStateImpl) GPIO1State() GPIOState {
	return this.gpio1
}
