package tortoise

import (
	"fmt"
	"io"
	"log"
	"os"
	"sync"
	"time"

	"github.com/ZacharyDuve/SwitchMachineDriverServer/app/controller/hardware"
)

type MockHardwareDriver interface {
	hardware.Driver
	SetRXData([]byte)
	SetOutputForTx(io.Writer)
}

type mockHardwareDriverImpl struct {
	baseTortoiseControllerDriver
	mockRXData []byte
	rxMutex    *sync.Mutex
	txWriter   io.Writer
	txMutex    *sync.Mutex
}

func NewMockTortoiseControllerDriver() MockHardwareDriver {
	driver := createMockDriverImpl()

	ticker := time.NewTicker(time.Second * 2)
	driver.rxTrigger = ticker.C
	clsFunc := func() error {
		ticker.Stop()
		return nil
	}

	driver.closeFunc = clsFunc

	return driver
}

func NewMockTortoiseControllerDriverWithExternalRXTrigger(trig chan time.Time) MockHardwareDriver {
	driver := createMockDriverImpl()
	driver.rxTrigger = trig
	driver.closeFunc = func() error {
		return nil
	}
	return driver
}

func createMockDriverImpl() *mockHardwareDriverImpl {
	driver := &mockHardwareDriverImpl{}
	driver.rxMutex = &sync.Mutex{}
	driver.txMutex = &sync.Mutex{}
	driver.txWriter = os.Stdout
	driver.mockRXData = make([]byte, 8)

	txFunc := func(w, r []byte) error {
		driver.txMutex.Lock()
		fmt.Fprintf(driver.txWriter, "% X\n", w)
		driver.txMutex.Unlock()
		return nil
	}

	driver.txFunc = txFunc

	rxFunc := func(w, r []byte) error {

		driver.rxMutex.Lock()
		numCoppied := copy(r, driver.mockRXData)
		driver.rxMutex.Unlock()
		log.Println("mock rxFunc Called", numCoppied, r, driver.mockRXData)
		return nil
	}
	driver.rxFunc = rxFunc

	return driver
}

func printTRXFunc(w, r []byte) error {

	return nil
}

func (this *mockHardwareDriverImpl) SetRXData(mockRx []byte) {
	this.rxMutex.Lock()
	copy(this.mockRXData, mockRx)
	this.rxMutex.Unlock()
}

func (this *mockHardwareDriverImpl) SetOutputForTx(w io.Writer) {
	this.txMutex.Lock()
	this.txWriter = w
	this.txMutex.Unlock()
}
