package workflow

import (
	"context"
	"testing"
)

func TestRun(t *testing.T) {
	ctx := context.Background()
	client := NewWorkflowClient(testBaseUrl, testApiKey)

	// run with image
	workflowReq := &RunRequest{
		Inputs: map[string]interface{}{
			"image_url_new": map[string]string{
				"type":            "image",
				"transfer_method": "remote_url",
				"url":             "https://localhost/1-1.jpg",
			},
		},
		ResponseMode: "blocking",
		User:         "test-user",
	}

	resp, err := client.Run(ctx, workflowReq)
	if err != nil {
		t.Fatalf("RunWorkflow encountered an error: %v", err)
	}
	t.Log(resp.String())

	// check basic rsp
	if resp.WorkflowRunId == "" {
		t.Errorf("Expected non-empty WorkflowRunID, got empty")
	}
	if resp.TaskId == "" {
		t.Errorf("Expected non-empty TaskID, got empty")
	}

	// check status
	if resp.Data.Status != "succeeded" {
		t.Errorf("Expected workflow status 'succeeded', got: %v", resp.Data.Status)
	}

	// check outputs and metadata
	if len(resp.Data.Outputs) == 0 {
		t.Errorf("Expected outputs, but got none")
	}
	if resp.Data.ElapsedTime <= 0 {
		t.Errorf("Expected positive ElapsedTime, but got: %v", resp.Data.ElapsedTime)
	}
	if resp.Data.TotalSteps <= 0 {
		t.Errorf("Expected positive TotalSteps, but got: %v", resp.Data.TotalSteps)
	}

	t.Logf("Received workflow response: %+v", resp)
}
