package switchmachine

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"

	apiModel "github.com/ZacharyDuve/SwitchMachineDriverServer/app/api/model"
	"github.com/ZacharyDuve/SwitchMachineDriverServer/app/controller"
	"github.com/ZacharyDuve/SwitchMachineDriverServer/app/controller/hardware"
	"github.com/ZacharyDuve/SwitchMachineDriverServer/app/controller/hardware/tortoise"
	"github.com/ZacharyDuve/SwitchMachineDriverServer/app/controller/switchmachine"
	"github.com/ZacharyDuve/SwitchMachineDriverServer/app/environment"
	env "github.com/ZacharyDuve/apireg/environment"
	"github.com/gorilla/mux"
)

type switchMachineHandler struct {
	controller controller.TortoiseController
}

const (
	idRequestKey  string = "id"
	smHandlerPath string = "/switchmachine"
)

func NewSwitchMachineHandler(rtr *mux.Router) {
	smHandler := &switchMachineHandler{}
	subRtr := rtr.PathPrefix(smHandlerPath).Subrouter()
	var driver hardware.Driver
	if environment.GetCurrent() == env.Prod {
		var err error
		driver, err = tortoise.NewPiTortoiseControllerDriver()
		if err != nil {
			panic(err)
		}
	} else {
		trigger := make(chan time.Time)
		mockDriver := tortoise.NewMockTortoiseControllerDriverWithExternalRXTrigger(trigger)
		driver = mockDriver
		//txOut := make([]byte, 0)
		subRtr.PathPrefix("/mockrxdata").Methods(http.MethodPost).HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
			rxData, err := ioutil.ReadAll(hex.NewDecoder(r.Body))

			if err != nil {
				rw.WriteHeader(http.StatusBadRequest)
				rw.Write([]byte(err.Error()))
			} else {
				mockDriver.SetRXData(rxData)
				log.Println("Sent RX data", rxData)
				trigger <- time.Now()
			}
		})
	}
	smHandler.controller = controller.NewTortoiseController(driver)
	RegsiterEventHandler(subRtr, smHandler.controller)
	subRtr.PathPrefix("/{" + idRequestKey + "}").Methods(http.MethodGet).HandlerFunc(smHandler.handleGetSwitchMachine)
	//For updating a switch machine we are just going to put to the base
	subRtr.Methods(http.MethodPut).HandlerFunc(smHandler.handleUpdateSwitchMachine)

	subRtr.Methods(http.MethodGet).HandlerFunc(smHandler.handleGetSwitchMachines)
}

func (this *switchMachineHandler) handleGetSwitchMachines(w http.ResponseWriter, r *http.Request) {
	switchMachines := this.controller.GetSwitchMachines()
	//TODO optimize the creation of the slice
	apiSMs := make([]apiModel.SwitchMachine, 0)
	for _, curSM := range switchMachines {
		apiSMs = append(apiSMs, *apiModel.NewAPISwitchMachineFromModel(curSM))
	}

	encodeErr := json.NewEncoder(w).Encode(apiSMs)

	if encodeErr != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (this *switchMachineHandler) handleGetSwitchMachine(w http.ResponseWriter, r *http.Request) {
	smId, err := getSMIdFromRequest(r)
	sm, err := this.controller.GetSwitchMachineById(smId)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	encodeErr := json.NewEncoder(w).Encode(apiModel.NewAPISwitchMachineFromModel(sm))

	if encodeErr != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (this *switchMachineHandler) handleUpdateSwitchMachine(w http.ResponseWriter, r *http.Request) {

	switchMachineReq := &apiModel.SwitchMachine{}
	err := json.NewDecoder(r.Body).Decode(switchMachineReq)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	err = this.controller.UpdateSwitchMachine(switchMachineReq)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	}
}

func getSMIdFromRequest(r *http.Request) (switchmachine.Id, error) {
	var err error
	var smId switchmachine.Id
	routeVars := mux.Vars(r)
	smIdStr, hasSMId := routeVars[idRequestKey]
	if !hasSMId {
		err = errors.New("Request is missing id")
	} else {
		var smIdInt int64
		smIdInt, err = strconv.ParseInt(smIdStr, 10, 0)

		if err != nil {
			err = errors.New("Malformed id in request")
		} else {
			smId = switchmachine.Id(smIdInt)
		}
	}
	return smId, err
}
