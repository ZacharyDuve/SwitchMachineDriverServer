package handler

import (
	"encoding/json"
	"net/http"

	"github.com/ZacharyDuve/SwitchMachineDriverServer/app/api/model"
	"github.com/ZacharyDuve/SwitchMachineDriverServer/app/controller"
	"github.com/gorilla/mux"
)

type switchMachineHandler struct {
	controller controller.TortoiseController
}

func NewSwitchMachineHandler(rtr *mux.Router) {
	smHandler := &switchMachineHandler{}
	subRtr := rtr.PathPrefix("/switchmachine").Subrouter()

	//subRtr.Methods(http.MethodGet).HandlerFunc(smHandler.handleGetSwitchMachines)
	subRtr.Methods(http.MethodPut).HandlerFunc(smHandler.handleUpdateSwitchMachine)
	smHandler.controller = controller.NewTortoiseController()
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

	this.controller.UpdateSwitchMachine(sm)
}
