package workflow

import (
	"context"
	"testing"

	"github.com/taadis/dify-sdk-go/client"
	"github.com/taadis/dify-sdk-go/env"
)

func TestGetParamaters(t *testing.T) {
	ctx := context.Background()

	req := &GetParametersRequest{}
	client := NewWorkflowClient(client.DifyCloud, env.GetDifyApiKey())
	rsp, err := client.GetParameters(ctx, req)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(rsp.MarshalIndent())
}
