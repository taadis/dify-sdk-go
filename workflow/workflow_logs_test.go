package workflow

import (
	"context"
	"testing"

	"github.com/taadis/dify-sdk-go/client"
	"github.com/taadis/dify-sdk-go/env"
)

func TestGetWorkflowLogs(t *testing.T) {
	ctx := context.Background()
	req := new(GetWorkflowLogsRequest)
	req.Keyword = "test"
	req.Page = 1
	req.Limit = 1
	req.Status = WorkflowStatusFailed

	client := NewWorkflowClient(client.DifyCloud, env.GetDifyApiKey())
	rsp, err := client.GetWorkflowLogs(ctx, req)
	if err != nil {
		t.Fatalf("GetWorkflowsLogs error: %v", err)
	}
	t.Log(rsp.String())
}
