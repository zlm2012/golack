package main

import (
	"context"
	"fmt"
	"github.com/oklahomer/golack"
	"github.com/oklahomer/golack/event"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())

	// Setup Golack client
	config := golack.NewConfig()
	config.Token = os.Getenv("APP_TOKEN")
	client := golack.New(config)

	// Establish WebSocket connection
	conn, err := client.ConnectRTM(ctx)
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
