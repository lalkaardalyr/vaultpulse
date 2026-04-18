package notify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

// WhatsAppClient sends alerts via the WhatsApp Business Cloud API.
type WhatsAppClient struct {
	token       string
	phoneID     string
	recipient   string
	httpClient  *http.Client
	endpoint    string
}

// NewWhatsAppClient creates a new WhatsAppClient.
func NewWhatsAppClient(token, phoneID, recipient string) (*WhatsAppClient, error) {
	if token == "" {
		return nil, fmt.Errorf("whatsapp: token must not be empty")
	}
	if phoneID == "" {
		return nil, fmt.Errorf("whatsapp: phone ID must not be empty")
	}
	if recipient == "" {
		return nil, fmt.Errorf("whatsapp: recipient must not be empty")
	}
	return &WhatsAppClient{
		token:      token,
		phoneID:    phoneID,
		recipient:  recipient,
		httpClient: &http.Client{},
		endpoint:   "https://graph.facebook.com/v17.0",
	}, nil
}

// Send posts a text message to the WhatsApp Business API.
func (c *WhatsAppClient) Send(msg string) error {
	url := fmt.Sprintf("%s/%s/messages", c.endpoint, c.phoneID)
	payload := map[string]interface{}{
		"messaging_product": "whatsapp",
		"to":                c.recipient,
		"type":              "text",
		"text":              map[string]string{"body": msg},
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("whatsapp: marshal: %w", err)
	}
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("whatsapp: new request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.token)
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("whatsapp: do request: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("whatsapp: unexpected status: %d", resp.StatusCode)
	}
	return nil
}
