package model

type SwitchMachineUpdateRequest struct {
	Id string `json:"id,omitempty"`

	Position *SwitchMachinePosition `json:"position,omitempty"`

	Gpio0 *GpioState `json:"gpio0,omitempty"`

	Gpio1 *GpioState `json:"gpio1,omitempty"`
}
