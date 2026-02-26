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

func NewSMDSApi(logger *slog.Logger, router *mux.Router, svrIDSvc serverid.ServerIdService) error {

	if logger == nil {
		return errors.New("error call to NewSMDSApi failed due to missing logger")
	}

	if svrIDSvc == nil {
		return errors.New("error call to NewSMDSApi failed due to missing ServerIDService")
	}

	logger.Debug("Start Creating NewSMDSApi")

	// TODO: Look to see what other service discovery tools exist as this was my own
	//reg, err := apireg.NewRegistry(environment.GetCurrent())

	// if err != nil {
	// 	panic(err)
	// }
	//api.apiRegistry = reg

	router.Use(loggingMiddleware(logger))

	//Need an API sub router to separate from web
	// router.PathPrefix("/api").Handler(switchmachine.NewSwitchMachineHandler(logger, nil))
	router.Path("/health").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("{\"Status\":\"On fire\"}"))
	})
	apiRouter := router.PathPrefix("/api").Subrouter()

	switchmachine.NewSwitchMachineHandler(logger, apiRouter, nil)
	//Make it so that we can get the server id
	//apiSubRouter.HandleFunc(serverid.GetHandlerFuncFromServerIdService(svrIDSvc))
	//Register the switch machine handler with the api sub router
	return nil
}

func loggingMiddleware(logger *slog.Logger) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			logger.Debug("http request", "method", r.Method, "host", r.Host, "url", r.URL.Path)
			next.ServeHTTP(w, r)
		})
	}
}
