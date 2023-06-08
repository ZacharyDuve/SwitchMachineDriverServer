package main

import "github.com/ZacharyDuve/SwitchMachineDriverServer/app/api"

func main() {
	api := api.NewSMDSApi()

	api.ListenAndServe(":8080")
}
