package switchmachine

import (
	"bytes"
	"encoding/json"
	"log"

	"github.com/ZacharyDuve/SwitchMachineDriverServer/app/api/model"
	"github.com/ZacharyDuve/SwitchMachineDriverServer/app/controller"
	"github.com/ZacharyDuve/SwitchMachineDriverServer/app/controller/event"
	"github.com/ZacharyDuve/eventsocket"
	"github.com/gorilla/mux"
)

const (
	eventHandlerSubPath string = "/event"
)

func RegsiterEventHandler(r *mux.Router, c controller.TortoiseController) {
	eventServer := eventsocket.NewEventServer()

	c.SetSwitchMachineEventListenerFunc(func(sme event.SwitchMachineEvent) {
		smEvent := &model.SwitchMachineEvent{}
		smEvent.SwitchMachineState = model.NewAPISwitchMachineFromModel(sme.State())
		smEvent.EventType = model.MapSMEventToAPISMEventType(sme)

		data, err := json.Marshal(smEvent)
		if err != nil {
			log.Println(err)
			return
		}

		eventServer.Send(bytes.NewBuffer(data))
	})
	r.HandleFunc(eventHandlerSubPath, eventServer.ServeHTTP)

}
