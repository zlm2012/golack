package webapi

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"golang.org/x/xerrors"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"
)

const (
	slackAPIEndpointFormat = "https://slack.com/api/%s"
)

type URLValuer interface {
	ToURLValues() url.Values
}

type Config struct {
	Token          string        `json:"token" yaml:"token"`
	RequestTimeout time.Duration `json:"request_timeout" yaml:"request_timeout"`
}

func NewConfig() *Config {
	return &Config{
		Token:          "",
		RequestTimeout: 3 * time.Second,
	}
}

type ClientOption func(*Client)

func WithHTTPClient(httpClient *http.Client) ClientOption {
	return func(c *Client) {
		c.httpClient = httpClient
	}
}

type Client struct {
	config     *Config
	httpClient *http.Client
}

func NewClient(config *Config, options ...ClientOption) *Client {
	c := &Client{config: config}
	for _, opt := range options {
		opt(c)
	}

	if c.httpClient == nil {
		c.httpClient = http.DefaultClient
	}

	return c
}

func buildEndpoint(slackMethod string, queryParams url.Values) *url.URL {
	requestURL, err := url.Parse(fmt.Sprintf(slackAPIEndpointFormat, slackMethod))
	if err != nil {
		panic(err.Error())
	}

	if queryParams != nil {
		requestURL.RawQuery = queryParams.Encode()
	}

	return requestURL
}

func (client *Client) Get(ctx context.Context, slackMethod string, queryParams url.Values, response interface{}) error {
	// Prepare request
	endpoint := buildEndpoint(slackMethod, queryParams)
	req, err := http.NewRequest(http.MethodGet, endpoint.String(), nil)
	if err != nil {
		return err
	}
	reqCtx, cancel := context.WithTimeout(ctx, client.config.RequestTimeout)
	defer cancel()
	req.WithContext(reqCtx)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", client.config.Token))

	// Do request
	resp, err := client.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Usually, the API returns a JSON structure with status code 200.
	// https://api.slack.com/web#evaluating_responses
	if resp.StatusCode != http.StatusOK {
		return statusErr(resp)
	}

	// Handle response body
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return err
	}
	return nil
}

func statusErr(resp *http.Response) error {
	reqDump := []byte("N/A")
	if resp.Request != nil {
		dump, err := httputil.DumpRequestOut(resp.Request, true)
		if err != nil {
			reqDump = dump
		}
	}

	resDump, resErr := httputil.DumpResponse(resp, true)
	if resErr != nil {
		resDump = []byte("N/A")
	}

	return fmt.Errorf("response status error. Status: %d.\nRequest: %s\nResponse: %s", resp.StatusCode, string(reqDump), string(resDump))
}

func (client *Client) Post(ctx context.Context, slackMethod string, payload interface{}, response interface{}) error {
	// Decide how the request should be treated depending on the slackMethod/payload
	p, err := genPayload(slackMethod, payload)
	if err != nil {
		return err
	}

	// Prepare request
	endpoint := buildEndpoint(slackMethod, nil)
	req, err := http.NewRequest("POST", endpoint.String(), bytes.NewReader(p.Body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", p.Type)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", client.config.Token))
	reqCtx, cancel := context.WithTimeout(ctx, client.config.RequestTimeout)
	defer cancel()
	req.WithContext(reqCtx)

	// Do request
	resp, err := client.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Usually, the API returns a JSON structure with status code 200.
	// https://api.slack.com/web#evaluating_responses
	if resp.StatusCode != http.StatusOK {
		return statusErr(resp)
	}

	// Handle response body
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return err
	}
	return nil
}

func genPayload(m string, p interface{}) (*payload, error) {
	// Other than specifically designed methods, most payloads must be sent as application/x-www-form-urlencoded
	// See https://api.slack.com/web#basics
	switch typed := p.(type) {
	case url.Values:
		return &payload{
			Type: "application/x-www-form-urlencoded",
			Body: []byte(typed.Encode()),
		}, nil

	case URLValuer:
		// When the endpoint does not support JSON formatted payload, the incoming struct should implement URLValuer.
		// See how PostMessage worked back in https://github.com/oklahomer/golack/blob/75ad9b2a4ace063b033241c514458789056bd874/webapi/request.go#L113-L144
		return &payload{
			Type: "application/x-www-form-urlencoded",
			Body: []byte(typed.ToURLValues().Encode()),
		}, nil

	default:
		// Send JSON formatted payload when the method supports: https://api.slack.com/web#methods_supporting_json
		supported := IsJSONPayloadSupportedMethod(m)
		if !supported {
			return nil, ErrJSONPayloadNotSupported
		}

		b, err := json.Marshal(p)
		if err != nil {
			return nil, xerrors.Errorf("failed to serialize payload: %w", err)
		}
		return &payload{
			Type: "application/json",
			Body: b,
		}, nil
	}
}

type payload struct {
	Type string
	Body []byte
}
