package webapi

import (
	"encoding/json"
	"fmt"
	"golang.org/x/net/context"
	"golang.org/x/net/context/ctxhttp"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"
)

const (
	slackAPIEndpointFormat = "https://slack.com/api/%s"
)

type Config struct {
	Token          string
	RequestTimeout time.Duration
}

type Client struct {
	config *Config
}

func NewClient(config *Config) *Client {
	return &Client{config: config}
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
	resp, err := ctxhttp.Get(reqCtx, http.DefaultClient, endpoint.String())
	if err != nil {
		return err
	}

	defer resp.Body.Close()
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
	reqDump, reqErr := httputil.DumpRequestOut(resp.Request, true)
	if reqErr != nil {
		reqDump = []byte("N/A")
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
	resp, err := ctxhttp.PostForm(reqCtx, http.DefaultClient, endpoint.String(), bodyParam)
	if err != nil {
		return err
	}

	defer resp.Body.Close()
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
