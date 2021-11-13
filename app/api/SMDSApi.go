package api

import (
	"net/http"

	"github.com/ZacharyDuve/SwitchMachineDriverServer/app/api/handler"
	"github.com/gorilla/mux"
)

type smdsAPI struct {
	router *mux.Router
}

func NewSMDSApi() *smdsAPI {
	api := &smdsAPI{}
	api.router = mux.NewRouter()
	handler.AddConfigHandlerToRouter(api.router)
	handler.NewSwitchMachineHandler(api.router)
	return api
}

func (this *smdsAPI) ListenAndServe(addr string) {
	http.ListenAndServe(addr, this.router)
}
