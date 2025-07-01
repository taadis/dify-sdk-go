package chatflow

import (
	"context"

	"github.com/taadis/dify-sdk-go/client"
)

type ChatflowClient interface {
	// Stop Advanced Chat Message Generation
	Stop(ctx context.Context, req *StopRequest) (*StopResponse, error)
}

type chatflowClient struct {
	*client.Client
}

func NewChatflowClient(baseUrl string, apiKey string) ChatflowClient {
	c := new(chatflowClient)
	c.Client = client.NewClient(baseUrl, apiKey)
	return c
}
