package notify

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

// MailgunClient sends alert messages via the Mailgun API.
type MailgunClient struct {
	domain  string
	apiKey  string
	from    string
	to      string
	httpClient *http.Client
}

// NewMailgunClient constructs a MailgunClient.
// domain, apiKey, from, and to are all required.
func NewMailgunClient(domain, apiKey, from, to string) (*MailgunClient, error) {
	if domain == "" {
		return nil, fmt.Errorf("mailgun: domain is required")
	}
	if apiKey == "" {
		return nil, fmt.Errorf("mailgun: api key is required")
	}
	if from == "" {
		return nil, fmt.Errorf("mailgun: from address is required")
	}
	if to == "" {
		return nil, fmt.Errorf("mailgun: to address is required")
	}
	return &MailgunClient{
		domain:     domain,
		apiKey:     apiKey,
		from:       from,
		to:         to,
		httpClient: &http.Client{},
	}, nil
}

// Send delivers the alert message via Mailgun's messages endpoint.
func (c *MailgunClient) Send(message string) error {
	endpoint := fmt.Sprintf("https://api.mailgun.net/v3/%s/messages", c.domain)

	form := url.Values{}
	form.Set("from", c.from)
	form.Set("to", c.to)
	form.Set("subject", "VaultPulse Alert")
	form.Set("text", message)

	req, err := http.NewRequest(http.MethodPost, endpoint, strings.NewReader(form.Encode()))
	if err != nil {
		return fmt.Errorf("mailgun: failed to build request: %w", err)
	}
	req.SetBasicAuth("api", c.apiKey)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("mailgun: request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("mailgun: unexpected status %d", resp.StatusCode)
	}
	return nil
}
