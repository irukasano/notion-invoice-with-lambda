package s3

import (
	"fmt"
	"strconv"
	"strings"
	"time"
	"unicode"
)

const issuerName = "イルカシステム"

type InvoicePDFKeyInput struct {
	BillingDate   time.Time
	InvoiceNumber string
	CustomerName  string
	TotalAmount   float64
}

func InvoicePDFKey(input InvoicePDFKeyInput) string {
	billingDate := input.BillingDate.Format("2006-01-02")
	yearMonth := input.BillingDate.Format("2006/01")

	return fmt.Sprintf(
		"%s/%s_%s_%s_%s_%s.pdf",
		yearMonth,
		billingDate,
		input.InvoiceNumber,
		NormalizeFileNameSegment(input.CustomerName),
		issuerName,
		formatAmount(input.TotalAmount),
	)
}

func NormalizeFileNameSegment(value string) string {
	var builder strings.Builder
	lastUnderscore := false

	for _, r := range strings.TrimSpace(value) {
		if isFileNameSafeRune(r) {
			builder.WriteRune(r)
			lastUnderscore = false
			continue
		}
		if !lastUnderscore {
			builder.WriteByte('_')
			lastUnderscore = true
		}
	}

	return strings.Trim(builder.String(), "_")
}

func isFileNameSafeRune(r rune) bool {
	switch {
	case unicode.IsLetter(r), unicode.IsDigit(r):
		return true
	case r == '-', r == '_':
		return true
	default:
		return false
	}
}

func formatAmount(value float64) string {
	if value == float64(int64(value)) {
		return strconv.FormatInt(int64(value), 10)
	}

	return strconv.FormatFloat(value, 'f', -1, 64)
}
