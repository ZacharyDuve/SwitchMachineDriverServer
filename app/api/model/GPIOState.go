package model

import (
	"git.zmanhobbies.com/software/SwitchMachineDriverServer/app/controller/switchmachine"
)

type GPIOState string

// List of GPIOState
const (
	OFF GPIOState = "off"
	ON  GPIOState = "on"
)

func MapModelGPIOToAPI(modelGPIO switchmachine.GPIOState) GPIOState {
	if modelGPIO == switchmachine.GPIOOFF {
		return OFF
	} else {
		return ON
	}
}
