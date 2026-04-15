package notify

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

const twilioAPIBase = "https://api.twilio.com/2010-04-01/Accounts"

// TwilioClient sends SMS alerts via the Twilio REST API.
type TwilioClient struct {
	accountSID string
	authToken  string
	from       string
	to         string
	httpClient *http.Client
}

// NewTwilioClient constructs a TwilioClient. accountSID, authToken, from, and
// to are all required.
func NewTwilioClient(accountSID, authToken, from, to string) (*TwilioClient, error) {
	if accountSID == "" {
		return nil, fmt.Errorf("twilio: accountSID must not be empty")
	}
	if authToken == "" {
		return nil, fmt.Errorf("twilio: authToken must not be empty")
	}
	if from == "" {
		return nil, fmt.Errorf("twilio: from number must not be empty")
	}
	if to == "" {
		return nil, fmt.Errorf("twilio: to number must not be empty")
	}
	return &TwilioClient{
		accountSID: accountSID,
		authToken:  authToken,
		from:       from,
		to:         to,
		httpClient: &http.Client{},
	}, nil
}

// Send delivers the alert message as an SMS via Twilio.
func (c *TwilioClient) Send(message string) error {
	endpoint := fmt.Sprintf("%s/%s/Messages.json", twilioAPIBase, c.accountSID)

	form := url.Values{}
	form.Set("From", c.from)
	form.Set("To", c.to)
	form.Set("Body", message)

	req, err := http.NewRequest(http.MethodPost, endpoint, strings.NewReader(form.Encode()))
	if err != nil {
		return fmt.Errorf("twilio: failed to create request: %w", err)
	}
	req.SetBasicAuth(c.accountSID, c.authToken)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("twilio: request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		var errBody struct {
			Message string `json:"message"`
		}
		_ = json.NewDecoder(resp.Body).Decode(&errBody)
		return fmt.Errorf("twilio: unexpected status %d: %s", resp.StatusCode, errBody.Message)
	}
	return nil
}
