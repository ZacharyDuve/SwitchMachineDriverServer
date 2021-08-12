package model

type Error struct {
	Code string `json:"code,omitempty"`

	Message string `json:"message,omitempty"`

	Causedby *Error `json:"causedby,omitempty"`
}
