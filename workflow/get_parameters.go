package workflow

import (
	"context"
	"encoding/json"
	"net/http"
)

type GetParametersRequest struct {
}

func (r *GetParametersRequest) String() string {
	if r == nil {
		return ""
	}
	bs, err := json.Marshal(r)
	if err != nil {
		return ""
	}
	return string(bs)
}

type GetParametersResponse struct {
	// Opening statement.
	OpeningStatement string `json:"opening_statement"`

	// List of suggested questions for the opening.
	SuggestedQuestions []string `json:"suggested_questions"`

	SuggestedQuestionsAfterAnswer struct {
		// Whether suggesting questions after an answer is enabled.
		Enabled bool `json:"enabled"`
	} `json:"suggested_questions_after_answer"`
	//
	SpeechToText struct {
		// Whether speech to text is enabled.
		Enabled bool `json:"enabled"`
	} `json:"speech_to_text"`
	//
	RetrieverResource struct {
		// Whether citation and attribution (retriever resource) is enabled.
		Enabled bool `json:"enabled"`
	} `json:"retriever_resource"`
	//
	AnnotationReply struct {
		// Whether annotation reply is enabled.
		Enabled bool `json:"enabled"`
	} `json:"annotation_reply"`
	// User input form configuration.
	// May Options:
	// Option 1 - Text input control.
	// Option 2 - Paragraph text input control.
	// Option 3 - Dropdown control.
	UserInputForm []map[string]struct {
		// Variable display label name.
		Label string `json:"label"`
		// Variable ID.
		Variable string `json:"variable"`
		// Whether it is required.
		Required bool `json:"required"`
		// Default value.
		Default string `json:"default"`
	} `json:"user_input_form"`
	// File upload configuration.
	FileUpload struct {
		// Image settings. Currently only supports image types: png, jpg, jpeg, webp, gif.
		Image struct {
			// Whether image upload is enabled.
			Enabled bool `json:"enabled"`
			// Image number limit, default is 3.
			NumberLimits int `json:"number_limits"`
			// Detail level for image processing (e.g., 'high').
			// From example, not in main description.
			Detail string `json:"detail"`
			// List of transfer methods, must choose at least one if enabled.
			TransferMethods []string `json:"transfer_methods"`
		} `json:"image"`
	} `json:"file_upload"`
	// Syetem parameters.
	SystemParameters struct {
		// Document upload size limit (MB).
		FileSizeLimit int `json:"file_size_limit"`
		// Image file upload size limit (MB).
		ImageFileSizeLimit int `json:"image_file_size_limit"`
		// Audio file upload size limit (MB).
		AudioFileSizeLimit int `json:"audio_file_size_limit"`
		// Video file upload size limit (MB).
		VideoFileSizeLimit int `json:"video_file_size_limit"`
	} `json:"system_parameters"`
}

func (r *GetParametersResponse) String() string {
	if r == nil {
		return ""
	}
	bs, err := json.Marshal(r)
	if err != nil {
		return ""
	}
	return string(bs)
}

func (r *GetParametersResponse) MarshalIndent() string {
	if r == nil {
		return ""
	}
	bs, err := json.MarshalIndent(r, "", "  ")
	if err != nil {
		return ""
	}
	return string(bs)
}

func (c *workflowClient) GetParameters(ctx context.Context, req *GetParametersRequest) (*GetParametersResponse, error) {
	httpReq, err := c.CreateBaseRequest(ctx, http.MethodGet, "/parameters", nil)
	if err != nil {
		return nil, err
	}

	var rsp GetParametersResponse
	err = c.SendJSONRequest(httpReq, &rsp)
	if err != nil {
		return nil, err
	}
	return &rsp, nil
}
