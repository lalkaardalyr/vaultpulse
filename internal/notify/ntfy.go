package notify

import (
	"bytes"
	"fmt"
	"net/http"
)

// NtfyClient sends notifications via ntfy.sh or a self-hosted ntfy server.
type NtfyClient struct {
	serverURL string
	topic     string
	httpClient *http.Client
}

// NewNtfyClient creates a new NtfyClient.
// serverURL should be the base URL (e.g. https://ntfy.sh) and topic the target topic.
func NewNtfyClient(serverURL, topic string) (*NtfyClient, error) {
	if serverURL == "" {
		return nil, fmt.Errorf("ntfy: server URL must not be empty")
	}
	if topic == "" {
		return nil, fmt.Errorf("ntfy: topic must not be empty")
	}
	return &NtfyClient{
		serverURL:  serverURL,
		topic:      topic,
		httpClient: &http.Client{},
	}, nil
}

// Send posts a plain-text message to the configured ntfy topic.
func (c *NtfyClient) Send(message string) error {
	url := fmt.Sprintf("%s/%s", c.serverURL, c.topic)
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBufferString(message))
	if err != nil {
		return fmt.Errorf("ntfy: failed to build request: %w", err)
	}
	req.Header.Set("Content-Type", "text/plain")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("ntfy: request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("ntfy: unexpected status %d", resp.StatusCode)
	}
	return nil
}
