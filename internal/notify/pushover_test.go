package notify

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestNewPushoverClient_EmptyToken_ReturnsError(t *testing.T) {
	_, err := NewPushoverClient("", "user123")
	if err == nil {
		t.Fatal("expected error for empty api token")
	}
}

func TestNewPushoverClient_EmptyUserKey_ReturnsError(t *testing.T) {
	_, err := NewPushoverClient("token123", "")
	if err == nil {
		t.Fatal("expected error for empty user key")
	}
}

func TestNewPushoverClient_ValidConfig_ReturnsClient(t *testing.T) {
	c, err := NewPushoverClient("token123", "user123")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if c == nil {
		t.Fatal("expected non-nil client")
	}
}

func TestPushoverClient_Send_PostsCorrectPayload(t *testing.T) {
	var received pushoverPayload
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(body, &received)
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	c, _ := NewPushoverClient("tok", "usr")
	c.apiURL = ts.URL

	if err := c.Send("test alert"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if received.Token != "tok" {
		t.Errorf("expected token 'tok', got %q", received.Token)
	}
	if received.User != "usr" {
		t.Errorf("expected user 'usr', got %q", received.User)
	}
	if received.Message != "test alert" {
		t.Errorf("expected message 'test alert', got %q", received.Message)
	}
	if !strings.Contains(received.Title, "VaultPulse") {
		t.Errorf("expected title to contain 'VaultPulse', got %q", received.Title)
	}
}

func TestPushoverClient_Send_NonOKStatus_ReturnsError(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
	}))
	defer ts.Close()

	c, _ := NewPushoverClient("tok", "usr")
	c.apiURL = ts.URL

	if err := c.Send("msg"); err == nil {
		t.Fatal("expected error on non-OK status")
	}
}

func TestPushoverClient_Send_UnreachableServer_ReturnsError(t *testing.T) {
	c, _ := NewPushoverClient("tok", "usr")
	c.apiURL = "http://127.0.0.1:0"

	if err := c.Send("msg"); err == nil {
		t.Fatal("expected error for unreachable server")
	}
}
