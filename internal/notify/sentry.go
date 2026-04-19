package notify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

// SentryClient sends alert notifications to a Sentry DSN endpoint.
type SentryClient struct {
	dsn    string
	httpClient *http.Client
}

type sentryPayload struct {
	Message string            `json:"message"`
	Level   string            `json:"level"`
	Extra   map[string]string `json:"extra,omitempty"`
}

// NewSentryClient creates a new SentryClient.
// dsn must not be empty.
func NewSentryClient(dsn string) (*SentryClient, error) {
	if dsn == "" {
		return nil, fmt.Errorf("sentry: DSN must not be empty")
	}
	return &SentryClient{
		dsn:        dsn,
		httpClient: &http.Client{},
	}, nil
}

// Send posts an alert message to Sentry.
func (c *SentryClient) Send(message string) error {
	payload := sentryPayload{
		Message: message,
		Level:   "error",
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("sentry: failed to marshal payload: %w", err)
	}
	resp, err := c.httpClient.Post(c.dsn, "application/json", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("sentry: request failed: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("sentry: unexpected status code: %d", resp.StatusCode)
	}
	return nil
}
