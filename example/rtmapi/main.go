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
	config := golack.NewConfig()
	config.Token = os.Getenv("TOKEN")

	client := golack.New(config)
	ctx, cancel := context.WithCancel(context.Background())
	session, err := client.StartRTMSession(ctx)
	if err != nil {
		panic(err)
	}

	conn, err := client.ConnectRTM(ctx, session.URL)
	if err != nil {
		panic(err)
	}

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
				}

				fmt.Printf("PAYLOAD: %+v", payload)
			}
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)

	select {
	case <-c:
		fmt.Println("FINISH")
		cancel()

	}
}
