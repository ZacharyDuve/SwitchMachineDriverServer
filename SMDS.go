package main

import (
	"log"
	"log/slog"
	"net/http"
	"os"

	"github.com/ZacharyDuve/SwitchMachineDriverServer/app/api"
	"github.com/ZacharyDuve/serverid"
)

func main() {

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))

	// Need to load in the server id in for the server id service.
	sIDSvc, err := serverid.NewFileServerIdService("")
	if err != nil {
		panic(err)
	}

	api, err := api.NewSMDSApi(logger, sIDSvc)

	if err != nil {
		log.Fatal("error starting api", "error", err)
	}

	http.ListenAndServe(":8080", api)
}
