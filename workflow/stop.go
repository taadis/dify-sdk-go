package workflow

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// StopTaskRequest is the request struct for stopping a workflow task generation.
type StopTaskRequest struct {
	TaskId string `json:"task_id"`
	User   string `json:"user"`
}

func (r *StopTaskRequest) String() string {
	if r == nil {
		return ""
	}
	bs, err := json.Marshal(r)
	if err != nil {
		return ""
	}
	return string(bs)
}

// StopTaskResponse is the response struct for stopping a workflow task generation.
type StopTaskResponse struct {
	Result string `json:"result"`
}

func (r *StopTaskResponse) String() string {
	if r == nil {
		return ""
	}
	bs, err := json.Marshal(r)
	if err != nil {
		return ""
	}
	return string(bs)
}

// StopTask stops a workflow task generation by task_id and user.
func (c *workflowClient) StopTask(ctx context.Context, req *StopTaskRequest) (*StopTaskResponse, error) {
	if req == nil || req.TaskId == "" || req.User == "" {
		return nil, fmt.Errorf("task_id and user are required")
	}
	url := fmt.Sprintf("/workflows/tasks/%s/stop", req.TaskId)
	httpReq, err := c.CreateBaseRequest(ctx, http.MethodPost, url, map[string]string{"user": req.User})
	if err != nil {
		return nil, err
	}

	var rsp StopTaskResponse
	err = c.SendJSONRequest(httpReq, &rsp)
	if err != nil {
		return nil, err
	}
	return &rsp, nil
}
