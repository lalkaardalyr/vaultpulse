package notify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

const telegramAPIBase = "https://api.telegram.org/bot"

// TelegramClient sends alert messages to a Telegram chat via the Bot API.
type TelegramClient struct {
	botToken string
	chatID   string
	endpoint string
	httpClient *http.Client
}

type telegramPayload struct {
	ChatID    string `json:"chat_id"`
	Text      string `json:"text"`
	ParseMode string `json:"parse_mode"`
}

// NewTelegramClient creates a new TelegramClient.
// botToken must be the token provided by BotFather.
// chatID is the target chat or channel ID.
func NewTelegramClient(botToken, chatID string) (*TelegramClient, error) {
	if botToken == "" {
		return nil, fmt.Errorf("telegram: bot token must not be empty")
	}
	if chatID == "" {
		return nil, fmt.Errorf("telegram: chat ID must not be empty")
	}
	return &TelegramClient{
		botToken:   botToken,
		chatID:     chatID,
		endpoint:   telegramAPIBase + botToken + "/sendMessage",
		httpClient: &http.Client{},
	}, nil
}

// Send posts the message to the configured Telegram chat.
func (c *TelegramClient) Send(message string) error {
	payload := telegramPayload{
		ChatID:    c.chatID,
		Text:      message,
		ParseMode: "Markdown",
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("telegram: failed to marshal payload: %w", err)
	}
	resp, err := c.httpClient.Post(c.endpoint, "application/json", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("telegram: request failed: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("telegram: unexpected status code %d", resp.StatusCode)
	}
	return nil
}
