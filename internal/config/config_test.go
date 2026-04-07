package config

import (
	"errors"
	"reflect"
	"slices"
	"testing"
)

func TestLoadFromLookupReturnsConfig(t *testing.T) {
	cfg, err := LoadFromLookup(lookupFromMap(map[string]string{
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
	}))
	if err != nil {
		t.Fatalf("LoadFromLookup returned error: %v", err)
	}

	want := Config{
		Env:                  "stg",
		S3Bucket:             "invoice-pdf-stg",
		NotionDBTransactions: "tx-db",
		NotionDBCustomers:    "customers-db",
		NotionDBInvoices:     "invoices-db",
		NotionDBCashflow:     "cashflow-db",
		NotionToken:          "secret",
		SendGridBaseURL:      "https://api.sendgrid.test",
		SendGridAPIKey:       "sendgrid-key",
		SendGridFromEmail:    "billing@example.com",
		SendGridFromName:     "Billing Bot",
		PDFTKPath:            defaultPDFTKPath,
	}
	if !reflect.DeepEqual(cfg, want) {
		t.Fatalf("config mismatch: got %#v want %#v", cfg, want)
	}
}

func TestLoadFromLookupUsesPDFTKOverride(t *testing.T) {
	cfg, err := LoadFromLookup(lookupFromMap(map[string]string{
		"ENV":                    "prod",
		"S3_BUCKET":              "invoice-pdf-prod",
		"NOTION_DB_TRANSACTIONS": "tx-db",
		"NOTION_DB_CUSTOMERS":    "customers-db",
		"NOTION_DB_INVOICES":     "invoices-db",
		"NOTION_DB_CASHFLOW":     "cashflow-db",
		"NOTION_TOKEN":           "secret",
		"SENDGRID_BASE_URL":      "https://api.sendgrid.test",
		"SENDGRID_API_KEY":       "sendgrid-key",
		"SENDGRID_FROM_EMAIL":    "billing@example.com",
		"SENDGRID_FROM_NAME":     "Billing Bot",
		"PDFTK_PATH":             "/opt/bin/pdftk",
	}))
	if err != nil {
		t.Fatalf("LoadFromLookup returned error: %v", err)
	}
	if cfg.PDFTKPath != "/opt/bin/pdftk" {
		t.Fatalf("PDFTKPath mismatch: got %q", cfg.PDFTKPath)
	}
}

func TestLoadFromLookupReturnsMissingEnvironmentError(t *testing.T) {
	_, err := LoadFromLookup(lookupFromMap(map[string]string{
		"ENV":                    "stg",
		"S3_BUCKET":              "invoice-pdf-stg",
		"NOTION_DB_TRANSACTIONS": "tx-db",
		"SENDGRID_BASE_URL":      "https://api.sendgrid.test",
		"SENDGRID_FROM_NAME":     "Billing Bot",
	}))
	if err == nil {
		t.Fatal("LoadFromLookup error = nil")
	}

	var missingErr MissingEnvironmentError
	if !errors.As(err, &missingErr) {
		t.Fatalf("error type mismatch: got %T", err)
	}

	want := []string{
		"NOTION_DB_CASHFLOW",
		"NOTION_DB_CUSTOMERS",
		"NOTION_DB_INVOICES",
		"NOTION_TOKEN",
		"SENDGRID_API_KEY",
		"SENDGRID_FROM_EMAIL",
	}
	if !slices.Equal(missingErr.Keys, want) {
		t.Fatalf("missing keys mismatch: got %v want %v", missingErr.Keys, want)
	}
}

func lookupFromMap(values map[string]string) func(string) (string, bool) {
	return func(key string) (string, bool) {
		value, ok := values[key]
		return value, ok
	}
}
