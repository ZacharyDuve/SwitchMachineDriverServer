package api

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/ZacharyDuve/SwitchMachineDriverServer/app/api/switchmachine"

	"github.com/ZacharyDuve/serverid"
	"github.com/gorilla/mux"
)

const (
	ApiName          string = "SMDriverServer"
	ApiVersionMajor  uint   = 0
	ApiVersionMinor  uint   = 10
	ApiVersionBugFix uint   = 0
)

// type smdsAPI struct {
// 	logger *slog.Logger
// 	router *mux.Router
// 	//apiSubRouter *mux.Router
// 	//apiRegistry apireg.ApiRegistry
// }

func NewSMDSApi(logger *slog.Logger, svrIDSvc serverid.ServerIdService) (http.Handler, error) {

	if logger == nil {
		return nil, errors.New("error call to NewSMDSApi failed due to missing logger")
	}

	if svrIDSvc == nil {
		return nil, errors.New("error call to NewSMDSApi failed due to missing ServerIDService")
	}

	logger.Debug("Start Creating NewSMDSApi")

	// TODO: Look to see what other service discovery tools exist as this was my own
	//reg, err := apireg.NewRegistry(environment.GetCurrent())

	// if err != nil {
	// 	panic(err)
	// }
	//api.apiRegistry = reg

	router := mux.NewRouter()
	//Need an API sub router to separate from web
	apiSubRouter := router.PathPrefix("/api").Subrouter()

	//Make it so that we can get the server id
	apiSubRouter.HandleFunc(serverid.GetHandlerFuncFromServerIdService(svrIDSvc))
	//Register the switch machine handler with the api sub router
	switchmachine.NewSwitchMachineHandler(apiSubRouter)

	return router, nil
}
