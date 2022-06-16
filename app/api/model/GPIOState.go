package model

import (
	"github.com/ZacharyDuve/SwitchMachineDriverServer/app/controller/switchmachine"
)

type GPIOState string

// List of GPIOState
const (
	OFF GPIOState = "OFF"
	ON  GPIOState = "ON"
)

func MapModelGPIOToAPI(modelGPIO switchmachine.GPIOState) GPIOState {
	if modelGPIO == switchmachine.GPIOOFF {
		return OFF
	} else {
		return ON
	}
}
