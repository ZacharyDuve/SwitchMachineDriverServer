package tortoise

import (
	"log"
	"time"

	"github.com/ZacharyDuve/SwitchMachineDriverServer/app/controller/event"
	"github.com/ZacharyDuve/SwitchMachineDriverServer/app/controller/hardware"
	"periph.io/x/conn/v3/physic"
	"periph.io/x/conn/v3/spi"
	"periph.io/x/conn/v3/spi/spireg"
	"periph.io/x/host/v3"
)

const (
	spiClockSpeed     physic.Frequency = physic.KiloHertz * 10
	spiBusDevPath     string           = "/dev/spidev0"
	spiDevPath        string           = spiBusDevPath + ".0"
	spiMode           spi.Mode         = spi.Mode2
	spiBitsPerWord    int              = 8
	busUpdateDuration time.Duration    = time.Millisecond * 100
)

type piTortoiseControllerDriver struct {
	writeSPIPort    spi.PortCloser
	writeConn       spi.Conn
	readSPIPort     spi.PortCloser
	readConn        spi.Conn
	busUpdateTicker *time.Ticker
}

func init() {
	_, err := host.Init()

	if err != nil {
		panic(err)
	}
}

func NewPiTortoiseControllerDriver(smEventListener event.SwitchMachineEventListener) (piDriver hardware.SwitchMachineDriver, err error) {
	return NewPiTortoiseControllerDriverWithSPIDevPath(spiDevPath, smEventListener)
}

func NewPiTortoiseControllerDriverWithSPIDevPath(sDevPath string, smEventListener event.SwitchMachineEventListener) (piDriver hardware.SwitchMachineDriver, err error) {
	log.Println("NewPiTortoiseControllerDriverWithSPIDevPath called")
	ticker := time.NewTicker(busUpdateDuration)

	var driver *baseTortoiseControllerDriver
	trxFunc, closeFunc, err := setupConnections(sDevPath)
	log.Println("SPI connections setup. Errored:", err)
	if err == nil {
		piCloseFunc := func() (clsErr error) {
			ticker.Stop()
			return closeFunc()
		}

		driver, err = newBaseTortiseControllerDriver(trxFunc, piCloseFunc, ticker.C, smEventListener)
	}
	return driver, err
}

func setupConnections(spiDevPath string) (trxFunc func(w, r []byte) error, clsFunc func() error, err error) {
	log.Println("Setting up connection to:", spiDevPath)
	var initErr error

	var spiPort spi.PortCloser
	var spiConn spi.Conn

	//Open port and connections
	spiPort, initErr = spireg.Open(spiDevPath)
	if initErr == nil {
		spiConn, initErr = spiPort.Connect(spiClockSpeed, spiMode, spiBitsPerWord)

		if initErr == nil {
			clsFunc = spiPort.Close
			trxFunc = spiConn.Tx
		}
	}

	return trxFunc, clsFunc, initErr
}
