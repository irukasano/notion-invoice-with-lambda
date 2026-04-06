package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"testing"

	"github.com/irukasano/notion-invoice-with-lambda/internal/logger"
)

func TestHandleLogsStartAndFinish(t *testing.T) {
	buf := &bytes.Buffer{}
	log := logger.New(logger.WithOutput(buf), logger.WithStage("staging"))
	h := handler{
		app: stubApp{logger: log},
	}
	ctx := logger.ContextWithRequestID(context.Background(), "req-mail")

	got, err := h.handle(ctx)
	if err != nil {
		t.Fatalf("handle returned error: %v", err)
	}
	if got != nil {
		t.Fatalf("handle response mismatch: got %v", got)
	}

	logs := decodeLogs(t, buf)
	if len(logs) != 2 {
		t.Fatalf("log count mismatch: got %d", len(logs))
	}
	assertLog(t, logs[0], "INFO", "invoice-mail started", "req-mail")
	assertLog(t, logs[1], "INFO", "invoice-mail finished", "req-mail")
}

func TestHandleLogsFailure(t *testing.T) {
	buf := &bytes.Buffer{}
	log := logger.New(logger.WithOutput(buf), logger.WithStage("staging"))
	h := handler{
		app: stubApp{logger: log, runErr: errors.New("boom")},
	}
	ctx := logger.ContextWithRequestID(context.Background(), "req-mail")

	got, err := h.handle(ctx)
	if err == nil {
		t.Fatal("handle error = nil")
	}
	if got != nil {
		t.Fatalf("handle response mismatch: got %v", got)
	}

	logs := decodeLogs(t, buf)
	if len(logs) != 2 {
		t.Fatalf("log count mismatch: got %d", len(logs))
	}
	assertLog(t, logs[0], "INFO", "invoice-mail started", "req-mail")
	assertLog(t, logs[1], "ERROR", "invoice-mail failed", "req-mail")
	if logs[1]["error"] != "boom" {
		t.Fatalf("error field mismatch: got %v", logs[1]["error"])
	}
}

type stubApp struct {
	logger *logger.Logger
	runErr error
}

func (s stubApp) Run(ctx context.Context) error {
	_ = ctx
	return s.runErr
}

func (s stubApp) Logger() *logger.Logger {
	return s.logger
}

func decodeLogs(t *testing.T, buf *bytes.Buffer) []map[string]any {
	t.Helper()

	decoder := json.NewDecoder(buf)
	var logs []map[string]any
	for decoder.More() {
		var entry map[string]any
		if err := decoder.Decode(&entry); err != nil {
			t.Fatalf("failed to decode log: %v", err)
		}
		logs = append(logs, entry)
	}
	return logs
}

func assertLog(t *testing.T, entry map[string]any, wantLevel, wantMessage, wantRequestID string) {
	t.Helper()

	if entry["level"] != wantLevel {
		t.Fatalf("level mismatch: got %v want %s", entry["level"], wantLevel)
	}
	if entry["message"] != wantMessage {
		t.Fatalf("message mismatch: got %v want %s", entry["message"], wantMessage)
	}
	if entry["request_id"] != wantRequestID {
		t.Fatalf("request_id mismatch: got %v want %s", entry["request_id"], wantRequestID)
	}
}
