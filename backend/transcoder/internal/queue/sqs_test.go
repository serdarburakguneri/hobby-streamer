package queue

import (
	"context"
	"encoding/json"
	"errors"
	"testing"
)

type mockSQSClient struct {
	ReceiveMessageFunc func(ctx context.Context) ([]QueueMessage, error)
}

type testHandler struct {
	called bool
	fail   bool
}

func (h *testHandler) Handle(msg QueueMessage) error {
	h.called = true
	if h.fail {
		return errors.New("fail")
	}
	return nil
}

func TestSQSConsumer_HandlesMessages(t *testing.T) {
	h := &testHandler{}
	msg := QueueMessage{Type: "foo", Payload: json.RawMessage(`{}`)}
	err := h.Handle(msg)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if !h.called {
		t.Error("expected handler to be called")
	}
}

func TestSQSConsumer_HandlerError(t *testing.T) {
	h := &testHandler{fail: true}
	msg := QueueMessage{Type: "foo", Payload: json.RawMessage(`{}`)}
	err := h.Handle(msg)
	if err == nil {
		t.Error("expected error, got nil")
	}
}

func TestSQSConsumer_UnmarshalError(t *testing.T) {
	var msg QueueMessage
	bad := []byte(`notjson`)
	err := json.Unmarshal(bad, &msg)
	if err == nil {
		t.Error("expected unmarshal error, got nil")
	}
}
