package tortoise

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/ZacharyDuve/SwitchMachineDriverServer/app/controller/event"
	"github.com/ZacharyDuve/SwitchMachineDriverServer/app/controller/model"
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
	MaxNumberAttachableMainControllerBoards uint = 1
	//DefaultThrowTime is the default time that tortoise driver board will active the motor for to throw a turnout
	DefaultThrowTime time.Duration = time.Second * 2

	numBitsPerPort uint = 2

	port0RxBitMask   byte = 0x30
	port0RxBitIndex  uint = 2
	port0RxBitOffset uint = port0RxBitIndex * numBitsPerPort

	port1RxBitMask   byte = 0x0C
	port1RxBitIndex  uint = 1
	port1RxBitOffset uint = port1RxBitIndex * numBitsPerPort

	port2RxBitMask   byte = 0x03
	port2RxBitIndex  uint = 0
	port2RxBitOffset uint = port2RxBitIndex * numBitsPerPort

	port3RxBitMask   byte = 0xC0
	port3RxBitIndex  uint = 3
	port3RxBitOffset uint = port3RxBitIndex * numBitsPerPort
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
	smEventListener event.SwitchMachineEventListener
	txBuffer        []byte
	txWasteRxBuffer []byte
	prevRxBuffer    []byte
	rxBuffer        []byte
	rxWasteTxBuffer []byte
	//Function that is attached that handles closing any connections in the implementing driver
	closeFunc func() error
	//Function that handles writing data to device while also reading data from it. Can return error if something goes wrong
	txFunc func(w, r []byte) error
	rxFunc func(w, r []byte) error
	//Channel to alert processLoop to exit
	processLoopExitChan chan bool
	//Channel to take in new SwitchMachine States to be processed.
	newSMStateChan chan model.SwitchMachineState
	//Channel that triggers bus updates when value appears
	rxTrigger <-chan time.Time
}

func newBaseTortiseControllerDriver(txFunc, rxFunc func(w, r []byte) error, clsFunc func() error, rxTrigger <-chan time.Time, smEventListner event.SwitchMachineEventListener) (driver *baseTortoiseControllerDriver, err error) {
	if txFunc == nil {
		err = errors.New("txFunc is a required parameter for baseTortiseControllerDriver")
	} else if rxFunc == nil {
		err = errors.New("rxFunc is a required parameter for baseTortiseControllerDriver")
	} else if clsFunc == nil {
		err = errors.New("clsFunc is a required parameter for baseTortiseControllerDriver")
	} else if rxTrigger == nil {
		err = errors.New("rxTrigger is a required parameter for baseTortiseControllerDriver")
	} else {
		driver = &baseTortoiseControllerDriver{}
		driver.initBuffers()
		driver.initChans()
		driver.smEventListener = smEventListner

		driver.txFunc = txFunc
		driver.rxFunc = rxFunc
		driver.closeFunc = clsFunc
		driver.rxTrigger = rxTrigger

		go driver.runLoop()
	}

	return driver, err
}

func (this *baseTortoiseControllerDriver) UpdateSwitchMachine(newState model.SwitchMachineState) {
	this.newSMStateChan <- newState
}

func (this *baseTortoiseControllerDriver) Close() error {
	this.processLoopExitChan <- false
	return this.closeFunc()
}

func (this *baseTortoiseControllerDriver) initChans() {
	this.processLoopExitChan = make(chan bool)
	this.newSMStateChan = make(chan model.SwitchMachineState)
}

func (this *baseTortoiseControllerDriver) initBuffers() {
	this.txBuffer = make([]byte, MaxNumberAttachableMainControllerBoards*numTxBytesPerBoard)
	this.txWasteRxBuffer = make([]byte, len(this.txBuffer))

	this.rxBuffer = make([]byte, MaxNumberAttachableMainControllerBoards*numRxBytesPerBoard)
	this.prevRxBuffer = make([]byte, len(this.rxBuffer))
	this.rxWasteTxBuffer = make([]byte, len(this.rxBuffer))
}

func (this *baseTortoiseControllerDriver) runLoop() {
	for {
		select {
		case _ = <-this.processLoopExitChan:
			return
		case _ = <-this.rxTrigger:
			this.handleBusRead()
		case newSMState := <-this.newSMStateChan:
			this.processSMStateUpdate(newSMState)
		}
	}
}

func (this *baseTortoiseControllerDriver) handleBusWrite() {
	this.txFunc(this.txBuffer, this.txWasteRxBuffer)
}

func (this *baseTortoiseControllerDriver) handleBusRead() {
	log.Println("Handling the bus read")
	this.rxFunc(this.rxWasteTxBuffer, this.rxBuffer)
	//Figure out what changed
	this.processRxBufferChanges()

	//We need to swap rx buffers so cur becomes prev and we can reuse old prev for next read since it completely overwrites
	this.swapRxBuffers()
}

func (this *baseTortoiseControllerDriver) processRxBufferChanges() {
	for curIndex, curRxByte := range this.rxBuffer {
		prevRxByte := this.prevRxBuffer[curIndex]
		//If the bytes are not equal then something changed
		if prevRxByte != curRxByte {
			this.handleRXByteChange(prevRxByte, curRxByte, curIndex)
		}
	}
}

