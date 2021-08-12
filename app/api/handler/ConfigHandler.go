package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/ZacharyDuve/SwitchMachineDriverServer/app/api/model"
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

	// responseBody, err := json.Marshal(responseConfig)
	err := json.NewEncoder(w).Encode(responseConfig)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Print(err)
		w.Write([]byte(err.Error()))
	}
}

func handlePutConfig(w http.ResponseWriter, r *http.Request) {
	config := &model.SMDSConfig{}
	err := json.NewDecoder(r.Body).Decode(config)

	if err == nil {
		err = smdsconfig.UpdateAndSaveSMDSConfig(config)
		if err == nil {
			err = json.NewEncoder(w).Encode(config)

		}
	} else {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
	}

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	}
}
