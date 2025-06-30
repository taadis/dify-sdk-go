package client

import (
	"net/http"
	"time"
)

type ClientConfig struct {
	// Base URL
	BaseUrl string
	// API-Key
	ApiKey string
	//
	Timeout time.Duration
	//
	Transport *http.Transport
}
