package workflow

import (
	"context"
	"testing"
)

func TestUploadFile(t *testing.T) {
	ctx := context.Background()

	filePath := "your-file-path"
	user := "your-user"
	if filePath == "your-file-path" || user == "your-user" {
		t.Skip("Set a valid file path and user to run this test.")
	}

	req := &UploadFileRequest{FilePath: filePath, User: user}
	client := NewWorkflowClient(testBaseUrl, testApiKey)
	rsp, err := client.UploadFile(ctx, req)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(rsp.String())
}
