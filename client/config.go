package client

import (
	"net/http"
	"time"
)

type ClientConfig struct {
	Host      string
	ApiKey    string
	Timeout   time.Duration
	Transport *http.Transport
}
