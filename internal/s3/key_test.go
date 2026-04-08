package s3

import (
	"testing"
	"time"
)

func TestInvoicePDFKey(t *testing.T) {
	got := InvoicePDFKey(InvoicePDFKeyInput{
		BillingDate:   time.Date(2025, 2, 3, 0, 0, 0, 0, time.UTC),
		InvoiceNumber: "S-AB202502-01",
		CustomerName:  "株式会社サンプル / 東京",
		TotalAmount:   120000,
	})

	want := "2025/02/2025-02-03_S-AB202502-01_株式会社サンプル_東京_イルカシステム_120000.pdf"
	if got != want {
		t.Fatalf("InvoicePDFKey mismatch: got %q want %q", got, want)
	}
}

func TestNormalizeFileNameSegment(t *testing.T) {
	got := NormalizeFileNameSegment(" ACME / 東京(本社) ")
	want := "ACME_東京_本社"
	if got != want {
		t.Fatalf("NormalizeFileNameSegment mismatch: got %q want %q", got, want)
	}
}
