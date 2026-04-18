package notify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

// ZendeskClient sends alert notifications as Zendesk tickets.
type ZendeskClient struct {
	subdomain string
	email     string
	apiToken  string
	httpClient *http.Client
}

type zendeskPayload struct {
	Ticket zendeskTicket `json:"ticket"`
}

type zendeskTicket struct {
	Subject string          `json:"subject"`
	Comment zendeskComment  `json:"comment"`
	Priority string         `json:"priority"`
}

type zendeskComment struct {
	Body string `json:"body"`
}

// NewZendeskClient creates a new ZendeskClient.
func NewZendeskClient(subdomain, email, apiToken string) (*ZendeskClient, error) {
	if subdomain == "" {
		return nil, fmt.Errorf("zendesk: subdomain must not be empty")
	}
	if email == "" {
		return nil, fmt.Errorf("zendesk: email must not be empty")
	}
	if apiToken == "" {
		return nil, fmt.Errorf("zendesk: api token must not be empty")
	}
	return &ZendeskClient{
		subdomain:  subdomain,
		email:      email,
		apiToken:   apiToken,
		httpClient: &http.Client{},
	}, nil
}

// Send posts an alert as a Zendesk ticket.
func (c *ZendeskClient) Send(msg string) error {
	url := fmt.Sprintf("https://%s.zendesk.com/api/v2/tickets.json", c.subdomain)
	payload := zendeskPayload{
		Ticket: zendeskTicket{
			Subject:  "VaultPulse Alert",
			Comment:  zendeskComment{Body: msg},
			Priority: "high",
		},
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("zendesk: marshal payload: %w", err)
	}
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("zendesk: create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth(c.email+"/token", c.apiToken)
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("zendesk: send request: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("zendesk: unexpected status %d", resp.StatusCode)
	}
	return nil
}
