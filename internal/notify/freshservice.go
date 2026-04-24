package notify

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
)

// FreshserviceClient sends alert tickets to Freshservice.
type FreshserviceClient struct {
	domain     string
	apiKey     string
	email      string
	httpClient *http.Client
}

type freshserviceTicket struct {
	Subject     string `json:"subject"`
	Description string `json:"description"`
	Email       string `json:"email"`
	Priority    int    `json:"priority"`
	Status      int    `json:"status"`
}

// NewFreshserviceClient creates a new FreshserviceClient.
// Returns an error if domain, apiKey, or email is empty.
func NewFreshserviceClient(domain, apiKey, email string) (*FreshserviceClient, error) {
	if domain == "" {
		return nil, fmt.Errorf("freshservice: domain must not be empty")
	}
	if apiKey == "" {
		return nil, fmt.Errorf("freshservice: API key must not be empty")
	}
	if email == "" {
		return nil, fmt.Errorf("freshservice: requester email must not be empty")
	}
	return &FreshserviceClient{
		domain:     domain,
		apiKey:     apiKey,
		email:      email,
		httpClient: &http.Client{},
	}, nil
}

// Send creates a new ticket in Freshservice with the given message.
func (c *FreshserviceClient) Send(message string) error {
	payload := freshserviceTicket{
		Subject:     "VaultPulse Alert",
		Description: message,
		Email:       c.email,
		Priority:    2,
		Status:      2,
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("freshservice: failed to marshal payload: %w", err)
	}
	url := fmt.Sprintf("https://%s.freshservice.com/api/v2/tickets", c.domain)
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("freshservice: failed to build request: %w", err)
	}
	creds := base64.StdEncoding.EncodeToString([]byte(c.apiKey + ":X"))
	req.Header.Set("Authorization", "Basic "+creds)
	req.Header.Set("Content-Type", "application/json")
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("freshservice: request failed: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("freshservice: unexpected status code: %d", resp.StatusCode)
	}
	return nil
}
