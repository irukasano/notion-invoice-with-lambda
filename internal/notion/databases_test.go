package notion

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/irukasano/notion-invoice-with-lambda/internal/domain"
)

func TestQueryTransactionsUsesTransactionsDatabaseID(t *testing.T) {
	var got requestSnapshot
	client := testClient(fakeDoer(func(r *http.Request) (*http.Response, error) {
		got = captureSnapshot(t, r)
		return encodeResponse(t, map[string]any{
			"results": []map[string]any{},
		}), nil
	}))
	_, err := client.QueryTransactions(context.Background(), QueryDatabaseRequest{
		PageSize: 25,
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
	if got.Body["page_size"] != float64(25) {
		t.Fatalf("page_size mismatch: got %#v", got.Body["page_size"])
	}
}

func TestUpsertCashflowCreatesPageInCashflowDatabase(t *testing.T) {
	var got requestSnapshot
	client := testClient(fakeDoer(func(r *http.Request) (*http.Response, error) {
		got = captureSnapshot(t, r)
		return encodeResponse(t, map[string]any{
			"id": "cashflow-1",
		}), nil
	}))
	record, err := client.UpsertCashflow(context.Background(), UpsertPageRequest[domain.CashflowProperties]{
		Properties: domain.CashflowProperties{
			Date: domain.DateProperty{
				Type: "date",
				Date: &domain.DateValue{Start: "2025-12-31"},
			},
		},
	})
	if err != nil {
		t.Fatalf("UpsertCashflow returned error: %v", err)
	}
	if record.ID != "cashflow-1" {
		t.Fatalf("record id mismatch: got %q", record.ID)
	}
	if got.Method != http.MethodPost {
		t.Fatalf("method mismatch: got %q", got.Method)
	}
	if got.Path != "/v1/pages" {
		t.Fatalf("path mismatch: got %q", got.Path)
	}
	parent, ok := got.Body["parent"].(map[string]any)
	if !ok {
		t.Fatalf("parent type mismatch: got %#v", got.Body["parent"])
	}
	if parent["database_id"] != "cashflow-db" {
		t.Fatalf("database mismatch: got %#v", parent["database_id"])
	}
}

func TestUpsertInvoiceUpdatesExistingPage(t *testing.T) {
	var got requestSnapshot
	client := testClient(fakeDoer(func(r *http.Request) (*http.Response, error) {
		got = captureSnapshot(t, r)
		return encodeResponse(t, map[string]any{
			"id": "invoice-1",
		}), nil
	}))
	record, err := client.UpsertInvoice(context.Background(), UpsertPageRequest[domain.InvoiceProperties]{
		PageID: "invoice-1",
		Properties: domain.InvoiceProperties{
			Title: domain.TitleProperty{
				Type: "title",
			},
		},
	})
	if err != nil {
		t.Fatalf("UpsertInvoice returned error: %v", err)
	}
	if record.ID != "invoice-1" {
		t.Fatalf("record id mismatch: got %q", record.ID)
	}
	if got.Method != http.MethodPatch {
		t.Fatalf("method mismatch: got %q", got.Method)
	}
	if got.Path != "/v1/pages/invoice-1" {
		t.Fatalf("path mismatch: got %q", got.Path)
	}
	if _, ok := got.Body["parent"]; ok {
		t.Fatalf("unexpected parent in update body: %#v", got.Body)
	}
}

type requestSnapshot struct {
	Method string
	Path   string
	Body   map[string]any
}

func testClient(httpClient HTTPDoer) *Client {
	return NewClient(ClientConfig{
		BaseURL: "https://api.notion.test",
		Token:   "secret",
		DatabaseIDs: DatabaseIDs{
			Transactions: "tx-db",
			Customers:    "customers-db",
			Invoices:     "invoices-db",
			Cashflow:     "cashflow-db",
		},
	}, httpClient)
}

func captureSnapshot(t *testing.T, r *http.Request) requestSnapshot {
	t.Helper()

	body, err := io.ReadAll(r.Body)
	if err != nil {
		t.Fatalf("io.ReadAll returned error: %v", err)
	}
	defer r.Body.Close()

	result := requestSnapshot{
		Method: r.Method,
		Path:   r.URL.Path,
		Body:   map[string]any{},
	}
	if len(body) == 0 {
		return result
	}
	if err := json.Unmarshal(body, &result.Body); err != nil {
		t.Fatalf("json.Unmarshal returned error: %v", err)
	}
	return result
}

func encodeResponse(t *testing.T, body map[string]any) *http.Response {
	t.Helper()

	recorder := &strings.Builder{}
	if err := json.NewEncoder(recorder).Encode(body); err != nil {
		t.Fatalf("json.NewEncoder returned error: %v", err)
	}

	return jsonResponse(http.StatusOK, recorder.String())
}
