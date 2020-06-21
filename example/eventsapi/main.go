package main

import (
	"github.com/oklahomer/golack/event"
	"github.com/oklahomer/golack/eventsapi"
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
	case *event.MessageGroups:
		// A message usually sent in a group
		log.Printf("Group Message: %s", typed.Text)

	case *event.MessageIM:
		// A message directly sent to the app
		log.Printf("IM Message: %s", typed.Text)

	default:
		log.Printf("Event: %+v", typed)
	}
}

var _ eventsapi.EventReceiver = (*Receiver)(nil)
