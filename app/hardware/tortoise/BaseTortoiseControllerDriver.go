package tortoise

import (
	"errors"
	"fmt"
	"time"

	"github.com/ZacharyDuve/SwitchMachineDriverServer/app/hardware"
)

type bitOrder bool

const (
	//Writing one Tx Byte can drive 2 Tortoise Driver boards
	numTxPortsPerByte uint = 2
	//Reading one byte can read 4 Tortoise Driver Boards
	numRxPortsPerByte uint = 4
	//Number of ports to connect to individual driver boards from the main board
	numDriverPortsPerBoard uint = 4
	//Number of write bytes per main driver board
	numTxBytesPerBoard uint = numDriverPortsPerBoard / numTxPortsPerByte
	//Number of read bytes per main driver board
	numRxBytesPerBoard uint = numDriverPortsPerBoard / numRxPortsPerByte
	//MaxNumberAttachableMainControllerBoards is the limit of boards that one driver can control from one computer. This number is arbitrailily decided
	MaxNumberAttachableMainControllerBoards uint = 16
	//DefaultThrowTime is the default time that tortoise driver board will active the motor for to throw a turnout
	DefaultThrowTime time.Duration = time.Second * 2

	//msbFirst
	msbFirst bitOrder = false
	lsbFirst bitOrder = true
)

const (
	motorStateBitMask byte = 0x0C
	motorIdleBits     byte = 0x00
	motorToPos0Bits   byte = 0x04
	motorToPos1Bits   byte = 0x08
	motorBrakeBits    byte = 0x0C

	gpioBitMask  byte = 0x03
	gpio0HighBit byte = 0x01
	gpio1HighBit byte = 0x02

	dataBitMask byte = gpioBitMask | motorStateBitMask

	positionBitMask      byte = 0x03
	positionUnknown      byte = 0x03
	position0            byte = 0x01
	position1            byte = 0x02
	positionDisconnected byte = 0x00
)

type baseTortoiseControllerDriver struct {
	attachedSwitchMachines map[hardware.SwitchMachineId]*tortoiseSwitchMachine
	txBuffer               []byte
	prevRxBuffer           []byte
	rxBuffer               []byte
	//Function that is attached that handles closing any connections in the implementing driver
	closeFunc func() error
	//Function that handles writing data to device while also reading data from it. Can return error if something goes wrong
	txRxFunc func(w, r []byte) error
	//Channel to alert processLoop to exit
	processLoopExitChan chan bool
	//Channel to take in new SwitchMachine States to be processed.
	newSMStateChan chan hardware.SwitchMachineState
	//Ticker that sends tick to update bus if needed - basically so not every event updates the bus and we don't constantly read from it
	busUpdateTicker *time.Ticker
}

func newBaseTortiseControllerDriver(trxFunc func(w, r []byte) error, clsFunc func() error) (driver *baseTortoiseControllerDriver, err error) {
	if trxFunc == nil {
		err = errors.New("trxFunc is a required parameter for baseTortiseControllerDriver")
	} else if clsFunc == nil {
		err = errors.New("clsFunc is a required parameter for baseTortiseControllerDriver")
	} else {
		driver = &baseTortoiseControllerDriver{}
		driver.initBuffers()
		driver.initChans()
		driver.attachedSwitchMachines = make(map[hardware.SwitchMachineId]*tortoiseSwitchMachine)
		driver.busUpdateTicker = time.NewTicker(time.Millisecond * 200)

		driver.txRxFunc = trxFunc
		driver.closeFunc = clsFunc

		go driver.runLoop()
	}

	return driver, err
}

func (this *baseTortoiseControllerDriver) GetNumberSwitchMachinesConnected() uint {
	return 0
}

func (this *baseTortoiseControllerDriver) UpdateSwitchMachine(newState hardware.SwitchMachineState) error {
	var err error
	//Need to validate that there is even a switch machine there to perform update on
	//if _, exists := this.attachedSwitchMachines[newState.Id()]; exists {
	this.newSMStateChan <- newState
	//} else {
	//	err = &TurnoutNotAvailableError{id: newState.Id()}
	//}
	return err
}

func (this *baseTortoiseControllerDriver) Close() error {
	this.busUpdateTicker.Stop()
	return this.closeFunc()
}

func (this *baseTortoiseControllerDriver) initChans() {
	this.processLoopExitChan = make(chan bool)
	this.newSMStateChan = make(chan hardware.SwitchMachineState)
}

func (this *baseTortoiseControllerDriver) initBuffers() {
	this.txBuffer = make([]byte, MaxNumberAttachableMainControllerBoards*numTxBytesPerBoard)

	this.rxBuffer = make([]byte, MaxNumberAttachableMainControllerBoards*numRxBytesPerBoard)
	this.prevRxBuffer = make([]byte, len(this.rxBuffer))
}

func (this *baseTortoiseControllerDriver) runLoop() {
	for {
		select {
		case _ = <-this.processLoopExitChan:
			return
		case _ = <-this.busUpdateTicker.C:
			this.handleBusUpdate()
		case newSMState := <-this.newSMStateChan:
			this.processSMStateUpdate(newSMState)
		}
	}
}

