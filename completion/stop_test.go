package completion

import (
	"context"
	"testing"
)

func TestStop(t *testing.T) {
	ctx := context.Background()

	req := &StopRequest{}
	req.TaskId = "test-task-id"
	req.User = "test-user"
	client := NewCompletionClient(testBaseUrl, testApiKey)
	rsp, err := client.Stop(ctx, req)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(rsp.MarshalIndent())
}
