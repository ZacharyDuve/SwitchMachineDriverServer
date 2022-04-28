package api

import (
	"net"
	"net/http"
	"strconv"

	"github.com/ZacharyDuve/SwitchMachineDriverServer/app/api/handler"
	"github.com/ZacharyDuve/SwitchMachineDriverServer/app/environment"
	"github.com/ZacharyDuve/apireg"
	"github.com/ZacharyDuve/apireg/api"
	"github.com/ZacharyDuve/serverid"
	"github.com/gorilla/mux"
)

const (
	ApiName          string = "SMDriverServer"
	ApiVersionMajor  uint   = 0
	ApiVersionMinor  uint   = 1
	ApiVersionBugFix uint   = 1
)

type smdsAPI struct {
	router      *mux.Router
	apiRegistry apireg.ApiRegistry
}

func NewSMDSApi() *smdsAPI {
	api := &smdsAPI{}
	reg, err := apireg.NewRegistry(environment.GetCurrent())
	if err != nil {
		panic(err)
	}
	api.apiRegistry = reg
	api.router = mux.NewRouter()
	sIdSvc, err := serverid.NewFileServerIdService("")
	if err != nil {
		panic(err)
	}
	api.router.HandleFunc(serverid.GetHandlerFuncFromServerIdService(sIdSvc))
	handler.NewSwitchMachineHandler(api.router)
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
	http.ListenAndServe(addr, this.router)
}
