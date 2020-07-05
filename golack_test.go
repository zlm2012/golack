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

	if golack.WebClient != webClient {
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
		WebClient: webClient,
	}
	response, err := golack.PostMessage(context.TODO(), postMessage)

	if err != nil {
		t.Errorf("something is wrong. %#v", err)
	}

	if response.OK != true {
		t.Errorf("OK status is wrong. %#v", response)
	}
}
