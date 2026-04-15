package notify

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// TestTwilioClient_MultiSender_Integration verifies that TwilioClient works
// correctly when composed inside a MultiSender alongside another sender.
func TestTwilioClient_MultiSender_Integration(t *testing.T) {
	var twilioReceived, webhookReceived string

	tsTwilio := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = r.ParseForm()
		twiloReceived = r.FormValue("Body")
		w.WriteHeader(http.StatusCreated)
	}))
	defer tsTwilio.Close()

	tsWebhook := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		webhookReceived = "received"
		w.WriteHeader(http.StatusOK)
	}))
	defer tsWebhook.Close()

	// Build Twilio client and redirect to test server.
	tc, err := NewTwilioClient("AC123", "token", "+15550001111", "+15559998888")
	if err != nil {
		t.Fatalf("NewTwilioClient: %v", err)
	}
	tc.httpClient = &http.Client{
		Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
			req.URL.Host = strings.TrimPrefix(tsTwilio.URL, "http://")
			req.URL.Scheme = "http"
			return http.DefaultTransport.RoundTrip(req)
		}),
	}

	wc, err := NewWebhookClient(tsWebhook.URL)
	if err != nil {
		t.Fatalf("NewWebhookClient: %v", err)
	}

	multi, err := NewMultiSender(tc, wc)
	if err != nil {
		t.Fatalf("NewMultiSender: %v", err)
	}

	const msg = "[vaultpulse] secret expires soon"
	if err := multi.Send(msg); err != nil {
		t.Fatalf("multi.Send: %v", err)
	}

	if twilioReceived != msg {
		t.Errorf("twilio body: want %q, got %q", msg, twilioReceived)
	}
	if webhookReceived != "received" {
		t.Error("webhook was not called")
	}
}
