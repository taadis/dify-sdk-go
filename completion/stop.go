package completion

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type StopRequest struct {
	// Task ID, can be obtained from the streaming chunk return of a /completion-messages request.
	TaskId string `json:"-"`
	// User identifier, used to define the identity of the end-user, must be consistent with the user passed in the send message interface.
	User string `json:"user"`
}

type StopResponse struct {
	// Example: "success"
	Result string `json:"result"`
}

func (r *StopResponse) String() string {
	if r == nil {
		return ""
	}
	bs, err := json.Marshal(r)
	if err != nil {
		return ""
	}
	return string(bs)
}

func (r *StopResponse) MarshalIndent() string {
	if r == nil {
		return ""
	}
	bs, err := json.MarshalIndent(r, "", "  ")
	if err != nil {
		return ""
	}
	return string(bs)
}

func (c *completionClient) Stop(ctx context.Context, req *StopRequest) (*StopResponse, error) {
	if req.TaskId == "" {
		return nil, fmt.Errorf("missing required task_id")
	}
	// %s={task_id}
	apiUrl := fmt.Sprintf("/completions-messages/%s/stop", req.TaskId)
	r, err := c.CreateBaseRequest(ctx, http.MethodPost, apiUrl, req)
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

	var ret StopResponse
	if err := json.NewDecoder(rsp.Body).Decode(&ret); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}
	return &ret, nil
}
