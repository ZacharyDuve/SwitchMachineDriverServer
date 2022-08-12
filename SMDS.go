package main

import "git.zmanhobbies.com/software/SwitchMachineDriverServer/app/api"

func main() {
	api := api.NewSMDSApi()

	api.ListenAndServe(":8080")
}
