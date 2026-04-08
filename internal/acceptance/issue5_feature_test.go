package acceptance_test

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/irukasano/notion-invoice-with-lambda/internal/sendgrid"
)

func TestIssue5Feature_SendGridClient(t *testing.T) {
	t.Run("mail send request includes headers body and base64 attachment", func(t *testing.T) {
		var got sendGridRequestCapture
		client := newTestSendGridClient(t, sendGridDoerFunc(func(r *http.Request) (*http.Response, error) {
			got = captureSendGridRequest(t, r)
			return sendGridJSONResponse(http.StatusAccepted, `{"message":"accepted"}`), nil
		}))

		err := client.SendMail(context.Background(), sendgrid.SendMailInput{
			ToEmail:   "billing@example.com",
			ToName:    "経理担当",
			Subject:   "[請求書] 株式会社サンプル様 2025-12分 ご請求のご案内",
			PlainText: "2025-12分の請求書をお送りします。",
			Attachments: []sendgrid.Attachment{
				{
					Filename:    "invoice.pdf",
					ContentType: "application/pdf",
					Content:     []byte("%PDF-1.4"),
				},
			},
		})
		if err != nil {
			t.Fatalf("SendMail returned error: %v", err)
		}

		if got.Method != http.MethodPost {
			t.Fatalf("method mismatch: got %q", got.Method)
		}
		if got.Path != "/v3/mail/send" {
			t.Fatalf("path mismatch: got %q", got.Path)
		}
		if got.Authorization != "Bearer sendgrid-secret" {
			t.Fatalf("authorization mismatch: got %q", got.Authorization)
		}
		if got.ContentType != "application/json" {
			t.Fatalf("content type mismatch: got %q", got.ContentType)
		}

		from, ok := got.Body["from"].(map[string]any)
		if !ok {
			t.Fatalf("from type mismatch: got %#v", got.Body["from"])
		}
		if from["email"] != "noreply@example.com" {
			t.Fatalf("from email mismatch: got %#v", from["email"])
		}
		if from["name"] != "イルカシステム" {
			t.Fatalf("from name mismatch: got %#v", from["name"])
		}

		personalizations, ok := got.Body["personalizations"].([]any)
		if !ok || len(personalizations) != 1 {
			t.Fatalf("personalizations mismatch: got %#v", got.Body["personalizations"])
		}
		personalization, ok := personalizations[0].(map[string]any)
		if !ok {
			t.Fatalf("personalization type mismatch: got %#v", personalizations[0])
		}
		toList, ok := personalization["to"].([]any)
		if !ok || len(toList) != 1 {
			t.Fatalf("to list mismatch: got %#v", personalization["to"])
		}
		to, ok := toList[0].(map[string]any)
		if !ok {
			t.Fatalf("to type mismatch: got %#v", toList[0])
		}
		if to["email"] != "billing@example.com" {
			t.Fatalf("to email mismatch: got %#v", to["email"])
		}
		if to["name"] != "経理担当" {
			t.Fatalf("to name mismatch: got %#v", to["name"])
		}

		if got.Body["subject"] != "[請求書] 株式会社サンプル様 2025-12分 ご請求のご案内" {
			t.Fatalf("subject mismatch: got %#v", got.Body["subject"])
		}

		contentList, ok := got.Body["content"].([]any)
		if !ok || len(contentList) != 1 {
			t.Fatalf("content mismatch: got %#v", got.Body["content"])
		}
		content, ok := contentList[0].(map[string]any)
		if !ok {
			t.Fatalf("content type mismatch: got %#v", contentList[0])
		}
		if content["type"] != "text/plain" {
			t.Fatalf("content type mismatch: got %#v", content["type"])
		}
		if content["value"] != "2025-12分の請求書をお送りします。" {
			t.Fatalf("content value mismatch: got %#v", content["value"])
		}

		attachments, ok := got.Body["attachments"].([]any)
		if !ok || len(attachments) != 1 {
			t.Fatalf("attachments mismatch: got %#v", got.Body["attachments"])
		}
		attachment, ok := attachments[0].(map[string]any)
		if !ok {
			t.Fatalf("attachment type mismatch: got %#v", attachments[0])
		}
		if attachment["filename"] != "invoice.pdf" {
			t.Fatalf("filename mismatch: got %#v", attachment["filename"])
		}
		if attachment["type"] != "application/pdf" {
			t.Fatalf("attachment type mismatch: got %#v", attachment["type"])
		}
		if attachment["disposition"] != "attachment" {
			t.Fatalf("disposition mismatch: got %#v", attachment["disposition"])
		}
		if attachment["content"] != "JVBERi0xLjQ=" {
			t.Fatalf("attachment content mismatch: got %#v", attachment["content"])
		}
	})

	t.Run("recipient validation fails before http request", func(t *testing.T) {
		requestCount := 0
		client := newTestSendGridClient(t, sendGridDoerFunc(func(r *http.Request) (*http.Response, error) {
			requestCount++
			return sendGridJSONResponse(http.StatusAccepted, `{"message":"accepted"}`), nil
		}))

		err := client.SendMail(context.Background(), sendgrid.SendMailInput{
			ToEmail:   "invalid-address",
			Subject:   "subject",
			PlainText: "body",
		})
		if err == nil {
			t.Fatal("SendMail returned nil error")
		}

		var validationErr *sendgrid.ValidationError
		if !errors.As(err, &validationErr) {
			t.Fatalf("expected ValidationError, got %T", err)
		}
		if validationErr.Field != "to_email" {
			t.Fatalf("field mismatch: got %q", validationErr.Field)
		}
		if requestCount != 0 {
			t.Fatalf("unexpected request count: got %d", requestCount)
		}
	})

	t.Run("error responses are classified", func(t *testing.T) {
		t.Run("429", func(t *testing.T) {
			client := newTestSendGridClient(t, sendGridDoerFunc(func(r *http.Request) (*http.Response, error) {
				return sendGridJSONResponse(http.StatusTooManyRequests, `{"errors":[{"message":"rate limit"}]}`), nil
			}))

			err := client.SendMail(context.Background(), sendgrid.SendMailInput{
				ToEmail:   "billing@example.com",
				Subject:   "subject",
				PlainText: "body",
			})
			if err == nil {
				t.Fatal("SendMail returned nil error")
			}

			var rateLimitErr *sendgrid.RateLimitError
			if !errors.As(err, &rateLimitErr) {
				t.Fatalf("expected RateLimitError, got %T", err)
			}
		})

		t.Run("5xx", func(t *testing.T) {
			client := newTestSendGridClient(t, sendGridDoerFunc(func(r *http.Request) (*http.Response, error) {
				return sendGridJSONResponse(http.StatusBadGateway, `{"errors":[{"message":"bad gateway"}]}`), nil
			}))

			err := client.SendMail(context.Background(), sendgrid.SendMailInput{
				ToEmail:   "billing@example.com",
				Subject:   "subject",
				PlainText: "body",
			})
			if err == nil {
				t.Fatal("SendMail returned nil error")
			}

			var serverErr *sendgrid.ServerError
			if !errors.As(err, &serverErr) {
				t.Fatalf("expected ServerError, got %T", err)
			}
		})
	})
}

