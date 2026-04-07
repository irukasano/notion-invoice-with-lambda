package acceptance_test

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/irukasano/notion-invoice-with-lambda/internal/domain"
	"github.com/irukasano/notion-invoice-with-lambda/internal/notion"
)

func TestIssue3Feature_NotionClient(t *testing.T) {
	t.Run("query uses notion database endpoint and parses records", func(t *testing.T) {
		var got requestCapture
		client := newTestClient(fakeDoer(func(r *http.Request) (*http.Response, error) {
			got = captureRequest(t, r)
			return encodeResponse(t, map[string]any{
				"results": []map[string]any{{
					"id": "tx-1",
					"properties": map[string]any{
						"title": map[string]any{
							"type": "title",
							"title": []map[string]any{{
								"type":       "text",
								"plain_text": "12月保守費",
								"text": map[string]any{
									"content": "12月保守費",
								},
							}},
						},
					},
				}},
				"has_more":    false,
				"next_cursor": nil,
			}), nil
		}))

		response, err := client.QueryTransactions(context.Background(), notion.QueryDatabaseRequest{
			Filter: map[string]any{
				"property": "status",
				"status": map[string]any{
					"equals": "未請求",
				},
			},
			PageSize: 10,
		})
		if err != nil {
			t.Fatalf("QueryTransactions returned error: %v", err)
		}
		if got.Method != http.MethodPost {
			t.Fatalf("method mismatch: got %q", got.Method)
		}
		if got.Path != "/v1/databases/tx-db/query" {
			t.Fatalf("path mismatch: got %q", got.Path)
		}
		if got.Authorization != "Bearer notion-secret" {
			t.Fatalf("authorization mismatch: got %q", got.Authorization)
		}
		if got.NotionVersion == "" {
			t.Fatal("Notion-Version header is empty")
		}
		if got.ContentType != "application/json" {
			t.Fatalf("content-type mismatch: got %q", got.ContentType)
		}
		if pageSize := got.Body["page_size"]; pageSize != float64(10) {
			t.Fatalf("page_size mismatch: got %#v", pageSize)
		}
		filter, ok := got.Body["filter"].(map[string]any)
		if !ok {
			t.Fatalf("filter type mismatch: got %#v", got.Body["filter"])
		}
		if filter["property"] != "status" {
			t.Fatalf("filter property mismatch: got %#v", filter["property"])
		}
		if len(response.Results) != 1 {
			t.Fatalf("result count mismatch: got %d", len(response.Results))
		}
		if response.Results[0].ID != "tx-1" {
			t.Fatalf("record id mismatch: got %q", response.Results[0].ID)
		}
	})

	t.Run("upsert maps create and update to pages endpoint", func(t *testing.T) {
		var calls []requestCapture
		client := newTestClient(fakeDoer(func(r *http.Request) (*http.Response, error) {
			calls = append(calls, captureRequest(t, r))
			switch len(calls) {
			case 1:
				return encodeResponse(t, map[string]any{
					"id": "customer-1",
					"properties": map[string]any{
						"title": map[string]any{
							"type": "title",
							"title": []map[string]any{{
								"type":       "text",
								"plain_text": "株式会社サンプル",
								"text": map[string]any{
									"content": "株式会社サンプル",
								},
							}},
						},
					},
				}), nil
			case 2:
				return encodeResponse(t, map[string]any{
					"id": "invoice-1",
					"properties": map[string]any{
						"title": map[string]any{
							"type": "title",
							"title": []map[string]any{{
								"type":       "text",
								"plain_text": "S-AB202512-01",
								"text": map[string]any{
									"content": "S-AB202512-01",
								},
							}},
						},
					},
				}), nil
			default:
				t.Fatalf("unexpected call count: %d", len(calls))
				return nil, nil
			}
		}))

		customer, err := client.UpsertCustomer(context.Background(), notion.UpsertPageRequest[domain.CustomerProperties]{
			Properties: domain.CustomerProperties{
				Title: titleProperty("株式会社サンプル"),
			},
		})
		if err != nil {
			t.Fatalf("UpsertCustomer returned error: %v", err)
		}
		if customer.ID != "customer-1" {
			t.Fatalf("customer id mismatch: got %q", customer.ID)
		}

		invoice, err := client.UpsertInvoice(context.Background(), notion.UpsertPageRequest[domain.InvoiceProperties]{
			PageID: "invoice-1",
			Properties: domain.InvoiceProperties{
				Title: titleProperty("S-AB202512-01"),
			},
		})
		if err != nil {
			t.Fatalf("UpsertInvoice returned error: %v", err)
		}
		if invoice.ID != "invoice-1" {
			t.Fatalf("invoice id mismatch: got %q", invoice.ID)
		}
		if len(calls) != 2 {
			t.Fatalf("call count mismatch: got %d", len(calls))
		}
		if calls[0].Method != http.MethodPost || calls[0].Path != "/v1/pages" {
			t.Fatalf("create endpoint mismatch: %#v", calls[0])
		}
		parent, ok := calls[0].Body["parent"].(map[string]any)
		if !ok {
			t.Fatalf("parent type mismatch: got %#v", calls[0].Body["parent"])
		}
		if parent["database_id"] != "customers-db" {
			t.Fatalf("parent database mismatch: got %#v", parent["database_id"])
		}
		if calls[1].Method != http.MethodPatch || calls[1].Path != "/v1/pages/invoice-1" {
			t.Fatalf("update endpoint mismatch: %#v", calls[1])
		}
		if _, exists := calls[1].Body["parent"]; exists {
			t.Fatalf("update request should not contain parent: %#v", calls[1].Body)
		}
	})

	t.Run("rate limit responses are classified", func(t *testing.T) {
		requestCount := 0
		client := newTestClient(fakeDoer(func(r *http.Request) (*http.Response, error) {
			requestCount++
			return jsonResponse(http.StatusTooManyRequests, `{"message":"rate limited"}`), nil
		}))
		_, err := client.QueryCashflow(context.Background(), notion.QueryDatabaseRequest{})
		if err == nil {
			t.Fatal("QueryCashflow returned nil error")
		}

		var rateLimitErr *notion.RateLimitError
		if !errors.As(err, &rateLimitErr) {
			t.Fatalf("expected RateLimitError, got %T", err)
		}
		if rateLimitErr.StatusCode != http.StatusTooManyRequests {
			t.Fatalf("status code mismatch: got %d", rateLimitErr.StatusCode)
		}
		if requestCount != 1 {
			t.Fatalf("unexpected retry count: got %d want 1", requestCount)
		}
	})

	t.Run("5xx responses are classified", func(t *testing.T) {
		requestCount := 0
		client := newTestClient(fakeDoer(func(r *http.Request) (*http.Response, error) {
			requestCount++
			return jsonResponse(http.StatusBadGateway, `{"message":"server error"}`), nil
		}))
		_, err := client.QueryInvoices(context.Background(), notion.QueryDatabaseRequest{})
		if err == nil {
			t.Fatal("QueryInvoices returned nil error")
		}

		var serverErr *notion.ServerError
		if !errors.As(err, &serverErr) {
			t.Fatalf("expected ServerError, got %T", err)
		}
		if serverErr.StatusCode != http.StatusBadGateway {
			t.Fatalf("status code mismatch: got %d", serverErr.StatusCode)
		}
		if requestCount != 1 {
			t.Fatalf("unexpected retry count: got %d want 1", requestCount)
		}
	})
}

