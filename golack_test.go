package golack

import (
	"context"
	"fmt"
	"github.com/oklahomer/golack/testutil"
	"github.com/oklahomer/golack/webapi"
	"golang.org/x/xerrors"
	"net"
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
	g := &Golack{}

	option(g)

	if g.WebClient != webClient {
		t.Errorf("Specified WebClient is not set.")
	}
}

func TestNew(t *testing.T) {
	config := &Config{}
	optionCalled := false

	g := New(config, func(_ *Golack) { optionCalled = true })

	if g == nil {
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
	g := &Golack{
		WebClient: webClient,
	}

	postMessage := webapi.NewPostMessage("channel", "my message")
	response, err := g.PostMessage(context.TODO(), postMessage)

	if err != nil {
		t.Errorf("something is wrong. %#v", err)
	}

	if response.OK != true {
		t.Errorf("OK status is wrong. %#v", response)
	}
}

func TestGolack_ConnectRTM(t *testing.T) {
	t.Run("Web API returns error status", func(t *testing.T) {
		expectedErr := xerrors.New("DUMMY")
		webClient := &DummyWebClient{
			GetFunc: func(_ context.Context, _ string, _ *url.Values, _ interface{}) error {
				return expectedErr
			},
		}
		g := &Golack{
			WebClient: webClient,
		}

		_, err := g.ConnectRTM(context.Background())
		if err == nil {
			t.Fatal("Error is not returned.")
		}
		if err != expectedErr {
			t.Fatalf("Expected error is not returned: %+v", err)
		}
	})

	t.Run("Web API returns error response", func(t *testing.T) {
		webClient := &DummyWebClient{
			GetFunc: func(ctx context.Context, slackMethod string, queryParams *url.Values, intf interface{}) error {
				response := intf.(*webapi.RTMStart)
				response.OK = false
				response.Error = "some error"
				return nil
			},
		}
		g := &Golack{
			WebClient: webClient,
		}

		_, err := g.ConnectRTM(context.Background())
		if err == nil {
			t.Fatal("Expected error is not returned.")
		}
	})

	t.Run("connect WebSocket server", func(t *testing.T) {
		testutil.RunWithWebSocket(func(addr net.Addr) {
			webClient := &DummyWebClient{
				GetFunc: func(ctx context.Context, slackMethod string, queryParams *url.Values, intf interface{}) error {
					response := intf.(*webapi.RTMStart)
					response.OK = true
					response.URL = fmt.Sprintf("ws://%s%s", addr, "/ping")
					response.Error = ""
					return nil
				},
			}
			g := &Golack{
				WebClient: webClient,
			}

			rtm, err := g.ConnectRTM(context.Background())
			if err != nil {
				t.Fatalf("Unexpected error is returned: %s", err.Error())
			}

			err = rtm.Ping()
			if err != nil {
				t.Fatalf("Unexpected error is returned on Ping: %s", err.Error())
			}
		})
	})
}
