package notify

import (
	"bytes"
	"fmt"
	"net/http"
	"net/url"
)

// ChatworkClient sends alert messages to a Chatwork room.
type ChatworkClient struct {
	token  string
	roomID string
	baseURL string
}

// NewChatworkClient creates a new ChatworkClient.
// token is the Chatwork API token and roomID is the target room.
func NewChatworkClient(token, roomID string) (*ChatworkClient, error) {
	if token == "" {
		return nil, fmt.Errorf("chatwork: token must not be empty")
	}
	if roomID == "" {
		return nil, fmt.Errorf("chatwork: roomID must not be empty")
	}
	return &ChatworkClient{
		token:   token,
		roomID:  roomID,
		baseURL: "https://api.chatwork.com/v2",
	}, nil
}

// Send posts a message to the configured Chatwork room.
func (c *ChatworkClient) Send(message string) error {
	endpoint := fmt.Sprintf("%s/rooms/%s/messages", c.baseURL, c.roomID)

	form := url.Values{}
	form.Set("body", message)

	req, err := http.NewRequest(http.MethodPost, endpoint, bytes.NewBufferString(form.Encode()))
	if err != nil {
		return fmt.Errorf("chatwork: failed to build request: %w", err)
	}
	req.Header.Set("X-ChatWorkToken", c.token)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("chatwork: request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("chatwork: unexpected status %d", resp.StatusCode)
	}
	return nil
}
