package sendgrid

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/mail"
	"strings"
)

const defaultBaseURL = "https://api.sendgrid.com/v3"

type HTTPDoer interface {
	Do(req *http.Request) (*http.Response, error)
}

type ClientConfig struct {
	BaseURL   string
	APIKey    string
	FromEmail string
	FromName  string
}

type Attachment struct {
	Filename    string
	ContentType string
	Content     []byte
}

type SendMailInput struct {
	ToEmail     string
	ToName      string
	Subject     string
	PlainText   string
	Attachments []Attachment
}

type Client struct {
	baseURL    string
	apiKey     string
	fromEmail  string
	fromName   string
	httpClient HTTPDoer
}

func NewClient(cfg ClientConfig, httpClient HTTPDoer) (*Client, error) {
	baseURL := strings.TrimRight(strings.TrimSpace(cfg.BaseURL), "/")
	if baseURL == "" {
		baseURL = defaultBaseURL
	}
	if strings.TrimSpace(cfg.APIKey) == "" {
		return nil, errors.New("sendgrid api key is required")
	}
	if _, err := parseEmail("from_email", cfg.FromEmail); err != nil {
		return nil, err
	}
	if strings.TrimSpace(cfg.FromName) == "" {
		return nil, errors.New("sendgrid from name is required")
	}
	if httpClient == nil {
		httpClient = http.DefaultClient
	}

	return &Client{
		baseURL:    baseURL,
		apiKey:     strings.TrimSpace(cfg.APIKey),
		fromEmail:  strings.TrimSpace(cfg.FromEmail),
		fromName:   strings.TrimSpace(cfg.FromName),
		httpClient: httpClient,
	}, nil
}

func (c *Client) SendMail(ctx context.Context, input SendMailInput) error {
	toEmail, err := parseEmail("to_email", input.ToEmail)
	if err != nil {
		return err
	}

	payload := mailSendPayload{
		From: emailAddress{
			Email: c.fromEmail,
			Name:  c.fromName,
		},
		Personalizations: []personalization{{
			To: []emailAddress{{
				Email: toEmail,
				Name:  strings.TrimSpace(input.ToName),
			}},
		}},
		Subject: strings.TrimSpace(input.Subject),
		Content: []content{{
			Type:  "text/plain",
			Value: input.PlainText,
		}},
		Attachments: make([]attachment, 0, len(input.Attachments)),
	}

	for _, item := range input.Attachments {
		payload.Attachments = append(payload.Attachments, attachment{
			Content:     base64.StdEncoding.EncodeToString(item.Content),
			Type:        strings.TrimSpace(item.ContentType),
			Filename:    strings.TrimSpace(item.Filename),
			Disposition: "attachment",
		})
	}

	req, err := c.newJSONRequest(ctx, http.MethodPost, "/mail/send", payload)
	if err != nil {
		return err
	}

	return c.do(req)
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

	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	req.Header.Set("Content-Type", "application/json")

	return req, nil
}

func (c *Client) do(req *http.Request) error {
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

	return nil
}

func parseEmail(field, value string) (string, error) {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return "", &ValidationError{
			Field:   field,
			Message: "is required",
		}
	}
	addr, err := mail.ParseAddress(trimmed)
	if err != nil {
		return "", &ValidationError{
			Field:   field,
			Message: "must be a valid email address",
		}
	}
	return addr.Address, nil
}

type mailSendPayload struct {
	Personalizations []personalization `json:"personalizations"`
	From             emailAddress      `json:"from"`
	Subject          string            `json:"subject"`
	Content          []content         `json:"content"`
	Attachments      []attachment      `json:"attachments,omitempty"`
}

type personalization struct {
	To []emailAddress `json:"to"`
}

type emailAddress struct {
	Email string `json:"email"`
	Name  string `json:"name,omitempty"`
}

type content struct {
	Type  string `json:"type"`
	Value string `json:"value"`
}

type attachment struct {
	Content     string `json:"content"`
	Type        string `json:"type,omitempty"`
	Filename    string `json:"filename,omitempty"`
	Disposition string `json:"disposition,omitempty"`
}
