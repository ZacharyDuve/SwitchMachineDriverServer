package handler

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

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
	if os.Getenv("environment") == "production" {
		smHandler.controller = controller.NewTortoiseController()
	} else {
		rxIn := make([]byte, 5, 5)
		txOut := make([]byte, 0)
		subRtr.PathPrefix("/mockrxdata").Methods(http.MethodPost).HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
			rxData, err := ioutil.ReadAll(hex.NewDecoder(r.Body))

			if err != nil {
				rw.WriteHeader(http.StatusBadRequest)
				rw.Write([]byte(err.Error()))
			} else {
				copy(rxIn, rxData)
				fmt.Println(rxIn)
			}
		})
		smHandler.controller = controller.NewTortoiseControllerWithMockDriver(txOut, rxIn)
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

	err = this.controller.UpdateSwitchMachine(sm)

	if err != nil {
		if controller.IsSwitchMachineNotExistError(err) {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(fmt.Sprint("Switch machine ", sm.SMId, " does not exist")))
		}
	}
}
