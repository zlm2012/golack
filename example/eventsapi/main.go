package main

import (
	"github.com/oklahomer/golack/v2/event"
	"github.com/oklahomer/golack/v2/eventsapi"
	"log"
	"net/http"
	"os"
)

func main() {
	optValidator := eventsapi.WithRequestValidator(&eventsapi.SignatureValidator{Secret: os.Getenv("APP_SECRET")})
	handler := eventsapi.SetupHandler(&Receiver{}, optValidator)
	err := http.ListenAndServe(":8080", handler)
	if err != nil {
		log.Printf(err.Error())
	}
}

type Receiver struct {
}

func (r *Receiver) Receive(wrapper *eventsapi.EventWrapper) {
	switch typed := wrapper.Event.(type) {
	case *event.ChannelMessage:
		// A message in a public channel
		log.Printf("Channel Message: %s", typed.Text)

	default:
		log.Printf("Event: %+v", typed)
	}
}

var _ eventsapi.EventReceiver = (*Receiver)(nil)
