package logger

import (
	"bytes"
	"context"
	"encoding/json"
	"testing"
)

func TestInfoContainsRequestIDAndLevel(t *testing.T) {
	buf := &bytes.Buffer{}
	log := New(WithOutput(buf), WithStage("staging"))

	err := log.WithRequestID("req-123").Info("hello", map[string]any{"foo": "bar"})
	if err != nil {
		t.Fatalf("log Info returned error: %v", err)
	}

	var got map[string]any
	if err := json.NewDecoder(buf).Decode(&got); err != nil {
		t.Fatalf("failed to decode log: %v", err)
	}

	if got["request_id"] != "req-123" {
		t.Fatalf("request_id mismatch: got %v", got["request_id"])
	}
	if got["level"] != string(LevelInfo) {
		t.Fatalf("level mismatch: got %v", got["level"])
	}
	if got["stage"] != "staging" {
		t.Fatalf("stage mismatch: got %v", got["stage"])
	}
	if got["foo"] != "bar" {
		t.Fatalf("field merge failed: got %v", got["foo"])
	}
}

func TestWithContextUsesInjectedRequestID(t *testing.T) {
	buf := &bytes.Buffer{}
	log := New(WithOutput(buf))
	ctx := ContextWithRequestID(context.Background(), "ctx-req")

	if err := log.WithContext(ctx).Warn("warn", nil); err != nil {
		t.Fatalf("log Warn returned error: %v", err)
	}

	var got map[string]any
	if err := json.NewDecoder(buf).Decode(&got); err != nil {
		t.Fatalf("failed to decode log: %v", err)
	}

	if got["request_id"] != "ctx-req" {
		t.Fatalf("request_id mismatch: got %v", got["request_id"])
	}
	if got["level"] != string(LevelWarn) {
		t.Fatalf("level mismatch: got %v", got["level"])
	}
}