func (this *baseTortoiseControllerDriver) handleRXByteChange(prevRxByte, curRxByte byte, byteIndex int) {
	for portNumber := 0; portNumber < int(numRxPortsPerByte); portNumber++ {
		//Mask off the bits so we only get one port worth of data
		prevRxBits := getRxBitsForPortNumber(prevRxByte, portNumber)
		curRxBits := getRxBitsForPortNumber(curRxByte, portNumber)
		//If they are different then lets handle the change
		if prevRxBits != curRxBits {
			log.Println("Bytes don't match need to figure out what changed", prevRxBits, " ", curRxBits)

			curSMId := model.SwitchMachineId(portNumber)
			wasAttached := isConnectedFromPositionBits(prevRxBits)
			isAttached := isConnectedFromPositionBits(curRxBits)

			log.Println("curSMId:", curSMId, "wasAttached:", wasAttached, "isAttached", isAttached)
			//We know we are removed
			if wasAttached && !isAttached {
				this.smEventListener.SwitchMachineRemoved(curSMId)
			} else {
				position := getSMPositionFromRxBits(curRxBits)

				state := model.NewSwitchMachineState(curSMId, position, model.MotorStateIdle, model.GPIOOFF, model.GPIOOFF)

				if !wasAttached {
					this.smEventListener.SwitchMachineAdded(state)
				} else {
					this.smEventListener.SwitchMachineUpdated(state)
				}
			}

		}
	}
}

func getRxBitsForPortNumber(rxByte byte, portNum int) byte {
	var rxBits byte
	switch portNum {
	case 0:
		rxBits = (rxByte & port0RxBitMask) >> byte(port0RxBitOffset)
	case 1:
		rxBits = (rxByte & port1RxBitMask) >> byte(port1RxBitOffset)
	case 2:
		rxBits = (rxByte & port2RxBitMask) >> byte(port2RxBitOffset)
	case 3:
		rxBits = (rxByte & port3RxBitMask) >> byte(port3RxBitOffset)
	default:
		panic("Invalid Port number")
	}

	return rxBits
}

func getSMPositionFromRxBits(rxBits byte) model.SwitchMachinePosition {
	var position model.SwitchMachinePosition
	if rxBits == positionUnknown {
		position = model.PositionUnknown
	} else if rxBits == position0 {
		position = model.Position0
	} else if rxBits == position1 {
		position = model.Position1
	} else {
		panic("Switch Machine is disconnected cant tell position")
	}
	return position
}

func isConnectedFromPositionBits(posBits byte) bool {
	return posBits != positionDisconnected
}

func (this *baseTortoiseControllerDriver) swapRxBuffers() {
	oldPrevRxBuff := this.prevRxBuffer
	this.prevRxBuffer = this.rxBuffer
	this.rxBuffer = oldPrevRxBuff
}

func (this *baseTortoiseControllerDriver) processSMStateUpdate(newState model.SwitchMachineState) {
	log.Println("Getting update", newState)
	var txBits byte

	if newState.GPIO0State() {
		txBits = gpio0HighBit
	}
	if newState.GPIO1State() {
		txBits |= gpio1HighBit
	}

	switch newState.MotorState() {
	case model.MotorStateIdle:
		txBits |= motorIdleBits
	case model.MotorStateToPos0:
		txBits |= motorToPos0Bits
	case model.MotorStateToPos1:
		txBits |= motorToPos1Bits
	case model.MotorStateBrake:
		txBits |= motorBrakeBits
	}
	bitMask := dataBitMask
	//If it is an even id number than we need to left shift 4 bits
	if newState.Id()%2 == 0 {
		txBits = txBits << 4
		bitMask = bitMask << 4
	}

	//byteIndex := uint(len(this.txBuffer)) - (uint(newState.Id()) / numTxPortsPerByte)
	byteIndex := uint(newState.Id()) / numTxPortsPerByte

	this.txBuffer[byteIndex] = (this.txBuffer[byteIndex] & ^bitMask) | txBits
	this.handleBusWrite()
	log.Println("this.txBuffer", this.txBuffer, "byteIndex", byteIndex, "bitMask", bitMask, "txBits", txBits)
}

type MaxMainDriverBoardLimitExceededError struct {
}

func (this *MaxMainDriverBoardLimitExceededError) Error() string {
	return fmt.Sprintf("Number of Main Tortoise Driver Boards attached to this driver has been exceeded. The max is %d.", MaxNumberAttachableMainControllerBoards)
}

type TurnoutNotAvailableError struct {
	id model.SwitchMachineId
}

func (this *TurnoutNotAvailableError) Error() string {
	return fmt.Sprintf("Turnout with id : %d, is not available to be set. Probably you are setting a turnout for a driver board not attached.", this.id)
}

type TurnoutRequestNilError struct {
}

func (this *TurnoutRequestNilError) Error() string {
	return "TurnoutRequest was nil."
}
