package workflow

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// FileInput 结构体
type FileInput struct {
	Type           string `json:"type"`                     // 目前仅支持 "image"
	TransferMethod string `json:"transfer_method"`          // "remote_url" 或 "local_file"
	URL            string `json:"url,omitempty"`            // 当 transfer_method 为 remote_url 时使用
	UploadFileID   string `json:"upload_file_id,omitempty"` // 当 transfer_method 为 local_file 时使用
}

type RunRequest struct {
	Inputs       map[string]interface{} `json:"inputs"`
	ResponseMode string                 `json:"response_mode"`
	User         string                 `json:"user"`
	Files        []FileInput            `json:"files,omitempty"`
}

type RunResponse struct {
	WorkflowRunId string `json:"workflow_run_id"`
	TaskId        string `json:"task_id"`
	Data          struct {
		Id          string                 `json:"id"`
		WorkflowId  string                 `json:"workflow_id"`
		Status      string                 `json:"status"`
		Outputs     map[string]interface{} `json:"outputs"`
		Error       *string                `json:"error,omitempty"`
		ElapsedTime float64                `json:"elapsed_time"`
		TotalTokens int                    `json:"total_tokens"`
		TotalSteps  int                    `json:"total_steps"`
		CreatedAt   int64                  `json:"created_at"`
		FinishedAt  int64                  `json:"finished_at"`
	} `json:"data"`
}

func (r *RunResponse) String() string {
	if r == nil {
		return ""
	}
	bs, err := json.Marshal(r)
	if err != nil {
		return ""
	}
	return string(bs)
}

func (c *workflowClient) Run(ctx context.Context, req *RunRequest) (*RunResponse, error) {
	r, err := c.CreateBaseRequest(ctx, http.MethodPost, "/workflows/run", req)
	if err != nil {
		return nil, fmt.Errorf("failed to create base request: %w", err)
	}

	rsp, err := c.SendRequest(r)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer rsp.Body.Close()

	if rsp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status %s: %s", rsp.Status, c.ReadResponseBody(rsp.Body))
	}

	var workflowResp RunResponse
	if err := json.NewDecoder(rsp.Body).Decode(&workflowResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &workflowResp, nil
}
