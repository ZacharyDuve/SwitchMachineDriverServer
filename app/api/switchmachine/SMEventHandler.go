package switchmachine

import (
	"encoding/json"
	"log"

	"github.com/ZacharyDuve/SwitchMachineDriverServer/app/api/model"
	"github.com/ZacharyDuve/SwitchMachineDriverServer/app/controller"
	"github.com/ZacharyDuve/SwitchMachineDriverServer/app/controller/event"
	"github.com/ZacharyDuve/SwitchMachineDriverServer/app/controller/switchmachine"
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
	c.SetSwitchMachineEventListenerFunc(func(sme event.SwitchMachineEvent) {
		jsonSMEvent := mapSMStateToJSONSMEventData(sme.State())
		eventType := mapSMEventTypeToAPISMEventType(sme.Type())

		data, err := json.Marshal(jsonSMEvent)
		if err != nil {
			log.Println(err)
			return
		}
		e, err := eventsocket.NewEvent(eventType, idSvc.GetServerId(), string(data))
		if err != nil {
			log.Println(err)
			return
		}
		eventServer.Send(e)
	})
	r.HandleFunc(eventHandlerSubPath, eventServer.ServeHTTP)
}

func mapSMStateToJSONSMEventData(s switchmachine.State) *jsonSMEventData {
	data := &jsonSMEventData{}
	data.Id = model.SwitchMachineId(s.Id())
	data.Position = model.MapModelPosToApiPos(s.Position())
	data.Motor = model.MapModelMStateToAPIMState(s.MotorState())
	data.GPIO0 = model.MapModelGPIOToAPI(s.GPIO0State())
	data.GPIO1 = model.MapModelGPIOToAPI(s.GPIO1State())
	return data
}

func mapSMEventTypeToAPISMEventType(t event.EventType) string {
	switch t {
	case event.SwitchMachineAdded:
		return SMAdded
	case event.SwitchMachineUpdated:
		return SMUpdated
	case event.SwitchMachineRemoved:
		return SMRemoved
	default:
		panic("Unable to convert event.EventType to SMEvent")
	}
}
