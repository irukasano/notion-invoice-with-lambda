package config

import (
	"fmt"
	"os"
	"sort"
	"strings"
)

const defaultPDFTKPath = "pdftk"

type Config struct {
	Env                  string
	S3Bucket             string
	NotionDBTransactions string
	NotionDBCustomers    string
	NotionDBInvoices     string
	NotionDBCashflow     string
	NotionToken          string
	SendGridBaseURL      string
	SendGridAPIKey       string
	SendGridFromEmail    string
	SendGridFromName     string
	PDFTKPath            string
}

type MissingEnvironmentError struct {
	Keys []string
}

func (e MissingEnvironmentError) Error() string {
	return fmt.Sprintf(
		"missing required environment variables: %s",
		strings.Join(e.Keys, ", "),
	)
}

func Load() (Config, error) {
	return LoadFromLookup(os.LookupEnv)
}

func LoadFromLookup(lookup func(string) (string, bool)) (Config, error) {
	cfg := Config{
		Env:                  lookupRequired(lookup, "ENV"),
		S3Bucket:             lookupRequired(lookup, "S3_BUCKET"),
		NotionDBTransactions: lookupRequired(lookup, "NOTION_DB_TRANSACTIONS"),
		NotionDBCustomers:    lookupRequired(lookup, "NOTION_DB_CUSTOMERS"),
		NotionDBInvoices:     lookupRequired(lookup, "NOTION_DB_INVOICES"),
		NotionDBCashflow:     lookupRequired(lookup, "NOTION_DB_CASHFLOW"),
		NotionToken:          lookupRequired(lookup, "NOTION_TOKEN"),
		SendGridBaseURL:      lookupRequired(lookup, "SENDGRID_BASE_URL"),
		SendGridAPIKey:       lookupRequired(lookup, "SENDGRID_API_KEY"),
		SendGridFromEmail:    lookupRequired(lookup, "SENDGRID_FROM_EMAIL"),
		SendGridFromName:     lookupRequired(lookup, "SENDGRID_FROM_NAME"),
		PDFTKPath:            lookupOptional(lookup, "PDFTK_PATH", defaultPDFTKPath),
	}

	missing := cfg.missingRequired()
	if len(missing) > 0 {
		return Config{}, MissingEnvironmentError{Keys: missing}
	}

	return cfg, nil
}

func lookupRequired(lookup func(string) (string, bool), key string) string {
	value, _ := lookup(key)
	return value
}

func lookupOptional(lookup func(string) (string, bool), key, fallback string) string {
	value, ok := lookup(key)
	if !ok || value == "" {
		return fallback
	}
	return value
}

func (c Config) missingRequired() []string {
	missing := make([]string, 0, 11)
	required := map[string]string{
		"ENV":                    c.Env,
		"NOTION_DB_CASHFLOW":     c.NotionDBCashflow,
		"NOTION_DB_CUSTOMERS":    c.NotionDBCustomers,
		"NOTION_DB_INVOICES":     c.NotionDBInvoices,
		"NOTION_DB_TRANSACTIONS": c.NotionDBTransactions,
		"NOTION_TOKEN":           c.NotionToken,
		"S3_BUCKET":              c.S3Bucket,
		"SENDGRID_API_KEY":       c.SendGridAPIKey,
		"SENDGRID_BASE_URL":      c.SendGridBaseURL,
		"SENDGRID_FROM_EMAIL":    c.SendGridFromEmail,
		"SENDGRID_FROM_NAME":     c.SendGridFromName,
	}

	for key, value := range required {
		if strings.TrimSpace(value) == "" {
			missing = append(missing, key)
		}
	}
	sort.Strings(missing)

	return missing
}
