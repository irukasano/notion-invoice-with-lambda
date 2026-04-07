package acceptance_test

import (
	"encoding/json"
	"testing"

	"github.com/irukasano/notion-invoice-with-lambda/internal/config"
	"github.com/irukasano/notion-invoice-with-lambda/internal/domain"
)

func TestIssue2Feature_ConfigAndNotionRecords(t *testing.T) {
	cfg, err := config.LoadFromLookup(func(key string) (string, bool) {
		values := map[string]string{
			"ENV":                    "stg",
			"S3_BUCKET":              "invoice-pdf-stg",
			"NOTION_DB_TRANSACTIONS": "tx-db",
			"NOTION_DB_CUSTOMERS":    "customers-db",
			"NOTION_DB_INVOICES":     "invoices-db",
			"NOTION_DB_CASHFLOW":     "cashflow-db",
			"NOTION_TOKEN":           "secret",
			"SENDGRID_BASE_URL":      "https://api.sendgrid.test",
			"SENDGRID_API_KEY":       "sendgrid-key",
			"SENDGRID_FROM_EMAIL":    "billing@example.com",
			"SENDGRID_FROM_NAME":     "Billing Bot",
		}
		value, ok := values[key]
		return value, ok
	})
	if err != nil {
		t.Fatalf("LoadFromLookup returned error: %v", err)
	}
	if cfg.PDFTKPath != "pdftk" {
		t.Fatalf("default PDFTKPath mismatch: got %q", cfg.PDFTKPath)
	}

	records := []any{
		domain.TransactionRecord{
			Properties: domain.TransactionProperties{
				Title: domain.TitleProperty{
					Type: "title",
					Title: []domain.RichTextObject{{
						Type:      "text",
						PlainText: "12月保守費",
						Text: &domain.TextContent{
							Content: "12月保守費",
						},
					}},
				},
			},
		},
		domain.CustomerRecord{
			Properties: domain.CustomerProperties{
				Title: domain.TitleProperty{
					Type: "title",
					Title: []domain.RichTextObject{{
						Type:      "text",
						PlainText: "株式会社サンプル",
						Text: &domain.TextContent{
							Content: "株式会社サンプル",
						},
					}},
				},
			},
		},
		domain.InvoiceRecord{
			Properties: domain.InvoiceProperties{
				Title: domain.TitleProperty{
					Type: "title",
					Title: []domain.RichTextObject{{
						Type:      "text",
						PlainText: "S-AB202512-01",
						Text: &domain.TextContent{
							Content: "S-AB202512-01",
						},
					}},
				},
			},
		},
		domain.CashflowRecord{
			Properties: domain.CashflowProperties{
				Date: domain.DateProperty{
					Type: "date",
					Date: &domain.DateValue{
						Start: "2025-12-31",
					},
				},
			},
		},
	}

	for _, record := range records {
		if _, err := json.Marshal(record); err != nil {
			t.Fatalf("json.Marshal returned error: %v", err)
		}
	}
}
