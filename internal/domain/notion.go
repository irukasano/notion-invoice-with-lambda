package domain

type TransactionRecord struct {
	ID         string                `json:"id,omitempty"`
	Properties TransactionProperties `json:"properties"`
}

type TransactionProperties struct {
	Title           TitleProperty    `json:"title"`
	Customer        RelationProperty `json:"customer,omitempty"`
	CashDate        DateProperty     `json:"cash_date,omitempty"`
	Amount          NumberProperty   `json:"amount,omitempty"`
	Category        SelectProperty   `json:"category,omitempty"`
	IncludeCashflow CheckboxProperty `json:"include_cashflow"`
	Description     RichTextProperty `json:"description,omitempty"`
	Qty             NumberProperty   `json:"qty,omitempty"`
	UnitPrice       NumberProperty   `json:"unit_price,omitempty"`
	LineTotal       FormulaProperty  `json:"line_total,omitempty"`
	BillingDate     DateProperty     `json:"billing_date,omitempty"`
	PaymentDue      DateProperty     `json:"payment_due,omitempty"`
	Invoice         RelationProperty `json:"invoice,omitempty"`
	Status          StatusProperty   `json:"status,omitempty"`
}

type CustomerRecord struct {
	ID         string             `json:"id,omitempty"`
	Properties CustomerProperties `json:"properties"`
}

type CustomerProperties struct {
	Title        TitleProperty    `json:"title"`
	BillingEmail EmailProperty    `json:"billing_email,omitempty"`
	BillingTo    RichTextProperty `json:"billing_to,omitempty"`
	SendMethod   SelectProperty   `json:"send_method,omitempty"`
	CustomerCode RichTextProperty `json:"customer_code,omitempty"`
	Memo         RichTextProperty `json:"memo,omitempty"`
}

type InvoiceRecord struct {
	ID         string            `json:"id,omitempty"`
	Properties InvoiceProperties `json:"properties"`
}

type InvoiceProperties struct {
	Title          TitleProperty    `json:"title"`
	Customer       RelationProperty `json:"customer,omitempty"`
	BillingDate    DateProperty     `json:"billing_date,omitempty"`
	BillingPeriod  RichTextProperty `json:"billing_period,omitempty"`
	TotalAmount    RollupProperty   `json:"total_amount,omitempty"`
	Items          RelationProperty `json:"items,omitempty"`
	PDFURL         URLProperty      `json:"pdf_url,omitempty"`
	ApprovalStatus StatusProperty   `json:"approval_status,omitempty"`
	EmailStatus    SelectProperty   `json:"email_status,omitempty"`
	PaymentDue     RollupProperty   `json:"payment_due,omitempty"`
}

type CashflowRecord struct {
	ID         string             `json:"id,omitempty"`
	Properties CashflowProperties `json:"properties"`
}

type CashflowProperties struct {
	Date           DateProperty    `json:"date,omitempty"`
	OpeningBalance NumberProperty  `json:"opening_balance,omitempty"`
	CashIn         NumberProperty  `json:"cash_in,omitempty"`
	CashOut        NumberProperty  `json:"cash_out,omitempty"`
	NetChange      NumberProperty  `json:"net_change,omitempty"`
	ClosingBalance NumberProperty  `json:"closing_balance,omitempty"`
	RiskFlag       FormulaProperty `json:"risk_flag,omitempty"`
}

type TitleProperty struct {
	Type  string           `json:"type,omitempty"`
	Title []RichTextObject `json:"title"`
}

type RichTextProperty struct {
	Type     string           `json:"type,omitempty"`
	RichText []RichTextObject `json:"rich_text"`
}

type RichTextObject struct {
	Type      string       `json:"type,omitempty"`
	PlainText string       `json:"plain_text,omitempty"`
	Text      *TextContent `json:"text,omitempty"`
}

type TextContent struct {
	Content string `json:"content"`
}

type RelationProperty struct {
	Type     string        `json:"type,omitempty"`
	Relation []RelationRef `json:"relation"`
}

type RelationRef struct {
	ID string `json:"id"`
}

type DateProperty struct {
	Type string     `json:"type,omitempty"`
	Date *DateValue `json:"date,omitempty"`
}

type DateValue struct {
	Start string `json:"start"`
	End   string `json:"end,omitempty"`
}

type NumberProperty struct {
	Type   string   `json:"type,omitempty"`
	Number *float64 `json:"number,omitempty"`
}

type SelectProperty struct {
	Type   string        `json:"type,omitempty"`
	Select *SelectOption `json:"select,omitempty"`
}

type StatusProperty struct {
	Type   string        `json:"type,omitempty"`
	Select *SelectOption `json:"select,omitempty"`
	Status *SelectOption `json:"status,omitempty"`
}

type SelectOption struct {
	Name string `json:"name"`
}

type CheckboxProperty struct {
	Type     string `json:"type,omitempty"`
	Checkbox bool   `json:"checkbox"`
}

type FormulaProperty struct {
	Type    string        `json:"type,omitempty"`
	Formula *FormulaValue `json:"formula,omitempty"`
}

type FormulaValue struct {
	Type   string   `json:"type"`
	String *string  `json:"string,omitempty"`
	Number *float64 `json:"number,omitempty"`
}

type RollupProperty struct {
	Type   string       `json:"type,omitempty"`
	Rollup *RollupValue `json:"rollup,omitempty"`
}

type RollupValue struct {
	Type   string     `json:"type"`
	Number *float64   `json:"number,omitempty"`
	Date   *DateValue `json:"date,omitempty"`
}

type EmailProperty struct {
	Type     string           `json:"type,omitempty"`
	Email    string           `json:"email,omitempty"`
	RichText []RichTextObject `json:"rich_text,omitempty"`
}

type URLProperty struct {
	Type  string       `json:"type,omitempty"`
	URL   string       `json:"url,omitempty"`
	Files []FileObject `json:"files,omitempty"`
}

type FileObject struct {
	Name     string        `json:"name,omitempty"`
	External *ExternalFile `json:"external,omitempty"`
}

type ExternalFile struct {
	URL string `json:"url"`
}
