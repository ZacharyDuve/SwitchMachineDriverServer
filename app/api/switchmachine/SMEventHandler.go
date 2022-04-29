package switchmachine

import (
	"encoding/json"
	"log"

	"github.com/ZacharyDuve/SwitchMachineDriverServer/app/controller"
	"github.com/ZacharyDuve/SwitchMachineDriverServer/app/controller/model"
	"github.com/ZacharyDuve/eventsocket"
	"github.com/ZacharyDuve/serverid"
	"github.com/gorilla/mux"
)

const (
	eventHandlerSubPath string = "/event"
)

func RegsiterEventHandler(r *mux.Router, c controller.TortoiseController) {
	eventServer := eventsocket.NewEventServer()
	idSvc, err := serverid.NewFileServerIdService("")
	if err != nil {
		panic(err)
	}
	c.SetSwitchMachineAddedListenerFunc(func(s model.SwitchMachineState) {
		jsonSMEvent := mapSMStateToJSONSMEventData(s)

		data, err := json.Marshal(jsonSMEvent)
		if err != nil {
			log.Println(err)
			return
		}
		e, err := eventsocket.NewEvent(SMAdded, idSvc.GetServerId(), string(data))
		if err != nil {
			log.Println(err)
			return
		}
		eventServer.Send(e)
	})
	c.SetSwitchMachineRemovedListenerFunc(func(smi model.SwitchMachineId) {
		jsonSMEvent := &jsonSMEventData{Id: smi}

		data, err := json.Marshal(jsonSMEvent)
		if err != nil {
			log.Println(err)
			return
		}
		e, err := eventsocket.NewEvent(SMRemoved, idSvc.GetServerId(), string(data))
		if err != nil {
			log.Println(err)
			return
		}
		eventServer.Send(e)
	})
	c.SetSwitchMachineUpdatedListenerFunc(func(s model.SwitchMachineState) {
		jsonSMEvent := mapSMStateToJSONSMEventData(s)

		data, err := json.Marshal(jsonSMEvent)
		if err != nil {
			log.Println(err)
			return
		}
		e, err := eventsocket.NewEvent(SMUpdated, idSvc.GetServerId(), string(data))
		if err != nil {
			log.Println(err)
			return
		}
		eventServer.Send(e)
	})
	r.HandleFunc(eventHandlerSubPath, eventServer.ServeHTTP)
}

func mapSMStateToJSONSMEventData(s model.SwitchMachineState) *jsonSMEventData {
	data := &jsonSMEventData{}
	data.Id = s.Id()
	data.Position = s.Position()
	data.Motor = s.MotorState()
	data.GPIO0 = s.GPIO0State()
	data.GPIO1 = s.GPIO1State()
	return data
}
