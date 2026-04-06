package domain

import (
	"encoding/json"
	"reflect"
	"testing"
)

func TestTransactionRecordJSONRoundTrip(t *testing.T) {
	record := TransactionRecord{
		ID: "tx-1",
		Properties: TransactionProperties{
			Title:           titleProperty("12月保守費"),
			Customer:        RelationProperty{Type: "relation", Relation: []RelationRef{{ID: "customer-1"}}},
			CashDate:        dateProperty("2025-12-31"),
			Amount:          numberProperty(120000),
			Category:        selectProperty("売上"),
			IncludeCashflow: CheckboxProperty{Type: "checkbox", Checkbox: true},
			Description:     richTextProperty("月次保守"),
			Qty:             numberProperty(1),
			UnitPrice:       numberProperty(120000),
			LineTotal:       formulaNumberProperty(120000),
			BillingDate:     dateProperty("2025-12-31"),
			PaymentDue:      dateProperty("2026-01-31"),
			Invoice:         RelationProperty{Type: "relation", Relation: []RelationRef{{ID: "invoice-1"}}},
			Status:          statusProperty("未請求"),
		},
	}

	assertJSONRoundTrip(t, record)
}

func TestCustomerRecordJSONRoundTrip(t *testing.T) {
	record := CustomerRecord{
		ID: "customer-1",
		Properties: CustomerProperties{
			Title: titleProperty("株式会社サンプル"),
			BillingEmail: EmailProperty{
				Type:  "email",
				Email: "billing@example.com",
			},
			BillingTo:    richTextProperty("株式会社サンプル 御中"),
			SendMethod:   selectProperty("メール"),
			CustomerCode: richTextProperty("AB"),
			Memo:         richTextProperty("優先顧客"),
		},
	}

	assertJSONRoundTrip(t, record)
}

func TestInvoiceRecordJSONRoundTrip(t *testing.T) {
	record := InvoiceRecord{
		ID: "invoice-1",
		Properties: InvoiceProperties{
			Title:         titleProperty("S-AB202512-01"),
			Customer:      RelationProperty{Type: "relation", Relation: []RelationRef{{ID: "customer-1"}}},
			BillingDate:   dateProperty("2025-12-31"),
			BillingPeriod: richTextProperty("2025-12"),
			TotalAmount:   rollupNumberProperty(120000),
			Items:         RelationProperty{Type: "relation", Relation: []RelationRef{{ID: "tx-1"}}},
			PDFURL: URLProperty{
				Type: "files",
				Files: []FileObject{{
					Name: "invoice.pdf",
					External: &ExternalFile{
						URL: "https://example.com/invoice.pdf",
					},
				}},
			},
			ApprovalStatus: selectOrStatusProperty("select", "ドラフト"),
			EmailStatus:    selectProperty("未送信"),
			PaymentDue:     rollupDateProperty("2026-01-31"),
		},
	}

	assertJSONRoundTrip(t, record)
}

func TestCashflowRecordJSONRoundTrip(t *testing.T) {
	record := CashflowRecord{
		ID: "cashflow-1",
		Properties: CashflowProperties{
			Date:           dateProperty("2025-12-31"),
			OpeningBalance: numberProperty(500000),
			CashIn:         numberProperty(120000),
			CashOut:        numberProperty(30000),
			NetChange:      numberProperty(90000),
			ClosingBalance: numberProperty(590000),
			RiskFlag:       formulaStringProperty("safe"),
		},
	}

	assertJSONRoundTrip(t, record)
}

func assertJSONRoundTrip[T any](t *testing.T, record T) {
	t.Helper()

	data, err := json.Marshal(record)
	if err != nil {
		t.Fatalf("json.Marshal returned error: %v", err)
	}

	var got T
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatalf("json.Unmarshal returned error: %v", err)
	}

	if !reflect.DeepEqual(got, record) {
		t.Fatalf("round-trip mismatch: got %#v want %#v", got, record)
	}
}

func titleProperty(content string) TitleProperty {
	return TitleProperty{
		Type:  "title",
		Title: []RichTextObject{textObject(content)},
	}
}

func richTextProperty(content string) RichTextProperty {
	return RichTextProperty{
		Type:     "rich_text",
		RichText: []RichTextObject{textObject(content)},
	}
}

func textObject(content string) RichTextObject {
	return RichTextObject{
		Type:      "text",
		PlainText: content,
		Text: &TextContent{
			Content: content,
		},
	}
}

func dateProperty(start string) DateProperty {
	return DateProperty{
		Type: "date",
		Date: &DateValue{Start: start},
	}
}

func numberProperty(value float64) NumberProperty {
	return NumberProperty{
		Type:   "number",
		Number: &value,
	}
}

func selectProperty(name string) SelectProperty {
	return SelectProperty{
		Type: "select",
		Select: &SelectOption{
			Name: name,
		},
	}
}

func statusProperty(name string) StatusProperty {
	return selectOrStatusProperty("status", name)
}

func selectOrStatusProperty(kind, name string) StatusProperty {
	property := StatusProperty{Type: kind}
	option := &SelectOption{Name: name}
	if kind == "select" {
		property.Select = option
		return property
	}
	property.Status = option
	return property
}

func formulaNumberProperty(value float64) FormulaProperty {
	return FormulaProperty{
		Type: "formula",
		Formula: &FormulaValue{
			Type:   "number",
			Number: &value,
		},
	}
}

func formulaStringProperty(value string) FormulaProperty {
	return FormulaProperty{
		Type: "formula",
		Formula: &FormulaValue{
			Type:   "string",
			String: &value,
		},
	}
}

func rollupNumberProperty(value float64) RollupProperty {
	return RollupProperty{
		Type: "rollup",
		Rollup: &RollupValue{
			Type:   "number",
			Number: &value,
		},
	}
}

func rollupDateProperty(start string) RollupProperty {
	return RollupProperty{
		Type: "rollup",
		Rollup: &RollupValue{
			Type: "date",
			Date: &DateValue{Start: start},
		},
	}
}
