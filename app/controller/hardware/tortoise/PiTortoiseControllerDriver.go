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
	spiClockSpeed     physic.Frequency = physic.MegaHertz * 1
	spiBusDevPath     string           = "/dev/spidev0"
	spiTxDevPath      string           = spiBusDevPath + ".0"
	spiRxDevPath      string           = spiBusDevPath + ".1"
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
}

func init() {
	_, err := host.Init()

	if err != nil {
		panic(err)
	}
}

func NewPiTortoiseControllerDriver(smEventListener event.SwitchMachineEventListener) (piDriver hardware.SwitchMachineDriver, err error) {
	return NewPiTortoiseControllerDriverWithSPIDevPath(spiTxDevPath, spiRxDevPath, smEventListener)
}

func NewPiTortoiseControllerDriverWithSPIDevPath(txDevPath, rxDevPath string, smEventListener event.SwitchMachineEventListener) (hardware.SwitchMachineDriver, error) {
	log.Println("NewPiTortoiseControllerDriverWithSPIDevPath called")
	ticker := time.NewTicker(busUpdateDuration)

	txFunc, txCloseFunc, txOpenErr := setupConnection(txDevPath, spiTxMode)
	log.Println("TX SPI connections setup. Errored:", txOpenErr)
	if txOpenErr != nil {
		log.Println("Error opening tx line", txOpenErr)
		return nil, txOpenErr
	}
	rxFunc, rxCloseFunc, rxOpenErr := setupConnection(rxDevPath, spiRxMode)
	log.Println("RX SPI connections setup. Errored:", rxOpenErr)
	if rxOpenErr != nil {
		log.Println("Error opening rx line", rxOpenErr)
		return nil, rxOpenErr
	}

	piCloseFunc := func() (clsErr error) {
		ticker.Stop()
		defer txCloseFunc()
		defer rxCloseFunc()
		return nil
	}

	return newBaseTortiseControllerDriver(txFunc, rxFunc, piCloseFunc, ticker.C, smEventListener)
}

func setupConnection(spiDevPath string, m spi.Mode) (xFunc func(w, r []byte) error, clsFunc func() error, err error) {
	log.Println("Setting up connection to:", spiDevPath)
	var initErr error

	var spiPort spi.PortCloser
	var spiConn spi.Conn

	//Open port and connections
	spiPort, initErr = spireg.Open(spiDevPath)
	if initErr == nil {
		log.Println("Going to start opening connection")
		spiConn, initErr = spiPort.Connect(spiClockSpeed, m, spiBitsPerWord)
		log.Println("Should have opened spi port")
		if initErr == nil {
			log.Println("Binding close and x functions")
			clsFunc = spiPort.Close
			xFunc = spiConn.Tx
		}
	}

	return xFunc, clsFunc, initErr
}
