package workflow

import (
	"context"
	"testing"
)

func TestSiteInfo(t *testing.T) {
	ctx := context.Background()

	req := &GetSiteRequest{}
	client := NewWorkflowClient(testBaseUrl, testApiKey)
	rsp, err := client.GetSite(ctx, req)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(rsp.MarshalIndent())
}
