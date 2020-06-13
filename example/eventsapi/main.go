package main

import (
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
	log.Printf("Recieved: %#v", wrapper)
}

var _ eventsapi.EventReceiver = (*Receiver)(nil)