func (this *baseTortoiseControllerDriver) handleBusUpdate() {
	this.txRxFunc(this.txBuffer, this.rxBuffer)
	//Figure out what changed
	//this.processRxBufferChanges()

	//We need to swap rx buffers so cur becomes prev and we can reuse old prev for next read since it completely overwrites
	//this.swapRxBuffers()
}

func (this *baseTortoiseControllerDriver) processRxBufferChanges() {
	for curIndex, curRxByte := range this.rxBuffer {
		prevRxByte := this.prevRxBuffer[curIndex]
		//If the bytes are not equal then something changed
		if prevRxByte != curRxByte {
			for bitIndex := 0; bitIndex < int(numRxPortsPerByte); bitIndex++ {
				//Mask off the bits so we only get one port worth of data
				prevRxBits := prevRxByte & dataBitMask
				curRxBits := curRxByte & dataBitMask
				//If they are different then lets handle the change
				if prevRxBits != curRxBits {

					portNumber := calcPortNumFromByteIndexAndBitIndex(curIndex, bitIndex)
					//See if we changed from disconnected to connected
					if prevRxBits == positionDisconnected {
						attachedSM := &tortoiseSwitchMachine{}
						attachedSM.id = hardware.SwitchMachineId(portNumber)
						attachedSM.gpio0 = hardware.GPIOOFF
						attachedSM.gpio1 = hardware.GPIOOFF
						attachedSM.motorState = hardware.MotorStateIdle
						this.attachedSwitchMachines[hardware.SwitchMachineId(portNumber)] = attachedSM
					}

					curSM := this.attachedSwitchMachines[hardware.SwitchMachineId(portNumber)]

					if curRxBits == positionUnknown {
						curSM.position = hardware.PositionUnknown
					} else if curRxBits == position0 {
						curSM.position = hardware.Position0
					} else if curRxBits == position1 {
						curSM.position = hardware.Position1
					} else {
						//We have become disconnected
						delete(this.attachedSwitchMachines, hardware.SwitchMachineId(portNumber))
					}
				}
			}
		}
	}
}

func calcPortNumFromByteIndexAndBitIndex(byteIndex, bitIndex int) int {
	var portNum int

	if bitIndex >= int(numRxPortsPerByte) {
		panic(errors.New("Unable to calc port as bit index is impossible value"))
	}

	if bitIndex == 0 {
		portNum = 2
	} else if bitIndex == 1 {
		portNum = 1
	} else if bitIndex == 2 {
		portNum = 0
	} else if bitIndex == 3 {
		portNum = 3
	}

	return portNum + byteIndex*int(numRxPortsPerByte)
}

func (this *baseTortoiseControllerDriver) isUpdateInRxBuffers() bool {
	var updateOccured bool
	for curIndex, curByte := range this.rxBuffer {
		if curByte != this.prevRxBuffer[curIndex] {
			updateOccured = true
			break
		}
	}

	return updateOccured
}

func (this *baseTortoiseControllerDriver) swapRxBuffers() {
	oldPrevRxBuff := this.prevRxBuffer
	this.prevRxBuffer = this.rxBuffer
	this.rxBuffer = oldPrevRxBuff
}

func (this *baseTortoiseControllerDriver) processSMStateUpdate(newState hardware.SwitchMachineState) {
	var txBits byte

	if newState.GPIO0State() {
		txBits = gpio0HighBit
	}
	if newState.GPIO1State() {
		txBits |= gpio1HighBit
	}

	switch newState.MotorState() {
	case hardware.MotorStateIdle:
		txBits |= motorIdleBits
	case hardware.MotorStateToPos0:
		txBits |= motorToPos0Bits
	case hardware.MotorStateToPos1:
		txBits |= motorToPos1Bits
	case hardware.MotorStateBrake:
		txBits |= motorBrakeBits
	}
	bitMask := dataBitMask
	//If it is an even id number than we need to left shift 4 bits
	if newState.Id()%2 == 0 {
		txBits = txBits << 4
		bitMask = bitMask << 4
	}

	byteIndex := newState.Id() / hardware.SwitchMachineId(numTxPortsPerByte)

	this.txBuffer[byteIndex] = (this.txBuffer[byteIndex] & bitMask) | txBits
}

type MaxMainDriverBoardLimitExceededError struct {
}

func (this *MaxMainDriverBoardLimitExceededError) Error() string {
	return fmt.Sprintf("Number of Main Tortoise Driver Boards attached to this driver has been exceeded. The max is %d.", MaxNumberAttachableMainControllerBoards)
}

type TurnoutNotAvailableError struct {
	id hardware.SwitchMachineId
}

func (this *TurnoutNotAvailableError) Error() string {
	return fmt.Sprintf("Turnout with id : %d, is not available to be set. Probably you are setting a turnout for a driver board not attached.", this.id)
}

type TurnoutRequestNilError struct {
}

func (this *TurnoutRequestNilError) Error() string {
	return "TurnoutRequest was nil."
}
