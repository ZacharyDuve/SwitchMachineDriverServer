package switchmachine

import (
	"log"
	"net/http"
	"sync"

	"github.com/ZacharyDuve/SwitchMachineDriverServer/app/api/model"
	"github.com/ZacharyDuve/SwitchMachineDriverServer/app/controller"
	"github.com/ZacharyDuve/SwitchMachineDriverServer/app/controller/event"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

const (
	eventHandlerSubPath string = "/event"
)

func RegsiterEventHandler(r *mux.Router, c controller.TortoiseController) {
	eventServer := newEventServer()

	c.SetSwitchMachineEventListenerFunc(func(sme event.SwitchMachineEvent) {
		smEvent := &model.SwitchMachineEvent{}
		smEvent.SwitchMachineState = model.NewAPISwitchMachineFromModel(sme.State())
		smEvent.EventType = model.MapSMEventToAPISMEventType(sme)
		eventServer.SendSwitchMachineEvent(smEvent)
	})
	r.HandleFunc(eventHandlerSubPath, eventServer.ServeHTTP)

}

type eventServer struct {
	upgrader     websocket.Upgrader
	clients      []*websocket.Conn
	clientsMutex *sync.Mutex
}

func newEventServer() *eventServer {
	eS := &eventServer{}
	eS.upgrader.CheckOrigin = func(r *http.Request) bool {
		return true
	}
	eS.clients = make([]*websocket.Conn, 0)
	eS.clientsMutex = &sync.Mutex{}
	return eS
}

func (this *eventServer) Close() error {
	this.clientsMutex.Lock()
	for _, curClient := range this.clients {
		curClient.Close()
	}
	this.clientsMutex.Unlock()
	return nil
}
func (this *eventServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c, err := this.upgrader.Upgrade(w, r, nil)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	this.clientsMutex.Lock()
	this.clients = append(this.clients, c)
	this.clientsMutex.Unlock()
}

func (this *eventServer) SendSwitchMachineEvent(sme *model.SwitchMachineEvent) {
	if sme != nil {
		var deadClients []*websocket.Conn
		this.clientsMutex.Lock()
		for _, curClient := range this.clients {
			err := curClient.WriteJSON(sme)
			if err != nil {
				log.Println(err)
				if websocket.IsUnexpectedCloseError(err) {
					if deadClients == nil {
						deadClients = make([]*websocket.Conn, 0)
					}
					deadClients = append(deadClients, curClient)
				}
			}
		}
		if len(deadClients) > 0 {
			this.removeDeadClients(deadClients)
		}
		this.clientsMutex.Unlock()
	}
}

func (this *eventServer) removeDeadClients(dClients []*websocket.Conn) {
	this.clientsMutex.Lock()
	newClientsSlice := make([]*websocket.Conn, 0, len(this.clients)-len(dClients))
	for _, curClient := range this.clients {
		isDead := false
		for _, curDeadClient := range dClients {
			if curDeadClient == curClient {
				isDead = true
				break
			}
		}
		if !isDead {
			newClientsSlice = append(newClientsSlice, curClient)
		}
	}

	this.clientsMutex.Unlock()
}
