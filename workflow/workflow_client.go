package workflow

import (
	"context"

	"github.com/taadis/dify-sdk-go/client"
)

type WorkflowClient interface {
	// Execute Workflow
	Run(ctx context.Context, req *RunRequest) (*RunResponse, error)
	// Get Workflow Logs
	GetWorkflowLogs(ctx context.Context, req *GetWorkflowLogsRequest) (*GetWorkflowLogsResponse, error)
	// Get Application Basic Information
	GetInfo(ctx context.Context, req *GetInfoRequest) (*GetInfoResponse, error)
	// Get Application Parameters Information
	GetParameters(ctx context.Context, req *GetParametersRequest) (*GetParametersResponse, error)
	// Stop Workflow Task Generation
	StopTask(ctx context.Context, req *StopTaskRequest) (*StopTaskResponse, error)
	// File Upload for Workflow
	UploadFile(ctx context.Context, req *UploadFileRequest) (*UploadFileResponse, error)
	// Get Application WebApp Settings
	GetSite(ctx context.Context, req *GetSiteRequest) (*GetSiteResponse, error)
}

type workflowClient struct {
	*client.Client
}

func NewWorkflowClient(baseUrl string, apiKey string) WorkflowClient {
	c := new(workflowClient)
	c.Client = client.NewClient(baseUrl, apiKey)
	return c
}
