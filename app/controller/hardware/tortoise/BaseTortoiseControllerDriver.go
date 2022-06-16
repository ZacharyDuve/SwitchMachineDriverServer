package tortoise

import (
	"fmt"
	"log"
	"time"

	"github.com/ZacharyDuve/SwitchMachineDriverServer/app/controller/hardware"
	"github.com/ZacharyDuve/SwitchMachineDriverServer/app/controller/switchmachine"
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
	MaxNumberAttachableMainControllerBoards uint = 8
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

	motorStateBitMask byte = 0x0C
	motorIdleBits     byte = 0x00
	motorToPos0Bits   byte = 0x04
	motorToPos1Bits   byte = 0x08
	motorBrakeBits    byte = 0x0C

	gpioBitMask  byte = 0x03
	gpio0HighBit byte = 0x01
	gpio1HighBit byte = 0x02

	dataBitMask byte = gpioBitMask | motorStateBitMask

	positionBitMask byte = 0x03
	positionUnknown byte = 0x03
	//A and D are opposite of B and C
	position0Port12      byte = 0x01
	position1Port12      byte = 0x02
	position0Port03      byte = 0x02
	position1Port03      byte = 0x01
	positionDisconnected byte = 0x00
)

type baseTortoiseControllerDriver struct {
	driverEventListener hardware.DriverEventListener
	txBuffer            []byte
	txWasteRxBuffer     []byte
	prevRxBuffer        []byte
	rxBuffer            []byte
	rxWasteTxBuffer     []byte
	//Function that is attached that handles closing any connections in the implementing driver
	closeFunc func() error
	//Function that handles writing data to device while also reading data from it. Can return error if something goes wrong
	txFunc func(w, r []byte) error
	rxFunc func(w, r []byte) error
	//Channel to alert processLoop to exit
	processLoopExitChan chan bool
	//Channel to take in new SwitchMachine States to be processed.
	newSMStateChan chan switchmachine.State
	//Channel that triggers bus updates when value appears
	rxTrigger <-chan time.Time
}

func (this *baseTortoiseControllerDriver) UpdateSwitchMachine(newState switchmachine.State) {
	this.newSMStateChan <- newState
}

func (this *baseTortoiseControllerDriver) Start(driverEventListener hardware.DriverEventListener) {
	this.driverEventListener = driverEventListener
	this.initChans()
	this.initBuffers()
	go this.runLoop()
}

func (this *baseTortoiseControllerDriver) Close() error {
	this.processLoopExitChan <- false
	return this.closeFunc()
}

func (this *baseTortoiseControllerDriver) initChans() {
	this.processLoopExitChan = make(chan bool)
	this.newSMStateChan = make(chan switchmachine.State)
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

			curSMId := switchmachine.Id(portNumber + byteIndex*int(numRxPortsPerByte))
			wasAttached := isConnectedFromPositionBits(prevRxBits)
			isAttached := isConnectedFromPositionBits(curRxBits)

			log.Println("curSMId:", curSMId, "wasAttached:", wasAttached, "isAttached", isAttached)
			//We know we are removed
			if wasAttached && !isAttached {
				this.driverEventListener.HandleDriverEvent(hardware.NewSwitchMachineRemovedEvent(curSMId))
			} else {
				position := getSMPositionFromRxBits(curRxBits, portNumber)

				state := switchmachine.NewState(curSMId, position, switchmachine.MotorStateIdle, switchmachine.GPIOOFF, switchmachine.GPIOOFF)

				if !wasAttached {
					this.driverEventListener.HandleDriverEvent(hardware.NewSwitchMachineAddedEvent(curSMId, state))
				} else {
					this.driverEventListener.HandleDriverEvent(hardware.NewSwitchMachinePositionChangedEvent(curSMId, state))
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

func getSMPositionFromRxBits(rxBits byte, portNumber int) switchmachine.Position {
	var position switchmachine.Position
	if rxBits == positionUnknown {
		position = switchmachine.PositionUnknown
	}
	if portNumber == 0 || portNumber == 3 {
		if rxBits == position0Port03 {
			position = switchmachine.Position0
		} else if rxBits == position1Port03 {
			position = switchmachine.Position1
		}
	} else {
		if rxBits == position0Port12 {
			position = switchmachine.Position0
		} else if rxBits == position1Port12 {
			position = switchmachine.Position1
		}
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

func (this *baseTortoiseControllerDriver) processSMStateUpdate(newState switchmachine.State) {
	log.Println("Getting update", switchmachine.StateToString(newState))
	var txBits byte

	if newState.GPIO0State() {
		txBits = gpio0HighBit
	}
	if newState.GPIO1State() {
		txBits |= gpio1HighBit
	}

	switch newState.MotorState() {
	case switchmachine.MotorStateIdle:
		txBits |= motorIdleBits
	case switchmachine.MotorStateToPos0:
		txBits |= motorToPos0Bits
	case switchmachine.MotorStateToPos1:
		txBits |= motorToPos1Bits
	case switchmachine.MotorStateBrake:
		txBits |= motorBrakeBits
	}
	bitMask := dataBitMask
	//If it is an even id number than we need to left shift 4 bits
	if newState.Id()%2 == 0 {
		txBits = txBits << 4
		bitMask = bitMask << 4
	}

	byteIndex := getTxIndexFromBufferLengthAndId(len(this.txBuffer), newState.Id())

	this.txBuffer[byteIndex] = (this.txBuffer[byteIndex] & ^bitMask) | txBits
	this.handleBusWrite()
	log.Println("this.txBuffer", this.txBuffer, "byteIndex", byteIndex, "bitMask", bitMask, "txBits", txBits)
}
func getTxIndexFromBufferLengthAndId(bLen int, id switchmachine.Id) uint {
	return uint(bLen-1) - calcTxByteOffsetFromId(id)
}

func calcTxByteOffsetFromId(id switchmachine.Id) uint {
	//This tells us which board we are on
	boardNumber := uint(id) / numDriverPortsPerBoard
	//Now we need to figure out which of the bytes for the board we land on
	byteIndex := boardNumber*numTxBytesPerBoard + (uint(id)%numDriverPortsPerBoard)/numTxPortsPerByte

	return byteIndex
}

type TurnoutNotAvailableError struct {
	id switchmachine.Id
}

func (this *TurnoutNotAvailableError) Error() string {
	return fmt.Sprintf("Turnout with id : %d, is not available to be set. Probably you are setting a turnout for a driver board not attached.", this.id)
}

type TurnoutRequestNilError struct {
}

func (this *TurnoutRequestNilError) Error() string {
	return "TurnoutRequest was nil."
}
