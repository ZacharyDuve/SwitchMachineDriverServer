package model

import "github.com/ZacharyDuve/SwitchMachineDriverServer/app/controller/model"

type GPIOState string

// List of GPIOState
const (
	OFF GPIOState = "OFF"
	ON  GPIOState = "ON"
)

func MapModelGPIOToAPI(modelGPIO model.GPIOState) GPIOState {
	if modelGPIO == model.GPIOOFF {
		return OFF
	} else {
		return ON
	}
}
