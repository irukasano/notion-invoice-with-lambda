package notion

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

const (
	defaultBaseURL       = "https://api.notion.com"
	defaultNotionVersion = "2022-06-28"
)

type HTTPDoer interface {
	Do(req *http.Request) (*http.Response, error)
}

type DatabaseIDs struct {
	Transactions string
	Customers    string
	Invoices     string
	Cashflow     string
}

type ClientConfig struct {
	BaseURL       string
	Token         string
	NotionVersion string
	DatabaseIDs   DatabaseIDs
}

type Client struct {
	baseURL       string
	token         string
	notionVersion string
	databaseIDs   DatabaseIDs
	httpClient    HTTPDoer
}

func NewClient(cfg ClientConfig, httpClient HTTPDoer) *Client {
	baseURL := strings.TrimRight(cfg.BaseURL, "/")
	if baseURL == "" {
		baseURL = defaultBaseURL
	}

	notionVersion := cfg.NotionVersion
	if notionVersion == "" {
		notionVersion = defaultNotionVersion
	}

	if httpClient == nil {
		httpClient = http.DefaultClient
	}

	return &Client{
		baseURL:       baseURL,
		token:         cfg.Token,
		notionVersion: notionVersion,
		databaseIDs:   cfg.DatabaseIDs,
		httpClient:    httpClient,
	}
}

func (c *Client) queryDatabase(ctx context.Context, databaseID string, payload QueryDatabaseRequest, out any) error {
	req, err := c.newJSONRequest(ctx, http.MethodPost, "/v1/databases/"+databaseID+"/query", payload)
	if err != nil {
		return err
	}

	return c.doJSON(req, out)
}

func (c *Client) upsertPage(ctx context.Context, databaseID, pageID string, properties any, out any) error {
	method := http.MethodPost
	path := "/v1/pages"
	body := map[string]any{
		"parent": map[string]string{
			"database_id": databaseID,
		},
		"properties": properties,
	}

	if pageID != "" {
		method = http.MethodPatch
		path = "/v1/pages/" + pageID
		body = map[string]any{
			"properties": properties,
		}
	}

	req, err := c.newJSONRequest(ctx, method, path, body)
	if err != nil {
		return err
	}

	return c.doJSON(req, out)
}

func (c *Client) newJSONRequest(ctx context.Context, method, path string, payload any) (*http.Request, error) {
	var body io.Reader
	if payload != nil {
		data, err := json.Marshal(payload)
		if err != nil {
			return nil, err
		}
		body = bytes.NewReader(data)
	}

	req, err := http.NewRequestWithContext(ctx, method, c.baseURL+path, body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+c.token)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Notion-Version", c.notionVersion)

	return req, nil
}

func (c *Client) doJSON(req *http.Request, out any) error {
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode >= http.StatusBadRequest {
		return classifyHTTPError(req, resp.StatusCode, body)
	}

	if out == nil || len(body) == 0 {
		return nil
	}

	if err := json.Unmarshal(body, out); err != nil {
		return fmt.Errorf("decode notion response: %w", err)
	}

	return nil
}
