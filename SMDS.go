package main

import (
	"log"
	"log/slog"
	"net/http"
	"os"

	"github.com/ZacharyDuve/SwitchMachineDriverServer/app/api"
	"github.com/ZacharyDuve/serverid"
	"github.com/gorilla/mux"
)

func main() {

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))

	// Need to load in the server id in for the server id service.
	sIDSvc, err := serverid.NewFileServerIdService("")
	if err != nil {
		panic(err)
	}

	rootRouter := mux.NewRouter()

	err = api.NewSMDSApi(logger, rootRouter, sIDSvc)

	if err != nil {
		log.Fatal("error starting api", "error", err)
	}

	http.ListenAndServe(":8080", rootRouter)
}
