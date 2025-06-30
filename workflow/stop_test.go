package workflow

import (
	"context"
	"testing"

	"github.com/taadis/dify-sdk-go/client"
	"github.com/taadis/dify-sdk-go/env"
)

func TestStopTask(t *testing.T) {
	ctx := context.Background()

	taskId := "your-task-id"
	user := "your-user"
	if taskId == "your-task-id" || user == "your-user" {
		t.Skip("Set a valid task_id and user to run this test.")
	}

	req := &StopTaskRequest{TaskId: taskId, User: user}
	client := NewWorkflowClient(client.DifyCloud, env.GetDifyApiKey())
	rsp, err := client.StopTask(ctx, req)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(rsp.String())
}
