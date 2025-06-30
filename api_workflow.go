package dify

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// 事件类型常量
const (
	EventWorkflowStarted  = "workflow_started"
	EventNodeStarted      = "node_started"
	EventNodeFinished     = "node_finished"
	EventWorkflowFinished = "workflow_finished"
	EventTTSMessage       = "tts_message"
	EventTTSMessageEnd    = "tts_message_end"
)

// FileInput 结构体
type FileInput struct {
	Type           string `json:"type"`                     // 目前仅支持 "image"
	TransferMethod string `json:"transfer_method"`          // "remote_url" 或 "local_file"
	URL            string `json:"url,omitempty"`            // 当 transfer_method 为 remote_url 时使用
	UploadFileID   string `json:"upload_file_id,omitempty"` // 当 transfer_method 为 local_file 时使用
}

// WorkflowRequest 结构体
type WorkflowRequest struct {
	Inputs       map[string]interface{} `json:"inputs"`
	ResponseMode string                 `json:"response_mode"`
	User         string                 `json:"user"`
	Files        []FileInput            `json:"files,omitempty"`
}

// StreamingResponse 结构体
type StreamingResponse struct {
	Event          string `json:"event"`
	TaskID         string `json:"task_id"`
	WorkflowRunID  string `json:"workflow_run_id"`
	SequenceNumber int    `json:"sequence_number"`
	Data           struct {
		ID                string                 `json:"id"`
		WorkflowID        string                 `json:"workflow_id,omitempty"`
		NodeID            string                 `json:"node_id,omitempty"`
		NodeType          string                 `json:"node_type,omitempty"`
		Title             string                 `json:"title,omitempty"`
		Index             int                    `json:"index"`
		Predecessor       string                 `json:"predecessor_node_id,omitempty"`
		Inputs            []interface{}          `json:"inputs,omitempty"`
		Outputs           map[string]interface{} `json:"outputs,omitempty"`
		Status            string                 `json:"status,omitempty"`
		Error             string                 `json:"error,omitempty"`
		ElapsedTime       float64                `json:"elapsed_time,omitempty"`
		ExecutionMetadata struct {
			TotalTokens int     `json:"total_tokens,omitempty"`
			TotalPrice  float64 `json:"total_price,omitempty"`
			Currency    string  `json:"currency,omitempty"`
		} `json:"execution_metadata,omitempty"`
		CreatedAt  int64 `json:"created_at"`
		FinishedAt int64 `json:"finished_at,omitempty"`
	} `json:"data"`
}

// TTSMessage 结构体
type TTSMessage struct {
	Event     string `json:"event"` // "tts_message" 或 "tts_message_end"
	TaskID    string `json:"task_id"`
	MessageID string `json:"message_id"`
	Audio     string `json:"audio"` // Base64 编码的音频数据
	CreatedAt int64  `json:"created_at"`
}

// EventHandler 接口
type EventHandler interface {
	HandleStreamingResponse(StreamingResponse)
	HandleTTSMessage(TTSMessage)
}

// DefaultEventHandler 结构体
type DefaultEventHandler struct {
	StreamHandler func(StreamingResponse)
}

func (h *DefaultEventHandler) HandleStreamingResponse(resp StreamingResponse) {
	if h.StreamHandler != nil {
		h.StreamHandler(resp)
	}
}

func (h *DefaultEventHandler) HandleTTSMessage(msg TTSMessage) {
	// 默认实现为空，如果用户不关心 TTS 消息可以忽略
}

// RunStreamWorkflow 方法
func (api *API) RunStreamWorkflow(ctx context.Context, request WorkflowRequest, handler func(StreamingResponse)) error {
	return api.RunStreamWorkflowWithHandler(ctx, request, &DefaultEventHandler{StreamHandler: handler})
}

// RunStreamWorkflowWithHandler 方法
func (api *API) RunStreamWorkflowWithHandler(ctx context.Context, request WorkflowRequest, handler EventHandler) error {
	req, err := api.createBaseRequest(ctx, http.MethodPost, "/v1/workflows/run", request)
	if err != nil {
		return err
	}

	resp, err := api.c.sendRequest(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("API request failed with status %s: %s", resp.Status, readResponseBody(resp.Body))
	}

	reader := bufio.NewReader(resp.Body)
	for {
		line, err := reader.ReadBytes('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			return fmt.Errorf("error reading streaming response: %w", err)
		}

		if len(line) > 6 && string(line[:6]) == "data: " {
			var event struct {
				Event string `json:"event"`
			}
			if err := json.Unmarshal(line[6:], &event); err != nil {
				fmt.Println("Error decoding event type:", err)
				continue
			}

			switch event.Event {
			case EventTTSMessage, EventTTSMessageEnd:
				var ttsMsg TTSMessage
				if err := json.Unmarshal(line[6:], &ttsMsg); err != nil {
					fmt.Println("Error decoding TTS message:", err)
					continue
				}
				handler.HandleTTSMessage(ttsMsg)
			default:
				var streamResp StreamingResponse
				if err := json.Unmarshal(line[6:], &streamResp); err != nil {
					fmt.Println("Error decoding streaming response:", err)
					continue
				}
				handler.HandleStreamingResponse(streamResp)
			}
		}
	}

	return nil
}

