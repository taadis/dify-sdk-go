package workflow

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"time"
)

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
	// 时间区间参数
	CreatedAtAfter  *time.Time `json:"created_at__after,omitempty"`
	CreatedAtBefore *time.Time `json:"created_at__before,omitempty"`
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

func (c *workflowClient) GetWorkflowLogs(ctx context.Context, req *GetWorkflowLogsRequest) (*GetWorkflowLogsResponse, error) {
	r, err := c.CreateBaseRequest(ctx, http.MethodGet, "/workflows/logs", nil)
	if err != nil {
		return nil, err
	}

	query := r.URL.Query()
	if req.Keyword != "" {
		query.Set("keyword", req.Keyword)
	}
	if req.Status != "" {
		query.Set("status", string(req.Status))
	}
	if req.Page > 0 {
		query.Set("page", strconv.FormatInt(int64(req.Page), 10))
	}
	if req.Limit > 0 {
		query.Set("limit", strconv.FormatInt(int64(req.Limit), 10))
	}
	if req.CreatedByEndUserSessionId != "" {
		query.Set("created_by_end_user_session_id", req.CreatedByEndUserSessionId)
	}
	if req.CreatedByAccount != "" {
		query.Set("created_by_account", req.CreatedByAccount)
	}
	if req.CreatedAtAfter != nil {
		query.Set("created_at__after", req.CreatedAtAfter.Format(time.RFC3339))
	}
	if req.CreatedAtBefore != nil {
		query.Set("created_at__before", req.CreatedAtBefore.Format(time.RFC3339))
	}
	r.URL.RawQuery = query.Encode()

	var rsp GetWorkflowLogsResponse
	err = c.SendJSONRequest(r, &rsp)
	if err != nil {
		return nil, err
	}
	return &rsp, nil
}
