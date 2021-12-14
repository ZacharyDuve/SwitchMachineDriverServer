package handler

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"

	apiModel "github.com/ZacharyDuve/SwitchMachineDriverServer/app/api/model"
	"github.com/ZacharyDuve/SwitchMachineDriverServer/app/controller"
	"github.com/ZacharyDuve/SwitchMachineDriverServer/app/controller/model"
	"github.com/gorilla/mux"
)

type switchMachineHandler struct {
	controller controller.TortoiseController
}

const (
	idRequestKey string = "id"
)

func NewSwitchMachineHandler(rtr *mux.Router) {
	smHandler := &switchMachineHandler{}
	subRtr := rtr.PathPrefix("/switchmachine").Subrouter()
	subRtr.PathPrefix("/{" + idRequestKey + "}/position").Methods(http.MethodPut).HandlerFunc(smHandler.handleUpdateSwitchMachinePosition)
	subRtr.PathPrefix("/{" + idRequestKey + "}/gpio").Methods(http.MethodPut).HandlerFunc(smHandler.handleUpdateSwitchMachineGPIO)
	subRtr.PathPrefix("/{" + idRequestKey + "}").Methods(http.MethodGet).HandlerFunc(smHandler.handleGetSwitchMachine)

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

func (this *switchMachineHandler) handleUpdateSwitchMachinePosition(w http.ResponseWriter, r *http.Request) {

	smId, err := getSMIdFromRequest(r)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	posUpdateReq := &apiModel.SwitchMachinePositionUpdateRequest{}
	err = json.NewDecoder(r.Body).Decode(posUpdateReq)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	err = this.controller.SetSwitchMachinePosition(smId, posUpdateReq.GetPosition())

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	}
}

func (this *switchMachineHandler) handleUpdateSwitchMachineGPIO(w http.ResponseWriter, r *http.Request) {
	smId, err := getSMIdFromRequest(r)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	gpioUpdateReq := &apiModel.SwitchMachineGPIOUpdateRequest{}

	err = json.NewDecoder(r.Body).Decode(gpioUpdateReq)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	err = this.controller.SetSwitchMachineGPIO(smId, gpioUpdateReq.GPIO0(), gpioUpdateReq.GPIO1())

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	}
}

func getSMIdFromRequest(r *http.Request) (model.SwitchMachineId, error) {
	var err error
	var smId model.SwitchMachineId
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
			smId = model.SwitchMachineId(smIdInt)
		}
	}
	return smId, err
}
