package notion

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"
)

func TestNewClientAppliesDefaults(t *testing.T) {
	client := NewClient(ClientConfig{Token: "secret"}, nil)

	if client.baseURL != defaultBaseURL {
		t.Fatalf("baseURL mismatch: got %q", client.baseURL)
	}
	if client.notionVersion != defaultNotionVersion {
		t.Fatalf("Notion version mismatch: got %q", client.notionVersion)
	}
	if client.httpClient == nil {
		t.Fatal("httpClient is nil")
	}
}

func TestNewJSONRequestSetsHeaders(t *testing.T) {
	client := NewClient(ClientConfig{
		BaseURL:       "https://api.notion.test/",
		Token:         "secret",
		NotionVersion: "2022-06-28",
	}, nil)

	req, err := client.newJSONRequest(context.Background(), http.MethodPost, "/v1/pages", map[string]string{
		"foo": "bar",
	})
	if err != nil {
		t.Fatalf("newJSONRequest returned error: %v", err)
	}

	if req.URL.String() != "https://api.notion.test/v1/pages" {
		t.Fatalf("url mismatch: got %q", req.URL.String())
	}
	if req.Header.Get("Authorization") != "Bearer secret" {
		t.Fatalf("authorization mismatch: got %q", req.Header.Get("Authorization"))
	}
	if req.Header.Get("Content-Type") != "application/json" {
		t.Fatalf("content-type mismatch: got %q", req.Header.Get("Content-Type"))
	}
	if req.Header.Get("Notion-Version") != "2022-06-28" {
		t.Fatalf("version mismatch: got %q", req.Header.Get("Notion-Version"))
	}
}

func TestDoJSONDecodesResponse(t *testing.T) {
	client := NewClient(ClientConfig{BaseURL: "https://api.notion.test"}, fakeDoer(func(req *http.Request) (*http.Response, error) {
		recorder := &strings.Builder{}
		if err := json.NewEncoder(recorder).Encode(map[string]string{"id": "page-1"}); err != nil {
			t.Fatalf("json.NewEncoder returned error: %v", err)
		}

		return jsonResponse(http.StatusOK, recorder.String()), nil
	}))
	req, err := client.newJSONRequest(context.Background(), http.MethodGet, "/v1/pages/page-1", nil)
	if err != nil {
		t.Fatalf("newJSONRequest returned error: %v", err)
	}

	var got struct {
		ID string `json:"id"`
	}
	if err := client.doJSON(req, &got); err != nil {
		t.Fatalf("doJSON returned error: %v", err)
	}
	if got.ID != "page-1" {
		t.Fatalf("id mismatch: got %q", got.ID)
	}
}

func TestDoJSONClassifiesErrors(t *testing.T) {
	tests := []struct {
		name   string
		status int
		assert func(*testing.T, error)
	}{
		{
			name:   "rate limit",
			status: http.StatusTooManyRequests,
			assert: func(t *testing.T, err error) {
				t.Helper()

				var target *RateLimitError
				if !errors.As(err, &target) {
					t.Fatalf("unexpected error type: got %T", err)
				}
			},
		},
		{
			name:   "server error",
			status: http.StatusBadGateway,
			assert: func(t *testing.T, err error) {
				t.Helper()

				var target *ServerError
				if !errors.As(err, &target) {
					t.Fatalf("unexpected error type: got %T", err)
				}
			},
		},
		{
			name:   "client error",
			status: http.StatusBadRequest,
			assert: func(t *testing.T, err error) {
				t.Helper()

				var target *ClientError
				if !errors.As(err, &target) {
					t.Fatalf("unexpected error type: got %T", err)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := NewClient(ClientConfig{BaseURL: "https://api.notion.test"}, fakeDoer(func(req *http.Request) (*http.Response, error) {
				return jsonResponse(tt.status, "boom"), nil
			}))
			req, err := client.newJSONRequest(context.Background(), http.MethodGet, "/v1/test", nil)
			if err != nil {
				t.Fatalf("newJSONRequest returned error: %v", err)
			}

			err = client.doJSON(req, nil)
			if err == nil {
				t.Fatal("doJSON returned nil error")
			}
			tt.assert(t, err)
		})
	}
}

type fakeDoer func(*http.Request) (*http.Response, error)

func (f fakeDoer) Do(req *http.Request) (*http.Response, error) {
	return f(req)
}

func jsonResponse(statusCode int, body string) *http.Response {
	return &http.Response{
		StatusCode: statusCode,
		Body:       io.NopCloser(strings.NewReader(body)),
		Header: http.Header{
			"Content-Type": []string{"application/json"},
		},
	}
}
