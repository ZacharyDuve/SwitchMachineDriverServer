package handler

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"

	"github.com/ZacharyDuve/SwitchMachineDriverServer/api/model"
	"github.com/ZacharyDuve/SwitchMachineDriverServer/app/smdsconfig"
	"github.com/gorilla/mux"
)

func AddConfigHandlerToRouter(router *mux.Router) {
	subRouter := router.PathPrefix("/config").Subrouter()
	subRouter.Methods(http.MethodGet).HandlerFunc(handleGetConfig)
	subRouter.Methods(http.MethodPut).HandlerFunc(handlePutConfig)
}

func handleGetConfig(w http.ResponseWriter, r *http.Request) {
	// bodyData,  := io.ReadAll(r.Body())
	config := smdsconfig.GetSMDSConfig()

	responseConfig := &model.SMDSConfig{NumControllerBoards: config.NumberControllerBoards()}

	responseBody, err := json.Marshal(responseConfig)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Print(err)
		w.Write([]byte(err.Error()))
	} else {
		w.WriteHeader(http.StatusOK)
		w.Write(responseBody)
	}
}

func handlePutConfig(w http.ResponseWriter, r *http.Request) {
	rqData, err := io.ReadAll(r.Body)

	if err == nil && len(rqData) != 0 {
		config := &model.SMDSConfig{}
		err = json.Unmarshal(rqData, config)
		if err == nil {
			err = smdsconfig.UpdateAndSaveSMDSConfig(config)
			if err == nil {
				retSMDSConfig := smdsconfig.GetSMDSConfig()
				config.NumControllerBoards = retSMDSConfig.NumberControllerBoards()
				var respData []byte
				respData, err = json.Marshal(config)

				if err == nil {
					w.WriteHeader(http.StatusOK)
					w.Write(respData)
				}
			}
		}
	} else if err == nil && len(rqData) == 0 {
		err = errors.New("Missing Request Body")
		w.Write([]byte(err.Error()))
		w.WriteHeader(http.StatusBadRequest)
	} else {
		w.Write([]byte(err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
	}

}
