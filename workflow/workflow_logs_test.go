package workflow

import (
	"context"
	"testing"
	"time"

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

	now := time.Now()
	todayMidnight := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	tomorrowMidnight := todayMidnight.Add(24 * time.Hour)

	req.CreatedAtAfter = &todayMidnight
	req.CreatedAtBefore = &tomorrowMidnight

	client := NewWorkflowClient(client.DifyCloud, env.GetDifyApiKey())
	rsp, err := client.GetWorkflowLogs(ctx, req)
	if err != nil {
		t.Fatalf("GetWorkflowsLogs error: %v", err)
		// maybe error:
		// [invalid_param]Unused components in ISO string
	}
	t.Log(rsp.String())
	t.Logf("rsp.length: %d", len(rsp.Data))
}
