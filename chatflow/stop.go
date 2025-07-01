package chatflow

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type StopRequest struct {
	// Task ID from the streaming chunk.
	TaskId string `json:"-"`
	// User identifier, consistent with the send message call.
	// Note: The Service API does not share conversations created by the WebApp.
	// Conversations created through the API are isolated from those created in the WebApp interface.
	User string `json:"user"`
}

// Operation successful.
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

func (c *chatflowClient) Stop(ctx context.Context, req *StopRequest) (*StopResponse, error) {
	if req.TaskId == "" {
		return nil, fmt.Errorf("missing required task_id")
	}
	// %s={task_id}
	apiUrl := fmt.Sprintf("/chat-messages/%s/stop", req.TaskId)
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
