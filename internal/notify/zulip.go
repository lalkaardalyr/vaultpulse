package notify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

// ZulipClient sends alert messages to a Zulip stream via the Zulip REST API.
type ZulipClient struct {
	baseURL  string
	email    string
	apiKey   string
	stream   string
	topic    string
	httpClient *http.Client
}

type zulipPayload struct {
	Type    string `json:"type"`
	To      string `json:"to"`
	Topic   string `json:"topic"`
	Content string `json:"content"`
}

// NewZulipClient creates a new ZulipClient. baseURL should be the Zulip server
// URL (e.g. https://yourorg.zulipchat.com), email is the bot email, apiKey is
// the bot API key, stream is the target stream, and topic is the message topic.
func NewZulipClient(baseURL, email, apiKey, stream, topic string) (*ZulipClient, error) {
	if baseURL == "" {
		return nil, fmt.Errorf("zulip: baseURL must not be empty")
	}
	if email == "" {
		return nil, fmt.Errorf("zulip: email must not be empty")
	}
	if apiKey == "" {
		return nil, fmt.Errorf("zulip: apiKey must not be empty")
	}
	if stream == "" {
		return nil, fmt.Errorf("zulip: stream must not be empty")
	}
	return &ZulipClient{
		baseURL:    baseURL,
		email:      email,
		apiKey:     apiKey,
		stream:     stream,
		topic:      topic,
		httpClient: &http.Client{},
	}, nil
}

// Send posts the message to the configured Zulip stream.
func (z *ZulipClient) Send(message string) error {
	payload := zulipPayload{
		Type:    "stream",
		To:      z.stream,
		Topic:   z.topic,
		Content: message,
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("zulip: marshal payload: %w", err)
	}
	req, err := http.NewRequest(http.MethodPost, z.baseURL+"/api/v1/messages", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("zulip: create request: %w", err)
	}
	req.SetBasicAuth(z.email, z.apiKey)
	req.Header.Set("Content-Type", "application/json")
	resp, err := z.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("zulip: send request: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("zulip: unexpected status: %d", resp.StatusCode)
	}
	return nil
}
