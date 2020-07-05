package golack

import (
	"context"
	"fmt"
	"github.com/oklahomer/golack/eventsapi"
	"golang.org/x/xerrors"
	"net/http"
	"net/url"
	"time"

	"github.com/oklahomer/golack/rtmapi"
	"github.com/oklahomer/golack/webapi"
)

type Config struct {
	AppSecret      string        `json:"app_secret" yaml:"app_secret"`
	Token          string        `json:"token" yaml:"token"`
	ListenPort     int           `json:"listen_port" yaml:"listen_port"`
	RequestTimeout time.Duration `json:"request_timeout" yaml:"request_timeout"`
}

// NewConfig returns initialized Config struct with default settings.
// AppSecret and Token are empty at this point. They can be set/updated by feeding this instance to json.Unmarshal/yaml.Unmarshal
// or by direct assignment.
func NewConfig() *Config {
	return &Config{
		AppSecret:      "",
		Token:          "",
		ListenPort:     8080,
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
		g.WebClient = wc
	}
}

type Golack struct {
	config    *Config
	WebClient WebClient
}

func New(config *Config, options ...Option) *Golack {
	g := &Golack{}
	g.config = config

	// Apply options to change specific behaviors
	for _, opt := range options {
		opt(g)
	}

	// If WebClient is not set with Option, then built one with default settings
	if g.WebClient == nil {
		apiConfig := webapi.NewConfig()
		apiConfig.Token = g.config.Token
		apiConfig.RequestTimeout = g.config.RequestTimeout
		g.WebClient = webapi.NewClient(apiConfig)
	}

	return g
}

func (g *Golack) PostMessage(ctx context.Context, postMessage *webapi.PostMessage) (*webapi.APIResponse, error) {
	response := &webapi.APIResponse{}
	err := g.WebClient.Post(ctx, "chat.postMessage", postMessage.ToURLValues(), response)
	if err != nil {
		return nil, err
	}

	if response.OK != true {
		return nil, fmt.Errorf("failed chat.postMessage request: %s", response.Error)
	}

	return response, nil
}

// ConnectRTM connects to Slack WebSocket server.
func (g *Golack) ConnectRTM(ctx context.Context) (rtmapi.Connection, error) {
	rtmStart := &webapi.RTMStart{}
	if err := g.WebClient.Get(ctx, "rtm.start", nil, rtmStart); err != nil {
		return nil, err
	}

	if rtmStart.OK != true {
		return nil, fmt.Errorf("failed rtm.start request: %s", rtmStart.Error)
	}

	return rtmapi.Connect(ctx, rtmStart.URL)
}

func (g *Golack) RunServer(ctx context.Context, receiver eventsapi.EventReceiver) <-chan error {
	errChan := make(chan error, 1)

	// Setup a request validator
	// For better security, this checks each request's signature
	appSecret := g.config.AppSecret
	if appSecret == "" {
		errChan <- xerrors.New("application secret is not set")
		return errChan
	}
	optValidator := eventsapi.WithRequestValidator(&eventsapi.SignatureValidator{Secret: appSecret})

	// Setup server and run it
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", g.config.ListenPort),
		Handler: eventsapi.SetupHandler(receiver, optValidator),
	}
	go func() {
		errChan <- srv.ListenAndServe()
	}()

	// Shutdown the server
	go func() {
		<-ctx.Done()
		//noinspection ALL
		srv.Shutdown(ctx)
	}()

	return errChan
}
