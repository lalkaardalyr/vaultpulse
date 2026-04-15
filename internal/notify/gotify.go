package notify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

// GotifyClient sends alert messages to a Gotify server.
type GotifyClient struct {
	url      string
	token    string
	priority int
	httpClient *http.Client
}

type gotifyPayload struct {
	Title    string `json:"title"`
	Message  string `json:"message"`
	Priority int    `json:"priority"`
}

// NewGotifyClient creates a new GotifyClient.
// url is the base URL of the Gotify server (e.g. https://gotify.example.com).
// token is the application token used for authentication.
// priority sets the message priority (0–10); values >=5 are treated as high priority.
func NewGotifyClient(url, token string, priority int) (*GotifyClient, error) {
	if url == "" {
		return nil, fmt.Errorf("gotify: server URL must not be empty")
	}
	if token == "" {
		return nil, fmt.Errorf("gotify: application token must not be empty")
	}
	if priority < 0 {
		priority = 0
	}
	return &GotifyClient{
		url:        url,
		token:      token,
		priority:   priority,
		httpClient: &http.Client{},
	}, nil
}

// Send posts an alert message to the Gotify server.
func (c *GotifyClient) Send(message string) error {
	payload := gotifyPayload{
		Title:    "VaultPulse Alert",
		Message:  message,
		Priority: c.priority,
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("gotify: failed to marshal payload: %w", err)
	}
	endpoint := fmt.Sprintf("%s/message?token=%s", c.url, c.token)
	resp, err := c.httpClient.Post(endpoint, "application/json", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("gotify: request failed: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("gotify: unexpected status code %d", resp.StatusCode)
	}
	return nil
}
