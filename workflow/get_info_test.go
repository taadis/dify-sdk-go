package workflow

import (
	"context"
	"testing"

	"github.com/taadis/dify-sdk-go/client"
	"github.com/taadis/dify-sdk-go/env"
)

func TestGetInfo(t *testing.T) {
	ctx := context.Background()

	// messages
	req := &GetInfoRequest{}

	client := NewWorkflowClient(client.DifyCloud, env.GetDifyApiKey())
	rsp, err := client.GetInfo(ctx, req)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(rsp.String())
}
