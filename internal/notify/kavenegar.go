package notify

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

const kavenegarBaseURL = "https://api.kavenegar.com/v1"

// KavenegarClient sends SMS alerts via the Kavenegar API.
type KavenegarClient struct {
	apiKey  string
	sender  string
	receptor string
	httpClient *http.Client
}

// NewKavenegarClient creates a new KavenegarClient.
func NewKavenegarClient(apiKey, sender, receptor string) (*KavenegarClient, error) {
	if apiKey == "" {
		return nil, errors.New("kavenegar: api key must not be empty")
	}
	if sender == "" {
		return nil, errors.New("kavenegar: sender must not be empty")
	}
	if receptor == "" {
		return nil, errors.New("kavenegar: receptor must not be empty")
	}
	return &KavenegarClient{
		apiKey:     apiKey,
		sender:     sender,
		receptor:   receptor,
		httpClient: &http.Client{},
	}, nil
}

// Send delivers the message as an SMS via Kavenegar.
func (c *KavenegarClient) Send(message string) error {
	endpoint := fmt.Sprintf("%s/%s/sms/send.json", kavenegarBaseURL, c.apiKey)

	params := url.Values{}
	params.Set("sender", c.sender)
	params.Set("receptor", c.receptor)
	params.Set("message", message)

	resp, err := c.httpClient.PostForm(endpoint, params)
	if err != nil {
		return fmt.Errorf("kavenegar: request failed: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()
	_, _ = io.Copy(io.Discard, resp.Body)

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("kavenegar: unexpected status %d", resp.StatusCode)
	}
	return nil
}
