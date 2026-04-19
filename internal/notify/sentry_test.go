package notify

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestNewSentryClient_EmptyDSN_ReturnsError(t *testing.T) {
	_, err := NewSentryClient("")
	if err == nil {
		t.Fatal("expected error for empty DSN")
	}
}

func TestNewSentryClient_ValidDSN_ReturnsClient(t *testing.T) {
	c, err := NewSentryClient("https://key@sentry.io/123")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if c == nil {
		t.Fatal("expected non-nil client")
	}
}

func TestSentryClient_Send_PostsCorrectPayload(t *testing.T) {
	var received string
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		buf := new(strings.Builder)
		_, _ = buf.ReadFrom(r.Body)
		received = buf.String()
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	c, _ := NewSentryClient(ts.URL)
	if err := c.Send("test alert"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(received, "test alert") {
		t.Errorf("expected payload to contain message, got: %s", received)
	}
}

func TestSentryClient_Send_NonOKStatus_ReturnsError(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer ts.Close()

	c, _ := NewSentryClient(ts.URL)
	if err := c.Send("alert"); err == nil {
		t.Fatal("expected error on non-OK status")
	}
}

func TestSentryClient_Send_UnreachableServer_ReturnsError(t *testing.T) {
	c, _ := NewSentryClient("http://127.0.0.1:0")
	if err := c.Send("alert"); err == nil {
		t.Fatal("expected error for unreachable server")
	}
}
