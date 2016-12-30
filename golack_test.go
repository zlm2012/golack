package golack

import (
	"github.com/jarcoal/httpmock"
	"github.com/oklahomer/golack/webapi"
	"golang.org/x/net/context"
	"testing"
)

func TestGolack_StartRTMSession(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	jsonResponder, _ := httpmock.NewJsonResponder(
		200,
		&webapi.RTMStart{
			APIResponse: webapi.APIResponse{OK: true},
			URL:         "https://localhost/foo",
			Self:        nil})

	httpmock.RegisterResponder(
		"GET",
		"https://slack.com/api/rtm.start",
		jsonResponder)

	golack := New(NewConfig())
	rtmStart, err := golack.StartRTMSession(context.TODO())

	if err != nil {
		t.Errorf("something went wrong. %#v", err)
	}

	if rtmStart.URL != "https://localhost/foo" {
		t.Errorf("URL is not returned properly. %#v", rtmStart)
	}
}

func TestGolack_PostMessage(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	response := &webapi.APIResponse{OK: true}
	jsonResponder, _ := httpmock.NewJsonResponder(200, response)
	httpmock.RegisterResponder(
		"POST",
		"https://slack.com/api/chat.postMessage",
		jsonResponder)

	postMessage := webapi.NewPostMessage("channel", "my message")
	golack := New(NewConfig())
	response, err := golack.PostMessage(context.TODO(), postMessage)

	if err != nil {
		t.Errorf("something is wrong. %#v", err)
	}

	if response.OK != true {
		t.Errorf("OK status is wrong. %#v", response)
	}
}
