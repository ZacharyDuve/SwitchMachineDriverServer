package api

import (
	"log"
	"net"
	"net/http"
	"strconv"

	"github.com/ZacharyDuve/SwitchMachineDriverServer/app/api/switchmachine"
	"github.com/ZacharyDuve/SwitchMachineDriverServer/app/environment"
	"github.com/ZacharyDuve/apireg"
	"github.com/ZacharyDuve/apireg/api"
	"github.com/ZacharyDuve/serverid"
	"github.com/gorilla/mux"
)

const (
	ApiName          string = "SMDriverServer"
	ApiVersionMajor  uint   = 0
	ApiVersionMinor  uint   = 10
	ApiVersionBugFix uint   = 0
)

type smdsAPI struct {
	router *mux.Router
	//apiSubRouter *mux.Router
	apiRegistry apireg.ApiRegistry
}

func NewSMDSApi() *smdsAPI {
	log.Println("Start Creating NewSMDSApi")
	api := &smdsAPI{}
	reg, err := apireg.NewRegistry(environment.GetCurrent())
	if err != nil {
		panic(err)
	}
	api.apiRegistry = reg
	api.router = mux.NewRouter()
	//Need an API sub router to separate from web
	apiSubRouter := api.router.PathPrefix("/api").Subrouter()
	//DO WE WANT ONLY FOR NON PRODS?
	apiSubRouter.Use(mux.CORSMethodMiddleware(apiSubRouter))
	sIdSvc, err := serverid.NewFileServerIdService("")
	if err != nil {
		panic(err)
	}
	//Make it so that we can get the server id
	apiSubRouter.HandleFunc(serverid.GetHandlerFuncFromServerIdService(sIdSvc))
	//Register the switch machine handler with the api sub router
	switchmachine.NewSwitchMachineHandler(apiSubRouter)
	//Need to serve any non api routes as web pages
	api.router.PathPrefix("/").Handler(http.FileServer(http.Dir("web-content")))
	log.Println("End Creating NewSMDSApi")
	return api
}

func (this *smdsAPI) ListenAndServe(addr string) {
	_, portString, err := net.SplitHostPort(addr)
	if err != nil {
		panic(err)
	}
	port, err := strconv.Atoi(portString)
	if err != nil {
		panic(err)
	}
	this.apiRegistry.RegisterApi(ApiName, &api.Version{Major: ApiVersionMajor, Minor: ApiVersionMinor, BugFix: ApiVersionBugFix}, port)
	log.Println("Server up and waiting for requests")
	log.Fatal(http.ListenAndServe(addr, this.router))
}
