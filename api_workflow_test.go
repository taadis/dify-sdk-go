package dify

import (
	"context"
	"sync"
	"testing"
	"time"
)

func TestRunWorkflowStreaming(t *testing.T) {
	client := NewClient(host, apiSecretKey)

	workflowReq := WorkflowRequest{
		Inputs: map[string]interface{}{
			"image_url_new": map[string]string{
				"type":            "image",
				"transfer_method": "remote_url",
				"url":             "https://localhost/1-1.jpg",
			},
		},
		ResponseMode: "streaming",
		User:         "Zhaokm@AWS",
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	var (
		mu               sync.Mutex
		workflowStarted  bool
		nodeStarted      bool
		nodeFinished     bool
		workflowFinished bool
		ttsReceived      bool
	)

	// 创建一个实现 EventHandler 接口的处理器
	handler := &testEventHandler{
		t:  t,
		mu: &mu,
		onStreamingResponse: func(resp StreamingResponse) {
			mu.Lock()
			defer mu.Unlock()

			switch resp.Event {
			case EventWorkflowStarted:
				workflowStarted = true
			case EventNodeStarted:
				nodeStarted = true
			case EventNodeFinished:
				nodeFinished = true
				if resp.Data.ExecutionMetadata.TotalTokens > 0 {
					t.Logf("Node used %d tokens", resp.Data.ExecutionMetadata.TotalTokens)
				}
			case EventWorkflowFinished:
				workflowFinished = true
				if resp.Data.Status != "succeeded" {
					t.Errorf("Expected workflow status 'succeeded', got: %v", resp.Data.Status)
				}
			}
		},
		onTTSMessage: func(msg TTSMessage) {
			mu.Lock()
			defer mu.Unlock()

			ttsReceived = true
			if msg.Audio == "" {
				t.Error("Expected non-empty audio data in TTS message")
			}
		},
	}

	err := client.API().RunStreamWorkflowWithHandler(ctx, workflowReq, handler)

	if err != nil {
		t.Fatalf("RunStreamWorkflow encountered an error: %v", err)
	}

	mu.Lock()
	defer mu.Unlock()

	// 验证是否收到所有预期的事件
	if !workflowStarted {
		t.Error("Expected workflow_started event, but didn't receive it")
	}
	if !nodeStarted {
		t.Error("Expected node_started event, but didn't receive it")
	}
	if !nodeFinished {
		t.Error("Expected node_finished event, but didn't receive it")
	}
	if !workflowFinished {
		t.Error("Expected workflow_finished event, but didn't receive it")
	}
	if !ttsReceived {
		t.Error("Expected TTS message, but didn't receive it")
	}

	t.Log("Streaming workflow test completed successfully")
}

func TestGetWorkflowRunDetail(t *testing.T) {
	ctx := context.Background()
	req := GetWorkflowRunDetailRequest{
		WorkflowRunId: "test-workflow-run-id",
	}

	client := NewClient(host, apiSecretKey)
	rsp, err := client.API().GetWorkflowRunDetail(ctx, &req)
	if err != nil {
		t.Fatalf("RunWorkflow encountered an error: %v", err)
	}
	t.Log(rsp.String())

}
