package notify

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// FreshdeskClient sends alert notifications by creating tickets in Freshdesk.
type FreshdeskClient struct {
	domain   string
	apiKey   string
	email    string
	httpClient *http.Client
}

type freshdeskTicket struct {
	Subject     string `json:"subject"`
	Description string `json:"description"`
	Email       string `json:"email"`
	Priority    int    `json:"priority"`
	Status      int    `json:"status"`
}

// NewFreshdeskClient creates a new FreshdeskClient.
func NewFreshdeskClient(domain, apiKey, email string) (*FreshdeskClient, error) {
	if domain == "" {
		return nil, fmt.Errorf("freshdesk: domain must not be empty")
	}
	if apiKey == "" {
		return nil, fmt.Errorf("freshdesk: api key must not be empty")
	}
	if email == "" {
		return nil, fmt.Errorf("freshdesk: email must not be empty")
	}
	return &FreshdeskClient{
		domain:     domain,
		apiKey:     apiKey,
		email:      email,
		httpClient: &http.Client{},
	}, nil
}

// Send creates a Freshdesk support ticket with the alert message.
func (c *FreshdeskClient) Send(ctx context.Context, message string) error {
	ticker := freshdeskTicket{
		Subject:     "VaultPulse Alert",
		Description: message,
		Email:       c.email,
		Priority:    2,
		Status:      2,
	}
	body, err := json.Marshal(ticker)
	if err != nil {
		return fmt.Errorf("freshdesk: marshal error: %w", err)
	}
	url := fmt.Sprintf("https://%s.freshdesk.com/api/v2/tickets", c.domain)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("freshdesk: request error: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth(c.apiKey, "X")
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("freshdesk: send error: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("freshdesk: unexpected status %d", resp.StatusCode)
	}
	return nil
}
