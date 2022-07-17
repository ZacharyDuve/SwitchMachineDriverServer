package tortoise

import (
	"log"
	"time"

	"github.com/ZacharyDuve/SwitchMachineDriverServer/app/controller/hardware"
	"periph.io/x/conn/v3/physic"
	"periph.io/x/conn/v3/spi"
	"periph.io/x/conn/v3/spi/spireg"
	"periph.io/x/host/v3"
)

const (
	spiClockSpeed     physic.Frequency = physic.KiloHertz * 1
	spiBusDevPath     string           = "/dev/spidev0"
	spiTxDevPath      string           = spiBusDevPath + ".1"
	spiRxDevPath      string           = spiBusDevPath + ".0"
	spiTxMode         spi.Mode         = spi.Mode2
	spiRxMode         spi.Mode         = spi.Mode0
	spiBitsPerWord    int              = 8
	busUpdateDuration time.Duration    = time.Millisecond * 100
)

type piTortoiseControllerDriver struct {
	writeSPIPort    spi.PortCloser
	writeConn       spi.Conn
	readSPIPort     spi.PortCloser
	readConn        spi.Conn
	busUpdateTicker *time.Ticker
	//baseTortoiseControllerDriver
}

func init() {
	_, err := host.Init()

	if err != nil {
		panic(err)
	}
}

func NewPiTortoiseControllerDriver() (piDriver hardware.Driver, err error) {
	return NewPiTortoiseControllerDriverWithSPIDevPath(spiTxDevPath, spiRxDevPath)
}

func NewPiTortoiseControllerDriverWithSPIDevPath(txDevPath, rxDevPath string) (hardware.Driver, error) {
	driver := &baseTortoiseControllerDriver{}
	log.Println("NewPiTortoiseControllerDriverWithSPIDevPath called")
	ticker := time.NewTicker(busUpdateDuration)
	driver.rxTrigger = ticker.C

	txFunc, txCloseFunc, txOpenErr := setupConnection(txDevPath, spiTxMode)
	log.Println("TX SPI connections setup.")
	if txOpenErr != nil {
		log.Println("Error opening tx line", txOpenErr)
		return nil, txOpenErr
	}

	driver.txFunc = txFunc

	rxFunc, rxCloseFunc, rxOpenErr := setupConnection(rxDevPath, spiRxMode)
	log.Println("RX SPI connections setup.")
	if rxOpenErr != nil {
		log.Println("Error opening rx line", rxOpenErr)
		return nil, rxOpenErr
	}

	driver.rxFunc = rxFunc

	piCloseFunc := func() (clsErr error) {
		ticker.Stop()
		defer txCloseFunc()
		defer rxCloseFunc()
		return nil
	}

	driver.closeFunc = piCloseFunc

	return driver, nil
}

func setupConnection(spiDevPath string, m spi.Mode) (xFunc func(w, r []byte) error, clsFunc func() error, err error) {
	log.Println("Setting up connection to:", spiDevPath)
	var initErr error

	var spiPort spi.PortCloser
	var spiConn spi.Conn

	//Open port and connections
	spiPort, initErr = spireg.Open(spiDevPath)
	if initErr == nil {
		spiConn, initErr = spiPort.Connect(spiClockSpeed, m, spiBitsPerWord)
		if initErr == nil {
			clsFunc = spiPort.Close
			xFunc = spiConn.Tx
		}
	}

	return xFunc, clsFunc, initErr
}
