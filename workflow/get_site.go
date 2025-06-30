package workflow

import (
	"context"
	"encoding/json"
	"net/http"
)

type GetSiteRequest struct {
}

func (r *GetSiteRequest) String() string {
	if r == nil {
		return ""
	}
	bs, err := json.Marshal(r)
	if err != nil {
		return ""
	}
	return string(bs)
}

type GetSiteResponse struct {
	// WebApp Name
	Title string `json:"title"`
	// Icon type.
	// Available options: emoji, image
	IconType string `json:"icon_type"`
	// Icon. If it's emoji type, it's an emoji symbol; if it's image type, it's an image URL.
	Icon string `json:"icon"`
	// Background color in hex format (e.g., #RRGGBB).
	IconBackground string `json:"icon_background"`
	// Icon URL (likely refers to image type if icon field is just a name/id).
	IconUrl string `json:"icon_url"`
	// Description
	Description string `json:"description"`
	// Copyright information.
	Copyright string `json:"copyright"`
	// Privacy policy link.
	PrivacyPolicy string `json:"privacy_policy"`
	// Custom disclaimer.
	CustomDisclaimer string `json:"custom_disclaimer"`
	// Default language (e.g., en-US).
	DefaultLanguage string `json:"default_language"`
	// Whether to show workflow details.
	ShowWorkflowSteps bool `json:"show_workflow_steps"`
	// Whether to replace 🤖 in chat with the WebApp icon.
	UseIconAsAnswerIcon bool `json:"use_icon_as_answer_icon"`
}

func (r *GetSiteResponse) String() string {
	if r == nil {
		return ""
	}
	bs, err := json.Marshal(r)
	if err != nil {
		return ""
	}
	return string(bs)
}

func (r *GetSiteResponse) MarshalIndent() string {
	if r == nil {
		return ""
	}
	bs, err := json.MarshalIndent(r, "", "  ")
	if err != nil {
		return ""
	}
	return string(bs)
}

func (c *workflowClient) GetSite(ctx context.Context, req *GetSiteRequest) (*GetSiteResponse, error) {
	r, err := c.CreateBaseRequest(ctx, http.MethodGet, "/site", nil)
	if err != nil {
		return nil, err
	}

	var rsp GetSiteResponse
	err = c.SendJSONRequest(r, &rsp)
	if err != nil {
		return nil, err
	}
	return &rsp, nil
}
