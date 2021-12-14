package model

import "github.com/ZacharyDuve/SwitchMachineDriverServer/app/controller/model"

type GpioState string

// List of GPIOState
const (
	OFF GpioState = "OFF"
	ON  GpioState = "ON"
)

func mapModelGPIOToAPI(modelGPIO model.GPIOState) GpioState {
	if modelGPIO == model.GPIOOFF {
		return OFF
	} else {
		return ON
	}
}