type requestCapture struct {
	Method        string
	Path          string
	Authorization string
	NotionVersion string
	ContentType   string
	Body          map[string]any
}

func newTestClient(httpClient notion.HTTPDoer) *notion.Client {
	return notion.NewClient(notion.ClientConfig{
		BaseURL: "https://api.notion.test",
		Token:   "notion-secret",
		DatabaseIDs: notion.DatabaseIDs{
			Transactions: "tx-db",
			Customers:    "customers-db",
			Invoices:     "invoices-db",
			Cashflow:     "cashflow-db",
		},
	}, httpClient)
}

func captureRequest(t *testing.T, r *http.Request) requestCapture {
	t.Helper()

	body, err := io.ReadAll(r.Body)
	if err != nil {
		t.Fatalf("io.ReadAll returned error: %v", err)
	}
	defer r.Body.Close()

	captured := requestCapture{
		Method:        r.Method,
		Path:          r.URL.Path,
		Authorization: r.Header.Get("Authorization"),
		NotionVersion: r.Header.Get("Notion-Version"),
		ContentType:   r.Header.Get("Content-Type"),
		Body:          map[string]any{},
	}
	if len(body) == 0 {
		return captured
	}
	if err := json.Unmarshal(body, &captured.Body); err != nil {
		t.Fatalf("json.Unmarshal returned error: %v", err)
	}
	return captured
}

type fakeDoer func(*http.Request) (*http.Response, error)

func (f fakeDoer) Do(req *http.Request) (*http.Response, error) {
	return f(req)
}

func encodeResponse(t *testing.T, body map[string]any) *http.Response {
	t.Helper()

	builder := &strings.Builder{}
	if err := json.NewEncoder(builder).Encode(body); err != nil {
		t.Fatalf("json.NewEncoder returned error: %v", err)
	}

	return jsonResponse(http.StatusOK, builder.String())
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

func titleProperty(content string) domain.TitleProperty {
	return domain.TitleProperty{
		Type: "title",
		Title: []domain.RichTextObject{{
			Type:      "text",
			PlainText: content,
			Text: &domain.TextContent{
				Content: content,
			},
		}},
	}
}

func richTextProperty(content string) domain.RichTextProperty {
	return domain.RichTextProperty{
		Type: "rich_text",
		RichText: []domain.RichTextObject{{
			Type:      "text",
			PlainText: content,
			Text: &domain.TextContent{
				Content: content,
			},
		}},
	}
}

func emailProperty(email string) domain.EmailProperty {
	return domain.EmailProperty{
		Type:  "email",
		Email: email,
	}
}

func dateProperty(start string) domain.DateProperty {
	return domain.DateProperty{
		Type: "date",
		Date: &domain.DateValue{
			Start: start,
		},
	}
}

func numberProperty(value float64) domain.NumberProperty {
	return domain.NumberProperty{
		Type:   "number",
		Number: &value,
	}
}

func selectProperty(name string) domain.SelectProperty {
	return domain.SelectProperty{
		Type: "select",
		Select: &domain.SelectOption{
			Name: name,
		},
	}
}

func statusProperty(name string) domain.StatusProperty {
	return domain.StatusProperty{
		Type: "status",
		Status: &domain.SelectOption{
			Name: name,
		},
	}
}

func relationProperty(id string) domain.RelationProperty {
	return domain.RelationProperty{
		Type: "relation",
		Relation: []domain.RelationRef{{
			ID: id,
		}},
	}
}
