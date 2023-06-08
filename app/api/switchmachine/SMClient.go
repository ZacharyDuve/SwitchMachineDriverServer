package switchmachine

// import (
// 	"encoding/json"
// 	"errors"
// 	"fmt"
// 	"io"
// 	"io/ioutil"
// 	"log"
// 	"net"
// 	"net/http"
// 	"sync"

// 	"github.com/ZacharyDuve/SwitchMachineDriverServer/app/api/model"
// 	"github.com/ZacharyDuve/apireg"
// 	"github.com/ZacharyDuve/apireg/api"
// 	"github.com/ZacharyDuve/apireg/environment"
// 	"github.com/ZacharyDuve/apireg/event"
// 	"github.com/ZacharyDuve/eventsocket"
// )

// //Client written in go for current SwitchMachineHandler
// type SwitchMachineClient interface {
// 	UpdateSwitchMachines([]*model.SwitchMachine) error
// 	io.Closer
// }

// type smClient struct {
// 	aReg           apireg.ApiRegistry
// 	eventFunc      func(SMEvent)
// 	serverIdToApis *sync.Map
// 	apisToServerId *sync.Map
// 	eventClients   *sync.Map
// }

// func NewClient(a net.Addr, eventFunc func(SMEvent)) (SwitchMachineClient, error) {
// 	var err error
// 	var c *smClient
// 	if eventFunc == nil {
// 		err = errors.New("eventFunc is required for a NewClient")
// 	} else {
// 		c = &smClient{}
// 		c.serverIdToApis = &sync.Map{}
// 		c.apisToServerId = &sync.Map{}
// 		c.eventClients = &sync.Map{}
// 		c.eventFunc = eventFunc
// 		//TODO this probably shouldn't be ALL
// 		c.aReg, err = apireg.NewRegistry(environment.All)

// 		c.aReg.AddListener(c)
// 	}

// 	return c, err
// }

// func (this *smClient) UpdateSwitchMachines(sMachines []*model.SwitchMachine) error {

// 	// for _, curSM := range sMachines {
// 	// 	api, this.serverIdToApis.Load()
// 	// }
// 	return nil
// }

// func (this *smClient) Close() error {
// 	var err error

// 	this.eventClients.Range(func(key, value any) bool {
// 		switch client := value.(type) {
// 		case *smClient:
// 			client.Close()
// 		default:
// 			if err == nil {
// 				err = errors.New("Error trying to close client but type is not *smClient")
// 			}
// 			log.Println("Error trying to close client")
// 		}
// 		return true
// 	})

// 	return err
// }

// func (this *smClient) HandleRegistration(e event.RegistrationEvent) {
// 	if isSMDSRegisterEvent(e) {
// 		this.handleNewSMDS(e)
// 	} else if isSMDSUnRegisterEvent(e) {
// 		this.handleRemoveSMDS(e)
// 	}
// }

// func isSMDSRegisterEvent(e event.RegistrationEvent) bool {
// 	return e.Type() == event.Added && e.Api().Name() == "SMDriverServer"
// }

// func isSMDSUnRegisterEvent(e event.RegistrationEvent) bool {
// 	return e.Type() == event.Removed && e.Api().Name() == "SMDriverServer"
// }

// func (this *smClient) handleNewSMDS(e event.RegistrationEvent) {
// 	sId, err := getServerId(e)

// 	if err == nil {
// 		this.serverIdToApis.Store(sId, e.Api())
// 		this.apisToServerId.Store(apiHostStringFromApi(e.Api()), sId)
// 		var smEventClient io.Closer
// 		smEventClient, err = NewSMEventClient(&net.TCPAddr{IP: e.Api().HostIP(), Port: e.Api().HostPort()}, this.handleEvent)
// 		if err == nil {
// 			this.eventClients.Store(sId, smEventClient)
// 			smEvents := getCurrentSMFromSMDS(e.Api())
// 			for _, curE := range smEvents {
// 				this.eventFunc(curE)
// 			}
// 		}
// 	}
// }

// func getCurrentSMFromSMDS(a api.Api) []SMEvent {
// 	var events []SMEvent

// 	r, err := http.Get(apiHostStringFromApi(a) + "/switchmachine")

// 	if err == nil {
// 		var switchMachines []*model.SwitchMachine
// 		err = json.NewDecoder(r.Body).Decode(&switchMachines)
// 		if err != nil {
// 			log.Println("Error decoding switch machines", err.Error())
// 		} else {
// 			for _, curSM := range switchMachines {
// 				event, _ := eventsocket.NewEvent(SMAdded, curSM.OriginServerId, "")
// 				events = append(events, &smEvent{id: curSM.SMId, pos: curSM.Pos, motor: curSM.Motor, gpio0: curSM.Gpio0, gpio1: curSM.Gpio1, e: event})
// 			}
// 		}
// 	}

// 	return events
// }

// func (this *smClient) handleRemoveSMDS(e event.RegistrationEvent) {
// 	hostString := apiHostStringFromApi(e.Api())
// 	//Get the server id
// 	sId, ok := this.apisToServerId.Load(hostString)

// 	if ok {
// 		c, hasClient := this.eventClients.Load(sId)
// 		if hasClient {
// 			switch client := c.(type) {
// 			case *smClient:
// 				client.Close()
// 			default:
// 				panic(errors.New("Error trying to close client but type is not *smClient"))
// 			}
// 			//Delete this client after closing it
// 			this.eventClients.Delete(sId)
// 		}
// 		//Delete the link server id to api
// 		this.serverIdToApis.Delete(sId)
// 	}
// 	//Delete the link from api to server id
// 	this.apisToServerId.Delete(hostString)
// }

// func (this *smClient) handleEvent(e SMEvent) {
// 	this.eventFunc(e)
// }

// func apiHostStringFromApi(a api.Api) string {
// 	return fmt.Sprintf("%s:%d", a.HostIP().String(), a.HostPort())
// }

// func getServerId(e event.RegistrationEvent) (string, error) {
// 	var id string
// 	r, err := http.DefaultClient.Get(apiHostStringFromApi(e.Api()) + "/serverid")

// 	if err == nil {
// 		var data []byte
// 		data, err = ioutil.ReadAll(r.Body)
// 		if err == nil {
// 			id = string(data)
// 		}
// 	}
// 	return id, err
// }
