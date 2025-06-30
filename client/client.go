package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/taadis/http2curl"
)

// from https://docs.dify.ai/en/openapi-api-access-readme#%F0%9F%94%91-api-access-configuration
var DifyCloud = "https://api.dify.ai/v1"

type Client struct {
	// Base URL
	baseUrl string
	// Service API-KEY
	apiKey string
	// HTTP Client
	httpClient *http.Client
}

func NewClientWithConfig(c *ClientConfig) *Client {
	var httpClient = &http.Client{}

	if c.Timeout != 0 {
		httpClient.Timeout = c.Timeout
	}
	if c.Transport != nil {
		httpClient.Transport = c.Transport
	}

	return &Client{
		baseUrl:    c.BaseUrl,
		apiKey:     c.ApiKey,
		httpClient: httpClient,
	}
}

func NewClient(baseUrl string, apiKey string) *Client {
	return NewClientWithConfig(&ClientConfig{
		BaseUrl: baseUrl,
		ApiKey:  apiKey,
	})
}

func (c *Client) SendRequest(req *http.Request) (*http.Response, error) {
	return c.sendRequest(req)
}

func (c *Client) sendRequest(req *http.Request) (*http.Response, error) {
	curlcmd, err := http2curl.GetCurlCommand(req)
	if err != nil {
		return nil, err
	}
	fmt.Println(curlcmd.String())

	return c.httpClient.Do(req)
}

func (c *Client) SendJSONRequest(req *http.Request, res interface{}) error {
	return c.sendJSONRequest(req, res)
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

func (c *Client) getBaseUrl() string {
	var baseUrl = strings.TrimSuffix(c.baseUrl, "/")
	return baseUrl
}

func (c *Client) getApiKey() string {
	return c.apiKey
}

func (c *Client) WithApiKey(apiKey string) *Client {
	c.apiKey = apiKey
	return c
}

func (c *Client) CreateBaseRequest(ctx context.Context, method string, apiUrl string, body interface{}) (*http.Request, error) {
	var b io.Reader
	if body != nil {
		reqBytes, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		b = bytes.NewBuffer(reqBytes)
	} else {
		b = http.NoBody
	}

	req, err := http.NewRequestWithContext(ctx, method, c.getBaseUrl()+apiUrl, b)
	if err != nil {
		return nil, err
	}
	fmt.Println("got api key=", c.getApiKey())
	req.Header.Set("Authorization", "Bearer "+c.getApiKey())
	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	return req, nil
}