// readResponseBody 辅助函数
func readResponseBody(body io.Reader) string {
	bodyBytes, err := io.ReadAll(body)
	if err != nil {
		return fmt.Sprintf("failed to read response body: %v", err)
	}
	return string(bodyBytes)
}

type GetWorkflowRunDetailRequest struct {
	WorkflowRunId string `json:"workflow_run_id"`
}

type GetWorkflowRunDetailResponse struct {
	// e.g.3c90c3cc-0d44-4b50-8888-8dd25736052a
	Id string `json:"id"`
	// e.g.3c90c3cc-0d44-4b50-8888-8dd25736052a
	WorkflowId string `json:"workflow_id"`
	// e.g.running,succeeded,failed,stopped
	Status string `json:"status"`
	// JSON string of input content.
	Inputs map[string]interface{} `json:"inputs"`
	// JSON object of output content.
	Outputs map[string]interface{} `json:"outputs,omitempty"`
	//
	Error *string `json:"error,omitempty"`
	// e.g.123
	TotalSteps int64 `json:"total_steps"`
	// e.g.123
	TotalTokens int64 `json:"total_tokens"`
	// e.g.123
	CreatedAt int64 `json:"created_at"`
	// e.g.123
	FinishedAt int64 `json:"finished_at"`
	// e.g.123
	ElapsedTime int64 `json:"elapsed_time"`
}

func (r *GetWorkflowRunDetailResponse) String() string {
	if r == nil {
		return ""
	}
	bs, _ := json.Marshal(r)
	return string(bs)
}

func (api *API) GetWorkflowRunDetail(ctx context.Context, req *GetWorkflowRunDetailRequest) (*GetWorkflowRunDetailResponse, error) {
	r, err := api.createBaseRequest(ctx, http.MethodGet, fmt.Sprintf("/v1/workflows/run/%s", req.WorkflowRunId), nil)
	if err != nil {
		return nil, err
	}

	var rsp GetWorkflowRunDetailResponse
	err = api.c.sendJSONRequest(r, &rsp)
	if err != nil {
		return nil, err
	}
	return &rsp, nil
}

type WorkflowStatus string

const (
	WorkflowStatusSucceeded WorkflowStatus = "succeeded"
	WorkflowStatusFailed    WorkflowStatus = "failed"
	WorkflowStatusStopped   WorkflowStatus = "stopped"
	WorkflowStatusRunning   WorkflowStatus = "running"
)

type GetWorkflowLogsRequest struct {
	// 关键字
	Keyword string `json:"keyword"`
	// 执行状态
	Status WorkflowStatus `json:"status"`
	// 当前页码,默认1
	Page int `json:"page"`
	// 每页条数,默认20
	Limit int `json:"limit"`
	// 由哪个endUser创建,例如:abc-123
	CreatedByEndUserSessionId string `json:"created_by_end_user_session_id"`
	// 由哪个邮箱账户创建,例如:lizb@test.com
	CreatedByAccount string `json:"created_by_account"`
}

func (r *GetWorkflowLogsRequest) String() string {
	if r == nil {
		return ""
	}
	bs, _ := json.Marshal(r)
	return string(bs)
}

type GetWorkflowLogsResponse struct {
	// 当前页码
	Page int `json:"page"`
	// 每页条数
	Limit int `json:"limit"`
	// 总条数
	Total int `json:"total"`
	// 是否还有更多数据
	HasMore bool `json:"has_more"`
	// 当前页码的数据
	Data []*WorkflowsLogItem `json:"data"`
}

type WorkflowsLogItem struct {
	// 标识
	Id string `json:"id"`
	// Workflow执行日志
	WorkflowRun struct {
		// 标识
		Id string `json:"id"`
		// 版本
		Version string `json:"version"`
		// 执行状态:running/succeeded/failed/stopped
		Status WorkflowStatus `json:"status"`
		// 错误(可选)
		Error string `json:"error,omitempty"`
		// 耗时,单位秒
		ElapsedTime float64 `json:"elapsed_time"`
		// 消耗的token数量
		TotalTokens int64 `json:"total_tokens"`
		// 执行的步骤长度
		TotalSteps int64 `json:"total_steps"`
		// 开始时间
		CreatedAt int64 `json:"created_at"`
		// 结束时间
		FinishedAt int64 `json:"finished_at"`
	} `json:"workflow_run"`

	// 来源
	CreatedFrom string `json:"created_from"`
	// 角色
	CreatedByRole string `json:"created_by_role"`
	// 账号(可选)
	CreatedByAccount string `json:"created_by_account,omitempty"`
	// 用户
	CreatedByEndUser struct {
		// 标识
		Id string `json:"id"`
		// 类型
		Type string `json:"type"`
		// 是否匿名
		IsAnonymous bool `json:"is_anonymous"`
		// 会话标识
		SessionId string `json:"session_id"`
	} `json:"created_by_end_user"`
	// 创建时间
	CreatedAt int64 `json:"created_at"`
}

func (r *GetWorkflowLogsResponse) String() string {
	if r == nil {
		return ""
	}
	bs, _ := json.Marshal(r)
	return string(bs)
}
