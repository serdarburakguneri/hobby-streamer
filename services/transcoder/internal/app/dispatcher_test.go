package app

import (
	"context"
	"encoding/json"
	"errors"
	"testing"

	"github.com/serdarburakguneri/hobby-streamer/services/transcoder/internal/job"
	"github.com/serdarburakguneri/hobby-streamer/services/transcoder/internal/queue"
)

type mockRunner struct {
	called bool
	fail   bool
}

func (m *mockRunner) Run(ctx context.Context, payload json.RawMessage) error {
	m.called = true
	if m.fail {
		return errors.New("fail")
	}
	return nil
}

func TestDispatcher_HandleMessage_RoutesCorrectly(t *testing.T) {
	r := job.NewRegistry()
	mock := &mockRunner{}
	r.Register("foo", mock)
	d := NewDispatcher(r)

	msg := queue.QueueMessage{Type: "foo", Payload: json.RawMessage(`{}`)}
	err := d.HandleMessage(msg)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if !mock.called {
		t.Error("expected runner to be called")
	}
}

func TestDispatcher_HandleMessage_UnknownType(t *testing.T) {
	r := job.NewRegistry()
	d := NewDispatcher(r)
	msg := queue.QueueMessage{Type: "bar", Payload: json.RawMessage(`{}`)}
	err := d.HandleMessage(msg)
	if err == nil {
		t.Error("expected error for unknown type, got nil")
	}
}
