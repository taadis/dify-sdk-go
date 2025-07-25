package dify

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/taadis/http2curl"
)

// from https://docs.dify.ai/en/openapi-api-access-readme#%F0%9F%94%91-api-access-configuration
var DifyCloud = "https://api.dify.ai"

type Client struct {
	host             string
	defaultAPISecret string
	httpClient       *http.Client
}

func NewClientWithConfig(c *ClientConfig) *Client {
	var httpClient = &http.Client{}

	if c.Timeout != 0 {
		httpClient.Timeout = c.Timeout
	}
	if c.Transport != nil {
		httpClient.Transport = c.Transport
	}

	secret := c.DefaultAPISecret
	if secret == "" {
		secret = c.ApiSecretKey
	}
	return &Client{
		host:             c.Host,
		defaultAPISecret: secret,
		httpClient:       httpClient,
	}
}

func NewClient(host, defaultAPISecret string) *Client {
	return NewClientWithConfig(&ClientConfig{
		Host:             host,
		DefaultAPISecret: defaultAPISecret,
	})
}

func (c *Client) sendRequest(req *http.Request) (*http.Response, error) {

	curlcmd, err := http2curl.GetCurlCommand(req)
	if err != nil {
		return nil, err
	}
	fmt.Println(curlcmd.String())

	return c.httpClient.Do(req)
}

func (c *Client) sendJSONRequest(req *http.Request, res interface{}) error {
	resp, err := c.sendRequest(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var errBody struct {
			Code    string `json:"code"`
			Message string `json:"message"`
			Status  int    `json:"status"`
		}
		err = json.NewDecoder(resp.Body).Decode(&errBody)
		if err != nil {
			return err
		}
		return fmt.Errorf("HTTP response error: [%v]%v", errBody.Code, errBody.Message)
	}

	err = json.NewDecoder(resp.Body).Decode(res)
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) getHost() string {
	var host = strings.TrimSuffix(c.host, "/")
	return host
}

func (c *Client) getAPISecret() string {
	return c.defaultAPISecret
}

// Api deprecated, use API() instead
func (c *Client) Api() *API {
	return c.API()
}

func (c *Client) API() *API {
	return &API{
		c: c,
	}
}
