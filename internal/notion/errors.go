package notion

import (
	"fmt"
	"net/http"
	"strings"
)

type HTTPError struct {
	Method     string
	Path       string
	StatusCode int
	Body       string
}

func (e *HTTPError) Error() string {
	return fmt.Sprintf(
		"notion api %s %s returned %d: %s",
		e.Method,
		e.Path,
		e.StatusCode,
		e.Body,
	)
}

type RateLimitError struct {
	HTTPError
}

type ServerError struct {
	HTTPError
}

type ClientError struct {
	HTTPError
}

func classifyHTTPError(req *http.Request, statusCode int, body []byte) error {
	httpErr := HTTPError{
		Method:     req.Method,
		Path:       req.URL.Path,
		StatusCode: statusCode,
		Body:       strings.TrimSpace(string(body)),
	}

	switch {
	case statusCode == http.StatusTooManyRequests:
		return &RateLimitError{HTTPError: httpErr}
	case statusCode >= http.StatusInternalServerError:
		return &ServerError{HTTPError: httpErr}
	default:
		return &ClientError{HTTPError: httpErr}
	}
}
