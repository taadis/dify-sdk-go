package dify

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"strings"
	"sync"
	"testing"
)

var (
	host         = os.Getenv("DIFY_API_HOST") // your dify api host
	apiSecretKey = os.Getenv("DIFY_API_KEY")  // your dify api secret key e.g.app-xxx
)

func TestMain(m *testing.M) {
	host = DifyCloud
	os.Exit(m.Run())
}

func TestApi3(t *testing.T) {
	var c = &ClientConfig{
		Host:         host,
		ApiSecretKey: apiSecretKey,
	}
	var client = NewClientWithConfig(c)

	ctx := context.Background()

	var (
		ch  = make(chan ChatMessageStreamChannelResponse)
		err error
	)

	ch, err = client.Api().ChatMessagesStream(ctx, &ChatMessageRequest{
		Query: "你是谁?",
		User:  "这里换成你创建的",
	})
	if err != nil {
		t.Fatal(err.Error())
	}

	var (
		strBuilder strings.Builder
		cId        string
	)
	for {
		select {
		case <-ctx.Done():
			t.Log("ctx.Done", strBuilder.String())
			return
		case r, isOpen := <-ch:
			if !isOpen {
				goto done
			}
			strBuilder.WriteString(r.Answer)
			cId = r.ConversationID
			log.Println("Answer2", r.Answer, r.ConversationID, cId, r.ID, r.TaskID)
		}
	}

done:
	t.Log(strBuilder.String())
	t.Log(cId)
}

func TestMessages(t *testing.T) {
	var cId = "ec373942-2d17-4f11-89bb-f9bbf863ebcc"
	var err error
	ctx := context.Background()

	// messages
	var messageReq = &MessagesRequest{
		ConversationID: cId,
		User:           "test-user",
	}

	var client = NewClient(host, apiSecretKey)

	var msg *MessagesResponse
	if msg, err = client.Api().Messages(ctx, messageReq); err != nil {
		t.Fatal(err.Error())
		return
	}
	j, _ := json.Marshal(msg)
	t.Log(string(j))
}

func TestMessagesFeedbacks(t *testing.T) {
	var client = NewClient(host, apiSecretKey)
	var err error
	ctx := context.Background()

	var id = "72d3dc0f-a6d5-4b5e-8510-bec0611a6048"

	var res *MessagesFeedbacksResponse
	if res, err = client.Api().MessagesFeedbacks(ctx, &MessagesFeedbacksRequest{
		MessageID: id,
		Rating:    FeedbackLike,
		User:      "test-user",
	}); err != nil {
		t.Fatal(err.Error())
	}

	j, _ := json.Marshal(res)

	log.Println(string(j))
}

func TestConversations(t *testing.T) {
	var client = NewClient(host, apiSecretKey)
	var err error
	ctx := context.Background()

	var res *ConversationsResponse
	if res, err = client.Api().Conversations(ctx, &ConversationsRequest{
		User: "test-user",
	}); err != nil {
		t.Fatal(err.Error())
	}

	j, _ := json.Marshal(res)

	log.Println(string(j))
}

func TestConversationsRename(t *testing.T) {
	var client = NewClient(host, apiSecretKey)
	var err error
	ctx := context.Background()

	var res *ConversationsRenamingResponse
	if res, err = client.Api().ConversationsRenaming(ctx, &ConversationsRenamingRequest{
		ConversationID: "ab383831-1e38-3d00-76cd-d6ppe752cedd",
		Name:           "rename!!!",
		User:           "test-user",
	}); err != nil {
		t.Fatal(err.Error())
	}

	j, _ := json.Marshal(res)

	log.Println(string(j))
}

func TestParameters(t *testing.T) {
	var client = NewClient(host, apiSecretKey)
	var err error
	ctx := context.Background()

	var res *ParametersResponse
	if res, err = client.Api().Parameters(ctx, &ParametersRequest{
		User: "test-user",
	}); err != nil {
		t.Fatal(err.Error())
	}

	j, _ := json.Marshal(res)

	log.Println(string(j))
}

// testEventHandler 实现 EventHandler 接口
type testEventHandler struct {
	t                   *testing.T
	mu                  *sync.Mutex
	onStreamingResponse func(StreamingResponse)
	onTTSMessage        func(TTSMessage)
}

func (h *testEventHandler) HandleStreamingResponse(resp StreamingResponse) {
	if h.onStreamingResponse != nil {
		h.onStreamingResponse(resp)
	}
}

func (h *testEventHandler) HandleTTSMessage(msg TTSMessage) {
	if h.onTTSMessage != nil {
		h.onTTSMessage(msg)
	}
}
