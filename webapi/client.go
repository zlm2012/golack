package webapi

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"time"
)

const (
	slackAPIEndpointFormat = "https://slack.com/api/%s"
)

type Config struct {
	Token          string
	RequestTimeout time.Duration
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

func buildEndpoint(slackMethod, token string, queryParams *url.Values) *url.URL {
	if queryParams == nil {
		queryParams = &url.Values{}
	}
	queryParams.Add("token", token)

	requestURL, err := url.Parse(fmt.Sprintf(slackAPIEndpointFormat, slackMethod))
	if err != nil {
		panic(err.Error())
	}
	requestURL.RawQuery = queryParams.Encode()

	return requestURL
}

func (client *Client) Get(ctx context.Context, slackMethod string, queryParams *url.Values, intf interface{}) error {
	// Prepare request
	endpoint := buildEndpoint(slackMethod, client.config.Token, queryParams)
	req, err := http.NewRequest(http.MethodGet, endpoint.String(), nil)
	if err != nil {
		return err
	}
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
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(body, &intf); err != nil {
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

func (client *Client) Post(ctx context.Context, slackMethod string, bodyParam url.Values, intf interface{}) error {
	// Prepare request
	endpoint := buildEndpoint(slackMethod, client.config.Token, nil)
	req, err := http.NewRequest("POST", endpoint.String(), strings.NewReader(bodyParam.Encode()))
	if err != nil {
		return err
	}
	reqCtx, cancel := context.WithTimeout(ctx, client.config.RequestTimeout)
	defer cancel()
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
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
	response, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(response, &intf); err != nil {
		return err
	}
	return nil
}
