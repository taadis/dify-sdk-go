package completion

import (
	"context"

	"github.com/taadis/dify-sdk-go/client"
)

type CompletionClient interface {
	// Stop Generate
	Stop(ctx context.Context, req *StopRequest) (*StopResponse, error)
}

type completionClient struct {
	*client.Client
}

func NewCompletionClient(baseUrl string, apiKey string) CompletionClient {
	c := new(completionClient)
	c.Client = client.NewClient(baseUrl, apiKey)
	return c
}
