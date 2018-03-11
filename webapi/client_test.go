package webapi

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"golang.org/x/net/context"
)

type GetResponseDummy struct {
	APIResponse
	Foo string
}

func TestWithHTTPClient(t *testing.T) {
	httpClient := &http.Client{}
	option := WithHTTPClient(httpClient)
	client := &Client{}

	option(client)

	if client.httpClient != httpClient {
		t.Errorf("Specified htt.Client is not set")
	}
}

func TestNewClient(t *testing.T) {
	config := &Config{Token: "abc", RequestTimeout: 1 * time.Second}
	optionCalled := false
	client := NewClient(config, func(*Client) { optionCalled = true })

	if client == nil {
		t.Fatal("Returned client is nil.")
	}

	if client.config != config {
		t.Errorf("Returned client does not have assigned config: %#v.", client.config)
	}

	if !optionCalled {
		t.Error("ClientOption is not called.")
	}

	if client.httpClient != http.DefaultClient {
		t.Errorf("When WithHTTPClient is not given, http.DefaultClient must be set: %+v", client.httpClient)
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
	returningResponse := &GetResponseDummy{
		APIResponse: APIResponse{OK: true},
		Foo:         "bar",
	}
	tripper := buildRoundTripper(http.MethodGet, "/api/foo", returningResponse)
	client := &Client{
		config: &Config{
			Token:          "abc",
			RequestTimeout: 3 * time.Second,
		},
		httpClient: &http.Client{Transport: tripper},
	}

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
	mux := http.NewServeMux()
	mux.HandleFunc("/api/foo", func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	})
	tripper := &localRoundTripper{mux: mux}
	client := &Client{
		config: &Config{
			Token:          "abc",
			RequestTimeout: 3 * time.Second,
		},
		httpClient: &http.Client{Transport: tripper},
	}

	err := client.Get(context.TODO(), "foo", nil, &GetResponseDummy{})

	if err == nil {
		t.Error("error should return when 404 is given.")
	}
}

func TestClient_Get_JSONError(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/foo", func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("invalid json"))
	})
	tripper := &localRoundTripper{mux: mux}
	client := &Client{
		config: &Config{
			Token:          "abc",
			RequestTimeout: 3 * time.Second,
		},
		httpClient: &http.Client{Transport: tripper},
	}

	err := client.Get(context.TODO(), "foo", nil, &GetResponseDummy{})

	if err == nil {
		t.Error("error should return")
	}
}

func TestClient_Post(t *testing.T) {
	tripper := buildRoundTripper(http.MethodPost, "/api/foo", &APIResponse{OK: true})
	client := &Client{
		config: &Config{
			Token:          "abc",
			RequestTimeout: 3 * time.Second,
		},
		httpClient: &http.Client{Transport: tripper},
	}

	returnedResponse := &APIResponse{}
	err := client.Post(context.TODO(), "foo", url.Values{}, returnedResponse)

	if err != nil {
		t.Errorf("something is wrong. %#v", err)
	}
}

func TestPostStatusError(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/foo", func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	})
	tripper := &localRoundTripper{mux: mux}
	client := &Client{
		config: &Config{
			Token:          "abc",
			RequestTimeout: 3 * time.Second,
		},
		httpClient: &http.Client{Transport: tripper},
	}

	returnedResponse := &APIResponse{}
	err := client.Post(context.TODO(), "foo", url.Values{}, returnedResponse)

	if err == nil {
		t.Error("error should return when 500 is given.")
	}
}

func buildRoundTripper(method string, path string, response interface{}) *localRoundTripper {
	mux := http.NewServeMux()
	mux.HandleFunc(path, func(w http.ResponseWriter, req *http.Request) {
		if req.Method != method {
			code := http.StatusMethodNotAllowed
			http.Error(w, http.StatusText(code), code)
		}
		bytes, err := json.Marshal(response)
		if err != nil {
			code := http.StatusInternalServerError
			http.Error(w, http.StatusText(code), code)
		}

		w.WriteHeader(http.StatusOK)
		_, err = w.Write(bytes)
		if err != nil {
			code := http.StatusInternalServerError
			http.Error(w, http.StatusText(code), code)
		}
	})
	return &localRoundTripper{mux: mux}
}

type localRoundTripper struct {
	mux *http.ServeMux
}

func (l *localRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	w := httptest.NewRecorder()
	l.mux.ServeHTTP(w, req)
	return w.Result(), nil
}
