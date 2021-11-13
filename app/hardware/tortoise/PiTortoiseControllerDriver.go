package tortoise

import (
	"github.com/ZacharyDuve/SwitchMachineDriverServer/app/hardware"
	"periph.io/x/conn/v3/physic"
	"periph.io/x/conn/v3/spi"
	"periph.io/x/conn/v3/spi/spireg"
	"periph.io/x/host/v3"
)

const (
	spiClockSpeed   physic.Frequency = physic.KiloHertz * 10
	spiBusDevPath   string           = "/dev/spidev0"
	spiWriteDevPath string           = spiBusDevPath + ".0"
	spiReadDevPath  string           = spiBusDevPath + ".1"
	spiWriteMode    spi.Mode         = spi.Mode2
	spiReadMode     spi.Mode         = spi.Mode2
	spiBitsPerWord  int              = 8
)

type piTortoiseControllerDriver struct {
	writeSPIPort spi.PortCloser
	writeConn    spi.Conn
	readSPIPort  spi.PortCloser
	readConn     spi.Conn
}

func init() {
	_, err := host.Init()

	if err != nil {
		panic(err)
	}
}

func NewPiTortoiseControllerDriver() (piDriver hardware.SwitchMachineDriver, err error) {
	return NewPiTortoiseControllerDriverWithSPIDevPath(spiBusDevPath)
}

func NewPiTortoiseControllerDriverWithSPIDevPath(spiDevPath string) (piDriver hardware.SwitchMachineDriver, err error) {

	var driver *baseTortoiseControllerDriver
	trxFunc, closeFunc, err := setupConnections(spiDevPath)

	if err == nil {
		driver, err = newBaseTortiseControllerDriver(trxFunc, closeFunc)
	}
	return driver, err
}

func setupConnections(spiDevPath string) (trxFunc func(w, r []byte) error, clsFunc func() error, err error) {
	var initErr error

	var spiPort spi.PortCloser
	var spiConn spi.Conn

	//Open port and connections
	spiPort, initErr = spireg.Open(spiDevPath)
	if initErr == nil {
		spiConn, initErr = spiPort.Connect(spiClockSpeed, spiWriteMode, spiBitsPerWord)

		if initErr == nil {
			clsFunc = func() error {
				return spiPort.Close()
			}
			trxFunc = spiConn.Tx
		}
	}

	return nil, clsFunc, initErr
}

//PI SPI is MSB

// type turnout struct {
// 	id uint
// 	currentPos TurnoutPosition
// 	currentDriveState driveState
// 	gpio0 bool
// 	gpio1 bool
// }

// type turnout struct {
// 	id         TurnoutID
// 	motorState      byte
// 	currentPos TurnoutPosition
// 	gpio0      GPIOState
// 	gpio1      GPIOState
// 	//Byte that could be written out spi, saved so we don't have to recalculate it everytime
// 	dataByte byte
// }

// func (this *turnout) GetDataByte() byte {

// }
