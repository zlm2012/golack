package golack

import (
	"context"
	"github.com/oklahomer/golack/webapi"
	"net/url"
	"testing"
)

type DummyWebClient struct {
	GetFunc  func(ctx context.Context, slackMethod string, queryParams *url.Values, intf interface{}) error
	PostFunc func(ctx context.Context, slackMethod string, bodyParam url.Values, intf interface{}) error
}

func (wc *DummyWebClient) Get(ctx context.Context, slackMethod string, queryParams *url.Values, intf interface{}) error {
	return wc.GetFunc(ctx, slackMethod, queryParams, intf)
}

func (wc *DummyWebClient) Post(ctx context.Context, slackMethod string, bodyParam url.Values, intf interface{}) error {
	return wc.PostFunc(ctx, slackMethod, bodyParam, intf)
}

func TestWithWebClient(t *testing.T) {
	webClient := &DummyWebClient{}
	option := WithWebClient(webClient)
	golack := &Golack{}

	option(golack)

	if golack.webClient != webClient {
		t.Errorf("Specified WebClient is not set.")
	}
}

func TestNew(t *testing.T) {
	config := &Config{}
	optionCalled := false

	golack := New(config, func(_ *Golack) { optionCalled = true })

	if golack == nil {
		t.Fatal("Returned *Golack is nil.")
	}

	if !optionCalled {
		t.Error("Option is not called.")
	}
}

func TestGolack_StartRTMSession(t *testing.T) {
	webClient := &DummyWebClient{
		GetFunc: func(_ context.Context, slackMethod string, _ *url.Values, intf interface{}) error {
			if slackMethod != "rtm.start" {
				t.Errorf("Requesting path is not correct: %s.", slackMethod)
			}
			start := intf.(*webapi.RTMStart)
			start.APIResponse = webapi.APIResponse{OK: true}
			start.URL = "https://localhost/foo"
			start.Self = nil
			return nil
		},
	}
	golack := &Golack{
		webClient: webClient,
	}

	rtmStart, err := golack.StartRTMSession(context.TODO())

	if err != nil {
		t.Errorf("something went wrong. %#v", err)
	}

	if rtmStart.URL != "https://localhost/foo" {
		t.Errorf("URL is not returned properly. %#v", rtmStart)
	}
}

func TestGolack_PostMessage(t *testing.T) {
	webClient := &DummyWebClient{
		PostFunc: func(ctx context.Context, slackMethod string, bodyParam url.Values, intf interface{}) error {
			response := intf.(*webapi.APIResponse)
			response.OK = true
			response.Error = ""
			return nil
		},
	}

	postMessage := webapi.NewPostMessage("channel", "my message")
	golack := &Golack{
		webClient: webClient,
	}
	response, err := golack.PostMessage(context.TODO(), postMessage)

	if err != nil {
		t.Errorf("something is wrong. %#v", err)
	}

	if response.OK != true {
		t.Errorf("OK status is wrong. %#v", response)
	}
}
