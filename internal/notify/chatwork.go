package notify

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

const chatworkBaseURL = "https://api.chatwork.com/v2"

// ChatworkClient sends alert messages to a Chatwork room.
type ChatworkClient struct {
	token  string
	roomID string
	baseURL string
	httpClient *http.Client
}

// NewChatworkClient creates a new ChatworkClient.
// Returns an error if token or roomID is empty.
func NewChatworkClient(token, roomID string) (*ChatworkClient, error) {
	if token == "" {
		return nil, fmt.Errorf("chatwork: API token must not be empty")
	}
	if roomID == "" {
		return nil, fmt.Errorf("chatwork: room ID must not be empty")
	}
	return &ChatworkClient{
		token:      token,
		roomID:     roomID,
		baseURL:    chatworkBaseURL,
		httpClient: &http.Client{},
	}, nil
}

// Send posts a message to the configured Chatwork room.
func (c *ChatworkClient) Send(message string) error {
	endpoint := fmt.Sprintf("%s/rooms/%s/messages", c.baseURL, c.roomID)

	form := url.Values{}
	form.Set("body", message)

	req, err := http.NewRequest(http.MethodPost, endpoint, strings.NewReader(form.Encode()))
	if err != nil {
		return fmt.Errorf("chatwork: failed to build request: %w", err)
	}
	req.Header.Set("X-ChatWorkToken", c.token)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("chatwork: request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("chatwork: unexpected status code %d", resp.StatusCode)
	}
	return nil
}