type sendGridRequestCapture struct {
	Method        string
	Path          string
	Authorization string
	ContentType   string
	Body          map[string]any
}

func newTestSendGridClient(t *testing.T, httpClient sendgrid.HTTPDoer) *sendgrid.Client {
	t.Helper()

	client, err := sendgrid.NewClient(sendgrid.ClientConfig{
		BaseURL:   "https://api.sendgrid.test/v3",
		APIKey:    "sendgrid-secret",
		FromEmail: "noreply@example.com",
		FromName:  "イルカシステム",
	}, httpClient)
	if err != nil {
		t.Fatalf("NewClient returned error: %v", err)
	}

	return client
}

func captureSendGridRequest(t *testing.T, r *http.Request) sendGridRequestCapture {
	t.Helper()

	body, err := io.ReadAll(r.Body)
	if err != nil {
		t.Fatalf("io.ReadAll returned error: %v", err)
	}
	defer r.Body.Close()

	var payload map[string]any
	if len(body) > 0 {
		if err := json.Unmarshal(body, &payload); err != nil {
			t.Fatalf("json.Unmarshal returned error: %v", err)
		}
	} else {
		payload = map[string]any{}
	}

	return sendGridRequestCapture{
		Method:        r.Method,
		Path:          r.URL.Path,
		Authorization: r.Header.Get("Authorization"),
		ContentType:   r.Header.Get("Content-Type"),
		Body:          payload,
	}
}

type sendGridDoerFunc func(*http.Request) (*http.Response, error)

func (f sendGridDoerFunc) Do(req *http.Request) (*http.Response, error) {
	return f(req)
}

func sendGridJSONResponse(status int, body string) *http.Response {
	return &http.Response{
		StatusCode: status,
		Header: http.Header{
			"Content-Type": []string{"application/json"},
		},
		Body: io.NopCloser(strings.NewReader(body)),
	}
}
