package notify

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestNewTwilioClient_EmptyAccountSID_ReturnsError(t *testing.T) {
	_, err := NewTwilioClient("", "token", "+1", "+2")
	if err == nil {
		t.Fatal("expected error for empty accountSID")
	}
}

func TestNewTwilioClient_EmptyAuthToken_ReturnsError(t *testing.T) {
	_, err := NewTwilioClient("AC123", "", "+1", "+2")
	if err == nil {
		t.Fatal("expected error for empty authToken")
	}
}

func TestNewTwilioClient_EmptyFrom_ReturnsError(t *testing.T) {
	_, err := NewTwilioClient("AC123", "token", "", "+2")
	if err == nil {
		t.Fatal("expected error for empty from")
	}
}

func TestNewTwilioClient_EmptyTo_ReturnsError(t *testing.T) {
	_, err := NewTwilioClient("AC123", "token", "+1", "")
	if err == nil {
		t.Fatal("expected error for empty to")
	}
}

func TestNewTwilioClient_ValidConfig_ReturnsClient(t *testing.T) {
	c, err := NewTwilioClient("AC123", "token", "+15550001111", "+15559998888")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if c == nil {
		t.Fatal("expected non-nil client")
	}
}

func TestTwilioClient_Send_PostsCorrectPayload(t *testing.T) {
	var capturedBody string
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = r.ParseForm()
		capturedBody = r.FormValue("Body")
		w.WriteHeader(http.StatusCreated)
	}))
	defer ts.Close()

	c, _ := NewTwilioClient("AC123", "token", "+15550001111", "+15559998888")
	// Override the base URL by pointing the httpClient at our test server.
	// We patch the endpoint indirectly by replacing the httpClient transport.
	c.httpClient = &http.Client{
		Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
			req.URL.Host = strings.TrimPrefix(ts.URL, "http://")
			req.URL.Scheme = "http"
			return http.DefaultTransport.RoundTrip(req)
		}),
	}

	if err := c.Send("test alert"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if capturedBody != "test alert" {
		t.Errorf("expected body 'test alert', got %q", capturedBody)
	}
}

func TestTwilioClient_Send_NonOKStatus_ReturnsError(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	}))
	defer ts.Close()

	c, _ := NewTwilioClient("AC123", "token", "+1", "+2")
	c.httpClient = ts.Client()
	// Point at the test server directly.
	c.httpClient = &http.Client{
		Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
			req.URL.Host = strings.TrimPrefix(ts.URL, "http://")
			req.URL.Scheme = "http"
			return http.DefaultTransport.RoundTrip(req)
		}),
	}

	if err := c.Send("msg"); err == nil {
		t.Fatal("expected error on non-2xx status")
	}
}
