package notify

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewZulipClient_EmptyBaseURL_ReturnsError(t *testing.T) {
	_, err := NewZulipClient("", "bot@example.com", "key", "alerts", "vault")
	if err == nil {
		t.Fatal("expected error for empty baseURL")
	}
}

func TestNewZulipClient_EmptyEmail_ReturnsError(t *testing.T) {
	_, err := NewZulipClient("https://example.zulipchat.com", "", "key", "alerts", "vault")
	if err == nil {
		t.Fatal("expected error for empty email")
	}
}

func TestNewZulipClient_EmptyAPIKey_ReturnsError(t *testing.T) {
	_, err := NewZulipClient("https://example.zulipchat.com", "bot@example.com", "", "alerts", "vault")
	if err == nil {
		t.Fatal("expected error for empty apiKey")
	}
}

func TestNewZulipClient_EmptyStream_ReturnsError(t *testing.T) {
	_, err := NewZulipClient("https://example.zulipchat.com", "bot@example.com", "key", "", "vault")
	if err == nil {
		t.Fatal("expected error for empty stream")
	}
}

func TestNewZulipClient_ValidConfig_ReturnsClient(t *testing.T) {
	c, err := NewZulipClient("https://example.zulipchat.com", "bot@example.com", "key", "alerts", "vault")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if c == nil {
		t.Fatal("expected non-nil client")
	}
}

func TestZulipClient_Send_PostsCorrectPayload(t *testing.T) {
	var received zulipPayload
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := json.NewDecoder(r.Body).Decode(&received); err != nil {
			http.Error(w, "bad body", http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	c, _ := NewZulipClient(ts.URL, "bot@example.com", "key", "alerts", "vault-expiry")
	if err := c.Send("test alert"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if received.Content != "test alert" {
		t.Errorf("expected content 'test alert', got %q", received.Content)
	}
	if received.To != "alerts" {
		t.Errorf("expected stream 'alerts', got %q", received.To)
	}
	if received.Topic != "vault-expiry" {
		t.Errorf("expected topic 'vault-expiry', got %q", received.Topic)
	}
}

func TestZulipClient_Send_NonOKStatus_ReturnsError(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	}))
	defer ts.Close()

	c, _ := NewZulipClient(ts.URL, "bot@example.com", "key", "alerts", "vault")
	if err := c.Send("msg"); err == nil {
		t.Fatal("expected error for non-OK status")
	}
}

func TestZulipClient_Send_UnreachableServer_ReturnsError(t *testing.T) {
	c, _ := NewZulipClient("http://127.0.0.1:19999", "bot@example.com", "key", "alerts", "vault")
	if err := c.Send("msg"); err == nil {
		t.Fatal("expected error for unreachable server")
	}
}
