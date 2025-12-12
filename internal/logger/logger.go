package logger

import (
	"context"
	"encoding/json"
	"io"
	"os"

	"github.com/aws/aws-lambda-go/lambdacontext"
)

type Level string

const (
	LevelDebug Level = "DEBUG"
	LevelInfo  Level = "INFO"
	LevelWarn  Level = "WARN"
	LevelError Level = "ERROR"

	defaultStage = "local"
)

type Option func(*Logger)

// WithStage sets the stage label for all logs.
func WithStage(stage string) Option {
	return func(l *Logger) {
		l.stage = stage
	}
}

// WithOutput overrides the output writer, primarily for tests.
func WithOutput(w io.Writer) Option {
	return func(l *Logger) {
		l.out = w
	}
}

type Logger struct {
	stage string
	out   io.Writer
}

// New constructs a logger with JSON output. Stage is read from the STAGE env var if not provided.
func New(opts ...Option) *Logger {
	l := &Logger{
		stage: os.Getenv("STAGE"),
		out:   os.Stdout,
	}

	if l.stage == "" {
		l.stage = defaultStage
	}

	for _, opt := range opts {
		opt(l)
	}

	return l
}

type Entry struct {
	logger    *Logger
	requestID string
	fields    map[string]any
}

// WithRequestID attaches a request ID to the log entry.
func (l *Logger) WithRequestID(requestID string) *Entry {
	return &Entry{
		logger:    l,
		requestID: requestID,
	}
}

// WithContext extracts a request ID from the context (Lambda or injected) and attaches it to the entry.
func (l *Logger) WithContext(ctx context.Context) *Entry {
	return l.WithRequestID(RequestIDFromContext(ctx))
}

// WithFields attaches structured fields to the entry.
func (e *Entry) WithFields(fields map[string]any) *Entry {
	merged := make(map[string]any, len(e.fields)+len(fields))
	for k, v := range e.fields {
		merged[k] = v
	}
	for k, v := range fields {
		merged[k] = v
	}
	return &Entry{
		logger:    e.logger,
		requestID: e.requestID,
		fields:    merged,
	}
}

func (e *Entry) Debug(message string, fields map[string]any) error {
	return e.log(LevelDebug, message, fields)
}

func (e *Entry) Info(message string, fields map[string]any) error {
	return e.log(LevelInfo, message, fields)
}

func (e *Entry) Warn(message string, fields map[string]any) error {
	return e.log(LevelWarn, message, fields)
}

func (e *Entry) Error(message string, fields map[string]any) error {
	return e.log(LevelError, message, fields)
}

func (e *Entry) log(level Level, message string, fields map[string]any) error {
	payload := map[string]any{
		"stage":   e.logger.stage,
		"level":   level,
		"message": message,
	}

	if e.requestID != "" {
		payload["request_id"] = e.requestID
	}

	for k, v := range e.fields {
		payload[k] = v
	}

	for k, v := range fields {
		payload[k] = v
	}

	encoder := json.NewEncoder(e.logger.out)
	return encoder.Encode(payload)
}

type contextKey struct{}

// RequestIDFromContext resolves a request ID from Lambda context or injected value.
func RequestIDFromContext(ctx context.Context) string {
	if lc, ok := lambdacontext.FromContext(ctx); ok && lc.AwsRequestID != "" {
		return lc.AwsRequestID
	}

	if v := ctx.Value(contextKey{}); v != nil {
		if requestID, ok := v.(string); ok {
			return requestID
		}
	}

	return ""
}

// ContextWithRequestID injects a request ID into context for local runs/tests.
func ContextWithRequestID(ctx context.Context, requestID string) context.Context {
	return context.WithValue(ctx, contextKey{}, requestID)
}
