package main

import (
	"context"
	"fmt"
	"github.com/oklahomer/golack/v2/event"
	"github.com/oklahomer/golack/v2/rtmapi"
	"github.com/oklahomer/golack/v2/webapi"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())

	// Setup Web API client
	config := webapi.NewConfig()
	config.Token = os.Getenv("APP_TOKEN")
	client := webapi.NewClient(config)

	// Retrieve RTM session information from Web API
	rtmStart := &webapi.RTMStart{}
	err := client.Get(ctx, "rtm.start", nil, rtmStart)
	if err != nil {
		panic(err)
	}
	if rtmStart.OK != true {
		panic(fmt.Errorf("failed rtm.start request: %s", rtmStart.Error))
	}

	// Create WebSocket connection with the returned URL
	conn, err := rtmapi.Connect(ctx, rtmStart.URL)
	if err != nil {
		panic(err)
	}

	// Receive events over WebSocket connection
	go func() {
		for {
			select {
			case <-ctx.Done():
				return

			default:
				payload, err := conn.Receive()

				if err == event.ErrEmptyPayload {
					continue
				}

				if err != nil {
					fmt.Printf("ERROR: %+v", err)
					continue
				}

				fmt.Printf("PAYLOAD: %+v", payload)
			}
		}
	}()

	// Stop interaction
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	select {
	case <-c:
		fmt.Println("FINISH")
		cancel()

	}
}
