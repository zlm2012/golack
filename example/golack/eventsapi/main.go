package main

import (
	"context"
	"fmt"
	"github.com/oklahomer/golack/v2"
	"github.com/oklahomer/golack/v2/event"
	"github.com/oklahomer/golack/v2/eventsapi"
	"github.com/oklahomer/golack/v2/webapi"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

func main() {
	config := golack.NewConfig()
	config.AppSecret = os.Getenv("APP_SECRET")
	config.Token = os.Getenv("APP_TOKEN")

	ctx, cancel := context.WithCancel(context.Background())

	g := golack.New(config)
	errChan := g.RunServer(ctx, &Receiver{client: g})

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)

	select {
	case <-c:
		fmt.Println("FINISH")
		cancel()

	case err := <-errChan:
		fmt.Printf("ERROR: %s\n", err.Error())
		cancel()
	}
}

type Receiver struct {
	client *golack.Golack
}

func (r *Receiver) Receive(wrapper *eventsapi.EventWrapper) {
	switch typed := wrapper.Event.(type) {
	case *event.MessageChannels:
		// A message is sent to a public channel
		log.Printf("Channel Message: %+v", typed)
		echoMsg := EchoMessage(typed.Text)
		if echoMsg == "" {
			return
		}

		message := webapi.NewPostMessage(typed.ChannelID, echoMsg).
			WithBlocks([]event.Block{
				event.NewContextBlock([]event.BlockElement{
					event.NewPlainTextObjectBlockElement("Hello! Wanna grab a *beer*? :beer:").
						WithEmoji(true),
				}),
				event.NewDividerBlock(),
				event.NewSectionBlock(event.NewMarkdownTextCompositionObject("Below message was sent by the user as *.echo* command")),
				event.NewDividerBlock(),
				event.NewSectionBlock(event.NewPlainTextCompositionObject(echoMsg)).
					WithFields([]*event.TextCompositionObject{
						event.NewPlainTextCompositionObject("User ID"),
						event.NewPlainTextCompositionObject(string(typed.UserID)),
						event.NewPlainTextCompositionObject("Channel ID"),
						event.NewPlainTextCompositionObject(string(typed.ChannelID)),
						event.NewPlainTextCompositionObject("Time"),
						event.NewPlainTextCompositionObject(typed.EventTimeStamp.Time.String()),
					}),
			})

		log.Printf("Payload: %+v\n", message)
		response, err := r.client.PostMessage(context.TODO(), message)
		if err != nil {
			log.Printf("Post error: %s", err.Error())
			return
		}

		if !response.OK {
			log.Printf("Response error: %s", response.Error)
			return
		}

	case *event.MessageGroups:
		// A message is sent to a group
		log.Printf("Group Message: %+v", typed)

	case *event.MessageIM:
		// A message directly sent to the app
		log.Printf("IM Message: %+v", typed)

	default:
		log.Printf("Event: %T, %+v", typed, typed)
	}
}

var _ eventsapi.EventReceiver = (*Receiver)(nil)

func EchoMessage(text string) string {
	echoMsg := strings.TrimPrefix(text, ".echo ")
	if echoMsg == text {
		// Original text did not starts with echo command
		return ""
	}

	return echoMsg
}
