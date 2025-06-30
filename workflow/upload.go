package workflow

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
)

type UploadFileRequest struct {
	FilePath string // Local file path
	User     string // User identifier
}

func (r *UploadFileRequest) String() string {
	if r == nil {
		return ""
	}
	bs, err := json.Marshal(r)
	if err != nil {
		return ""
	}
	return string(bs)
}

type UploadFileResponse struct {
	Id        string `json:"id"`
	Name      string `json:"name"`
	Size      int64  `json:"size"`
	Extension string `json:"extension"`
	MimeType  string `json:"mime_type"`
	CreatedBy string `json:"created_by"`
	CreatedAt int64  `json:"created_at"`
}

func (r *UploadFileResponse) String() string {
	if r == nil {
		return ""
	}
	bs, err := json.Marshal(r)
	if err != nil {
		return ""
	}
	return string(bs)
}

// UploadFileErrorResponse represents an error response from the upload API.
type UploadFileErrorResponse struct {
	Status  int    `json:"status"`
	Code    string `json:"code"`
	Message string `json:"message"`
}

// UploadFile uploads a file for workflow use
func (c *workflowClient) UploadFile(ctx context.Context, req *UploadFileRequest) (*UploadFileResponse, error) {
	if req == nil || req.FilePath == "" || req.User == "" {
		return nil, fmt.Errorf("file path and user are required")
	}

	file, err := os.Open(req.FilePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	pr, pw := io.Pipe()
	writer := multipart.NewWriter(pw)
	go func() {
		defer pw.Close()
		// file field
		part, err := writer.CreateFormFile("file", filepath.Base(req.FilePath))
		if err != nil {
			pw.CloseWithError(err)
			return
		}
		if _, err := io.Copy(part, file); err != nil {
			pw.CloseWithError(err)
			return
		}
		// user field
		if err := writer.WriteField("user", req.User); err != nil {
			pw.CloseWithError(err)
			return
		}
		writer.Close()
	}()

	url := "/files/upload"
	httpReq, err := c.CreateBaseRequest(ctx, http.MethodPost, url, nil)
	if err != nil {
		return nil, err
	}
	httpReq.Body = io.NopCloser(pr)
	httpReq.Header.Set("Content-Type", writer.FormDataContentType())

	resp, err := c.SendRequest(httpReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		var errBody UploadFileErrorResponse
		_ = json.NewDecoder(resp.Body).Decode(&errBody)
		return nil, fmt.Errorf("HTTP response error: [%v]%v", errBody.Code, errBody.Message)
	}

	var rsp UploadFileResponse
	if err := json.NewDecoder(resp.Body).Decode(&rsp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}
	return &rsp, nil
}
