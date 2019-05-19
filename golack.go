package golack

import (
	"context"
	"fmt"
	"net/url"
	"time"

	"github.com/oklahomer/golack/rtmapi"
	"github.com/oklahomer/golack/webapi"
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

type WebClient interface {
	Get(ctx context.Context, slackMethod string, queryParams *url.Values, intf interface{}) error
	Post(ctx context.Context, slackMethod string, bodyParam url.Values, intf interface{}) error
}

type Option func(*Golack)

func WithWebClient(wc WebClient) Option {
	return func(g *Golack) {
		g.webClient = wc
	}
}

type Golack struct {
	webClient WebClient
}

func New(config *Config, options ...Option) *Golack {
	g := &Golack{}
	for _, opt := range options {
		opt(g)
	}

	if g.webClient == nil {
		apiConfig := &webapi.Config{
			Token:          config.Token,
			RequestTimeout: config.RequestTimeout,
		}
		g.webClient = webapi.NewClient(apiConfig)
	}

	return g
}

func (g *Golack) StartRTMSession(ctx context.Context) (*webapi.RTMStart, error) {
	rtmStart := &webapi.RTMStart{}
	if err := g.webClient.Get(ctx, "rtm.start", nil, rtmStart); err != nil {
		return nil, err
	}

	if rtmStart.OK != true {
		return nil, fmt.Errorf("failed rtm.start request: %s", rtmStart.Error)
	}

	return rtmStart, nil
}

func (g *Golack) PostMessage(ctx context.Context, postMessage *webapi.PostMessage) (*webapi.APIResponse, error) {
	response := &webapi.APIResponse{}
	err := g.webClient.Post(ctx, "chat.postMessage", postMessage.ToURLValues(), response)
	if err != nil {
		return nil, err
	}

	if response.OK != true {
		return nil, fmt.Errorf("failed chat.postMessage request: %s", response.Error)
	}

	return response, nil
}

// Connect connects to Slack WebSocket server.
func (g *Golack) ConnectRTM(ctx context.Context, url string) (rtmapi.Connection, error) {
	return rtmapi.Connect(ctx, url)
}
