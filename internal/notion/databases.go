package notion

import (
	"context"
	"errors"

	"github.com/irukasano/notion-invoice-with-lambda/internal/domain"
)

var errDatabaseIDRequired = errors.New("notion database id is required")

func (c *Client) QueryTransactions(ctx context.Context, request QueryDatabaseRequest) (QueryDatabaseResponse[domain.TransactionRecord], error) {
	return queryRecords[domain.TransactionRecord](ctx, c, c.databaseIDs.Transactions, request)
}

func (c *Client) QueryCustomers(ctx context.Context, request QueryDatabaseRequest) (QueryDatabaseResponse[domain.CustomerRecord], error) {
	return queryRecords[domain.CustomerRecord](ctx, c, c.databaseIDs.Customers, request)
}

func (c *Client) QueryInvoices(ctx context.Context, request QueryDatabaseRequest) (QueryDatabaseResponse[domain.InvoiceRecord], error) {
	return queryRecords[domain.InvoiceRecord](ctx, c, c.databaseIDs.Invoices, request)
}

func (c *Client) QueryCashflow(ctx context.Context, request QueryDatabaseRequest) (QueryDatabaseResponse[domain.CashflowRecord], error) {
	return queryRecords[domain.CashflowRecord](ctx, c, c.databaseIDs.Cashflow, request)
}

func (c *Client) UpsertTransaction(ctx context.Context, request UpsertPageRequest[domain.TransactionProperties]) (domain.TransactionRecord, error) {
	return upsertRecord[domain.TransactionRecord, domain.TransactionProperties](ctx, c, c.databaseIDs.Transactions, request)
}

func (c *Client) UpsertCustomer(ctx context.Context, request UpsertPageRequest[domain.CustomerProperties]) (domain.CustomerRecord, error) {
	return upsertRecord[domain.CustomerRecord, domain.CustomerProperties](ctx, c, c.databaseIDs.Customers, request)
}

func (c *Client) UpsertInvoice(ctx context.Context, request UpsertPageRequest[domain.InvoiceProperties]) (domain.InvoiceRecord, error) {
	return upsertRecord[domain.InvoiceRecord, domain.InvoiceProperties](ctx, c, c.databaseIDs.Invoices, request)
}

func (c *Client) UpsertCashflow(ctx context.Context, request UpsertPageRequest[domain.CashflowProperties]) (domain.CashflowRecord, error) {
	return upsertRecord[domain.CashflowRecord, domain.CashflowProperties](ctx, c, c.databaseIDs.Cashflow, request)
}

func queryRecords[T any](ctx context.Context, client *Client, databaseID string, request QueryDatabaseRequest) (QueryDatabaseResponse[T], error) {
	var response QueryDatabaseResponse[T]
	if databaseID == "" {
		return response, errDatabaseIDRequired
	}

	if err := client.queryDatabase(ctx, databaseID, request, &response); err != nil {
		return QueryDatabaseResponse[T]{}, err
	}

	return response, nil
}

func upsertRecord[T any, P any](ctx context.Context, client *Client, databaseID string, request UpsertPageRequest[P]) (T, error) {
	var response T
	if databaseID == "" {
		return response, errDatabaseIDRequired
	}

	if err := client.upsertPage(ctx, databaseID, request.PageID, request.Properties, &response); err != nil {
		return response, err
	}

	return response, nil
}
