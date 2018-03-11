package webapi

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"

	"golang.org/x/net/context"
	"golang.org/x/net/context/ctxhttp"
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
	endpoint := buildEndpoint(slackMethod, client.config.Token, queryParams)

	reqCtx, cancel := context.WithTimeout(ctx, client.config.RequestTimeout)
	defer cancel()
	resp, err := ctxhttp.Get(reqCtx, client.httpClient, endpoint.String())
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	// Usually, the API returns a JSON structure with status code 200.
	// https://api.slack.com/web#evaluating_responses
	if resp.StatusCode != http.StatusOK {
		return statusErr(resp)
	}

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
	endpoint := buildEndpoint(slackMethod, client.config.Token, nil)

	reqCtx, cancel := context.WithTimeout(ctx, client.config.RequestTimeout)
	defer cancel()
	resp, err := ctxhttp.PostForm(reqCtx, client.httpClient, endpoint.String(), bodyParam)
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	// Usually, the API returns a JSON structure with status code 200.
	// https://api.slack.com/web#evaluating_responses
	if resp.StatusCode != http.StatusOK {
		return statusErr(resp)
	}

	response, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(response, &intf); err != nil {
		return err
	}

	return nil
}
