package workflow

import (
	"context"

	"github.com/taadis/dify-sdk-go/client"
)

type WorkflowClient interface {
	// Get Workflow Logs
	GetWorkflowLogs(ctx context.Context, req *GetWorkflowLogsRequest) (*GetWorkflowLogsResponse, error)
}

type workflowClient struct {
	*client.Client
}

func NewWorkflowClient(host string, apiKey string) WorkflowClient {
	c := new(workflowClient)
	c.Client = client.NewClient(host, apiKey)
	return c
}
