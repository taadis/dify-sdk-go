package workflow

import (
	"context"
	"encoding/json"
	"net/http"
)

type GetInfoRequest struct {
}

func (r *GetInfoRequest) String() string {
	if r == nil {
		return ""
	}
	bs, err := json.Marshal(r)
	if err != nil {
		return ""
	}
	return string(bs)
}

type GetInfoResponse struct {
	// application name.
	Name string `json:"name"`

	// application description.
	Description string `json:"description"`

	// application tags.
	Tags []string `json:"tags"`
}

func (r *GetInfoResponse) String() string {
	if r == nil {
		return ""
	}
	bs, err := json.Marshal(r)
	if err != nil {
		return ""
	}
	return string(bs)
}

func (c *workflowClient) GetInfo(ctx context.Context, req *GetInfoRequest) (*GetInfoResponse, error) {
	httpReq, err := c.CreateBaseRequest(ctx, http.MethodGet, "/info", nil)
	if err != nil {
		return nil, err
	}

	var rsp GetInfoResponse
	err = c.SendJSONRequest(httpReq, &rsp)
	if err != nil {
		return nil, err
	}
	return &rsp, nil
}
