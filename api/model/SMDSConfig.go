package model

type SMDSConfig struct {
	NumControllerBoards uint `json"numberControllerBoards"`
}

func (this *SMDSConfig) NumberControllerBoards() uint {
	return this.NumControllerBoards
}
