package switchmachine

import (
	"encoding/json"
	"io"
	"log"
	"net"
	"net/url"

	"github.com/ZacharyDuve/eventsocket"
)

const (
	eventHandlerFullPath = smHandlerPath + eventHandlerSubPath
)

type SMEventClient interface {
	io.Closer
	SetHandleSMAddFunc(func(SMEvent))
	SetHandleSMRemoveFunc(func(SMEvent))
	SetHandleSMUpdateFunc(func(SMEvent))
}

type smEventClient struct {
	eClient    io.Closer
	addFunc    func(SMEvent)
	removeFunc func(SMEvent)
	updateFunc func(SMEvent)
}

func NewSMEventClient(a net.Addr) (SMEventClient, error) {
	client := &smEventClient{}
	c, e := eventsocket.DialEventServer(&url.URL{Host: a.String(), Path: eventHandlerFullPath}, client.handleEvent)
	if e == nil {
		client.eClient = c
	}

	return client, e
}

func (this *smEventClient) handleEvent(e eventsocket.Event) {
	if isValidEvent(e) {
		jsonSMData := &jsonSMEventData{}
		err := json.Unmarshal([]byte(e.Data()), jsonSMData)
		if err != nil {
			log.Println(err)
			return
		}
		smE := &smEvent{e: e, id: jsonSMData.Id, pos: jsonSMData.Position}
		if this.addFunc != nil && e.Name() == SMAdded {
			this.addFunc(smE)
		} else if this.updateFunc != nil && e.Name() == SMUpdated {
			this.updateFunc(smE)
		} else if this.removeFunc != nil && e.Name() == SMRemoved {
			this.removeFunc(smE)
		}
	}

}

func isValidEvent(e eventsocket.Event) bool {
	return e.Name() == SMAdded || e.Name() == SMUpdated || e.Name() == SMRemoved
}

func (this *smEventClient) Close() error {
	return this.eClient.Close()
}

func (this *smEventClient) SetHandleSMAddFunc(aFunc func(SMEvent)) {
	this.addFunc = aFunc
}
func (this *smEventClient) SetHandleSMRemoveFunc(rFunc func(SMEvent)) {
	this.removeFunc = rFunc
}
func (this *smEventClient) SetHandleSMUpdateFunc(uFunc func(SMEvent)) {
	this.updateFunc = uFunc
}
