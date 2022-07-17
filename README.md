# SwitchMachineDriverServer
Server to drive switch machine controller circuit boards that I have developed via REST API

To build for raspberry pi 1 use command GOOS=linux GOARCH=arm GOARM=5 go build -o "builds/SMDS-$(git describe --tags --abbrev=0)-linux-arm-5" SMDS.go
