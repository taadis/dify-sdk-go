package workflow

import (
	"context"
	"testing"
)

func TestGetRun(t *testing.T) {
	ctx := context.Background()

	workflowRunId := "your-workflow-run-id"
	if workflowRunId == "your-workflow-run-id" {
		t.Skip("Set a valid workflow_run_id to run this test.")
	}

	req := &GetRunRequest{WorkflowRunId: workflowRunId}
	client := NewWorkflowClient(testBaseUrl, testApiKey)
	rsp, err := client.GetRun(ctx, req)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(rsp.MarshalIndent())
}
