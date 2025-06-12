package dify

import (
	"context"
	"net/http"
)

type GetApplicationBasicInformationRequest struct {
}

type GetApplicationBasicInformationResponse struct {
	Name string `json:"name"`

	Description string `json:"description"`

	Tags []string `json:"tags"`
}

// used to get basic information about this application.
func (api *API) GetApplicationBasicInformation(ctx context.Context, req *GetApplicationBasicInformationRequest) (*GetApplicationBasicInformationResponse, error) {
	httpReq, err := api.createBaseRequest(ctx, http.MethodGet, "/v1/info", nil)
	if err != nil {
		return nil, err
	}

	var rsp GetApplicationBasicInformationResponse
	err = api.c.sendJSONRequest(httpReq, &rsp)
	if err != nil {
		return nil, err
	}
	return &rsp, nil
}
