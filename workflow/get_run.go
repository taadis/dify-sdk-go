package workflow

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// GetRunRequest is the request struct for retrieving workflow run details.
type GetRunRequest struct {
	WorkflowRunId string `json:"workflow_run_id"`
}

func (r *GetRunRequest) String() string {
	if r == nil {
		return ""
	}
	bs, err := json.Marshal(r)
	if err != nil {
		return ""
	}
	return string(bs)
}

// GetRunResponse is the response struct for workflow run details.
type GetRunResponse struct {
	Id          string                 `json:"id"`
	WorkflowId  string                 `json:"workflow_id"`
	Status      string                 `json:"status"`
	Inputs      string                 `json:"inputs"`
	Outputs     map[string]interface{} `json:"outputs"`
	Error       *string                `json:"error"`
	TotalSteps  int                    `json:"total_steps"`
	TotalTokens int                    `json:"total_tokens"`
	CreatedAt   int64                  `json:"created_at"`
	FinishedAt  *int64                 `json:"finished_at"`
	ElapsedTime *float64               `json:"elapsed_time"`
}

func (r *GetRunResponse) String() string {
	if r == nil {
		return ""
	}
	bs, err := json.Marshal(r)
	if err != nil {
		return ""
	}
	return string(bs)
}

// GetRun retrieves the current execution results of a workflow task by workflow_run_id.
func (c *workflowClient) GetRun(ctx context.Context, req *GetRunRequest) (*GetRunResponse, error) {
	if req == nil || req.WorkflowRunId == "" {
		return nil, fmt.Errorf("workflow_run_id is required")
	}
	url := fmt.Sprintf("/workflows/run/%s", req.WorkflowRunId)
	httpReq, err := c.CreateBaseRequest(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	var rsp GetRunResponse
	err = c.SendJSONRequest(httpReq, &rsp)
	if err != nil {
		return nil, err
	}
	return &rsp, nil
}
