package model

type Error struct {
	Message string `json:"message,omitempty"`

	Causedby *Error `json:"causedby,omitempty"`
}
