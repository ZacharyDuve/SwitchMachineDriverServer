package handler

import (
	"encoding/json"
	"net/http"

	"github.com/ZacharyDuve/SwitchMachineDriverServer/app/api/model"
	"github.com/ZacharyDuve/SwitchMachineDriverServer/app/hardware"
	"github.com/ZacharyDuve/SwitchMachineDriverServer/app/hardware/tortoise"
	"github.com/gorilla/mux"
)

type switchMachineHandler struct {
	driver hardware.SwitchMachineDriver
}

func NewSwitchMachineHandler(rtr *mux.Router) {
	smHandler := &switchMachineHandler{}
	subRtr := rtr.PathPrefix("/switchmachine").Subrouter()

	//subRtr.Methods(http.MethodGet).HandlerFunc(smHandler.handleGetSwitchMachines)
	subRtr.Methods(http.MethodPut).HandlerFunc(smHandler.handleUpdateSwitchMachine)
	var err error
	smHandler.driver, err = tortoise.NewPiTortoiseControllerDriver()

	if err != nil {
		panic(err)
	}
}

func (this *switchMachineHandler) handleGetSwitchMachines(w http.ResponseWriter, r *http.Request) {

}

func (this *switchMachineHandler) handleUpdateSwitchMachine(w http.ResponseWriter, r *http.Request) {
	sm := &model.SwitchMachine{}
	err := json.NewDecoder(r.Body).Decode(sm)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	this.driver.UpdateSwitchMachine(sm)
}
