package dify

import (
	"context"
	"encoding/json"
	"testing"
)

func TestApplication(t *testing.T) {
	ctx := context.Background()

	// messages
	req := &GetApplicationBasicInformationRequest{}

	client := NewClient(host, apiSecretKey)
	rsp, err := client.API().GetApplicationBasicInformation(ctx, req)
	if err != nil {
		t.Fatal(err)
	}

	j, _ := json.Marshal(rsp)
	t.Log(string(j))
}
