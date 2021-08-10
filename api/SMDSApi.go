package api

import (
	"net/http"

	"github.com/ZacharyDuve/SwitchMachineDriverServer/api/handler"
	"github.com/gorilla/mux"
)

type smdsAPI struct {
	router *mux.Router
}

func NewSMDSApi() *smdsAPI {
	api := &smdsAPI{}
	api.router = mux.NewRouter()
	handler.AddConfigHandlerToRouter(api.router)
	return api
}

func (this *smdsAPI) ListenAndServe(addr string) {
	http.ListenAndServe(addr, this.router)
}
