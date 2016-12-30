package golack

import (
	"fmt"
	"github.com/oklahomer/golack/rtmapi"
	"github.com/oklahomer/golack/webapi"
	"golang.org/x/net/context"
	"time"
)

type Config struct {
	Token          string        `json:"token" yaml:"token"`
	RequestTimeout time.Duration `json:"request_timeout" yaml:"request_timeout"`
}

// NewConfig returns initialized Config struct with default settings.
// Token is empty at this point. Token and other settings can be set/updated by feeding this instance to json.Unmarshal/yaml.Unmarshal,
// or assigning directly.
func NewConfig() *Config {
	return &Config{
		Token:          "",
		RequestTimeout: 3 * time.Second,
	}
}

type Golack struct {
	webClient *webapi.Client
}

func New(config *Config) *Golack {
	webClient := webapi.NewClient(&webapi.Config{Token: config.Token, RequestTimeout: config.RequestTimeout})
	return &Golack{
		webClient: webClient,
	}
}

func (g *Golack) StartRTMSession(ctx context.Context) (*webapi.RTMStart, error) {
	rtmStart := &webapi.RTMStart{}
	if err := g.webClient.Get(ctx, "rtm.start", nil, &rtmStart); err != nil {
		return nil, err
	}

	if rtmStart.OK != true {
		return nil, fmt.Errorf("Error on rtm.start : %s", rtmStart.Error)
	}

	return rtmStart, nil
}

func (g *Golack) PostMessage(ctx context.Context, postMessage *webapi.PostMessage) (*webapi.APIResponse, error) {
	response := &webapi.APIResponse{}
	err := g.webClient.Post(ctx, "chat.postMessage", postMessage.ToURLValues(), &response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

// Connect connects to Slack WebSocket server.
func (g *Golack) ConnectRTM(ctx context.Context, url string) (rtmapi.Connection, error) {
	return rtmapi.Connect(ctx, url)
}
