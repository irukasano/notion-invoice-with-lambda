package sendgrid

import (
	"context"
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"
)

func TestNewClientRequiresMandatoryFields(t *testing.T) {
	t.Run("api key is required", func(t *testing.T) {
		_, err := NewClient(ClientConfig{
			BaseURL:   "https://api.sendgrid.test/v3",
			FromEmail: "noreply@example.com",
			FromName:  "イルカシステム",
		}, doerFunc(nil))
		if err == nil {
			t.Fatal("NewClient error = nil")
		}
	})

	t.Run("from email must be valid", func(t *testing.T) {
		_, err := NewClient(ClientConfig{
			BaseURL:   "https://api.sendgrid.test/v3",
			APIKey:    "secret",
			FromEmail: "invalid",
			FromName:  "イルカシステム",
		}, doerFunc(nil))
		if err == nil {
			t.Fatal("NewClient error = nil")
		}
	})
}

func TestSendMailReturnsValidationErrorBeforeHTTP(t *testing.T) {
	requestCount := 0
	client := mustNewClient(t, doerFunc(func(r *http.Request) (*http.Response, error) {
		requestCount++
		return jsonResponse(http.StatusAccepted, `{"message":"accepted"}`), nil
	}))

	err := client.SendMail(context.Background(), SendMailInput{
		ToEmail:   "",
		Subject:   "subject",
		PlainText: "body",
	})
	if err == nil {
		t.Fatal("SendMail error = nil")
	}

	var validationErr *ValidationError
	if !errors.As(err, &validationErr) {
		t.Fatalf("expected ValidationError, got %T", err)
	}
	if validationErr.Field != "to_email" {
		t.Fatalf("field mismatch: got %q", validationErr.Field)
	}
	if requestCount != 0 {
		t.Fatalf("unexpected request count: got %d", requestCount)
	}
}

func TestSendMailClassifiesClientError(t *testing.T) {
	client := mustNewClient(t, doerFunc(func(r *http.Request) (*http.Response, error) {
		return jsonResponse(http.StatusBadRequest, `{"errors":[{"message":"bad request"}]}`), nil
	}))

	err := client.SendMail(context.Background(), SendMailInput{
		ToEmail:   "billing@example.com",
		Subject:   "subject",
		PlainText: "body",
	})
	if err == nil {
		t.Fatal("SendMail error = nil")
	}

	var clientErr *ClientError
	if !errors.As(err, &clientErr) {
		t.Fatalf("expected ClientError, got %T", err)
	}
	if clientErr.StatusCode != http.StatusBadRequest {
		t.Fatalf("status code mismatch: got %d", clientErr.StatusCode)
	}
}

func mustNewClient(t *testing.T, httpClient HTTPDoer) *Client {
	t.Helper()

	client, err := NewClient(ClientConfig{
		BaseURL:   "https://api.sendgrid.test/v3",
		APIKey:    "secret",
		FromEmail: "noreply@example.com",
		FromName:  "イルカシステム",
	}, httpClient)
	if err != nil {
		t.Fatalf("NewClient returned error: %v", err)
	}

	return client
}

type doerFunc func(*http.Request) (*http.Response, error)

func (f doerFunc) Do(req *http.Request) (*http.Response, error) {
	return f(req)
}

func jsonResponse(status int, body string) *http.Response {
	return &http.Response{
		StatusCode: status,
		Header: http.Header{
			"Content-Type": []string{"application/json"},
		},
		Body: io.NopCloser(strings.NewReader(body)),
	}
}
