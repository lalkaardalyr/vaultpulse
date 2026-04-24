package notify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

const pagerTreeAPIURL = "https://api.pagertree.com/integration/"

// PagerTreeClient sends alerts to PagerTree via its integration API.
type PagerTreeClient struct {
	integrationID string
	endpoint       string
	httpClient    *http.Client
}

// NewPagerTreeClient creates a new PagerTreeClient.
// Returns an error if integrationID is empty.
func NewPagerTreeClient(integrationID string) (*PagerTreeClient, error) {
	if integrationID == "" {
		return nil, fmt.Errorf("pagertree: integration ID must not be empty")
	}
	return &PagerTreeClient{
		integrationID: integrationID,
		endpoint:       pagerTreeAPIURL + integrationID,
		httpClient:    &http.Client{},
	}, nil
}

type pagerTreePayload struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

// Send delivers msg to the configured PagerTree integration.
func (c *PagerTreeClient) Send(msg string) error {
	p := pagerTreePayload{
		Title:       "VaultPulse Alert",
		Description: msg,
	}
	body, err := json.Marshal(p)
	if err != nil {
		return fmt.Errorf("pagertree: marshal payload: %w", err)
	}

	resp, err := c.httpClient.Post(c.endpoint, "application/json", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("pagertree: post: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("pagertree: unexpected status %d", resp.StatusCode)
	}
	return nil
}
