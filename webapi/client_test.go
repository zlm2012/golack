package webapi

import (
	"github.com/jarcoal/httpmock"
	"golang.org/x/net/context"
	"net/http"
	"net/url"
	"testing"
	"time"
)

type GetResponseDummy struct {
	APIResponse
	Foo string
}

func TestNewClient(t *testing.T) {
	config := &Config{Token: "abc", RequestTimeout: 1 * time.Second}
	client := NewClient(config)

	if client == nil {
		t.Fatal("client is nil.")
	}

	if client.config != config {
		t.Errorf("returned client does not have assigned config: %#v.", client.config)
	}

}

func Test_buildEndpoint(t *testing.T) {
	token := "abc"
	params := &url.Values{
		"foo": []string{"bar", "buzz"},
	}

	method := "rtm.start"
	endpoint := buildEndpoint(method, token, params)

	if endpoint == nil {
		t.Fatal("url is not returned.")
	}

	fooParam, _ := endpoint.Query()["foo"]
	if fooParam == nil || fooParam[0] != "bar" || fooParam[1] != "buzz" {
		t.Errorf("expected query parameter was not returned: %#v.", fooParam)
	}

	if endpoint.Query().Get("token") != token {
		t.Errorf("expected token is not returned: %s.", endpoint.Query().Get("token"))
	}
}

func TestClient_Get(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	returningResponse := &GetResponseDummy{
		APIResponse: APIResponse{OK: true},
		Foo:         "bar",
	}
	responder, _ := httpmock.NewJsonResponder(200, returningResponse)
	httpmock.RegisterResponder("GET", "https://slack.com/api/foo", responder)

	client := &Client{config: &Config{Token: "abc"}}
	returnedResponse := &GetResponseDummy{}
	err := client.Get(context.TODO(), "foo", nil, returnedResponse)

	if err != nil {
		t.Errorf("something went wrong. %#v", err)
	}

	if returnedResponse.OK != true {
		t.Errorf("OK status is wrong. %#v", returnedResponse)
	}

	if returnedResponse.Foo != "bar" {
		t.Errorf("foo value is wrong. %#v", returnedResponse)
	}
}

func TestClient_Get_StatusError(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	statusCode := 404
	responder := func(req *http.Request) (*http.Response, error) {
		resp := httpmock.NewStringResponse(statusCode, "foo bar")
		resp.Request = req // To let *http.Response.Request work
		return resp, nil
	}
	httpmock.RegisterResponder("GET", "https://slack.com/api/foo", responder)

	client := &Client{config: &Config{Token: "abc"}}
	returnedResponse := &GetResponseDummy{}
	err := client.Get(context.TODO(), "foo", nil, returnedResponse)

	if err == nil {
		t.Errorf("error should return when %d is given.", statusCode)
	}
}

func TestClient_Get_JSONError(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	responder := httpmock.NewStringResponder(200, "invalid json")
	httpmock.RegisterResponder("GET", "https://slack.com/api/foo", responder)

	client := &Client{config: &Config{Token: "abc"}}
	returnedResponse := &GetResponseDummy{}
	err := client.Get(context.TODO(), "foo", nil, returnedResponse)

	if err == nil {
		t.Error("error should return")
	}
}

func TestClient_Post(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	responder, _ := httpmock.NewJsonResponder(200, &APIResponse{OK: true})
	httpmock.RegisterResponder("POST", "https://slack.com/api/foo", responder)

	client := &Client{config: &Config{Token: "abc"}}
	returnedResponse := &APIResponse{}
	err := client.Post(context.TODO(), "foo", url.Values{}, returnedResponse)

	if err != nil {
		t.Errorf("something is wrong. %#v", err)
	}
}

func TestPostStatusError(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	statusCode := 404
	responder := func(req *http.Request) (*http.Response, error) {
		resp := httpmock.NewStringResponse(statusCode, "foo bar")
		resp.Request = req // To let *http.Response.Request work
		return resp, nil
	}
	httpmock.RegisterResponder("POST", "https://slack.com/api/foo", responder)

	client := &Client{config: &Config{Token: "abc"}}
	returnedResponse := &APIResponse{}
	err := client.Post(context.TODO(), "foo", url.Values{}, returnedResponse)

	if err == nil {
		t.Errorf("error should return when %d is given.", statusCode)
	}
}
