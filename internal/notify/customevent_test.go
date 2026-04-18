package notify

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewCustomEventClient_EmptyURL_ReturnsError(t *testing.T) {
	_, err := NewCustomEventClient("", "key", "vaultpulse")
	if err == nil {
		t.Fatal("expected error for empty URL")
	}
}

func TestNewCustomEventClient_EmptyAPIKey_ReturnsError(t *testing.T) {
	_, err := NewCustomEventClient("http://example.com", "", "vaultpulse")
	if err == nil {
		t.Fatal("expected error for empty api key")
	}
}

func TestNewCustomEventClient_EmptySource_ReturnsError(t *testing.T) {
	_, err := NewCustomEventClient("http://example.com", "key", "")
	if err == nil {
		t.Fatal("expected error for empty source")
	}
}

func TestNewCustomEventClient_ValidConfig_ReturnsClient(t *testing.T) {
	c, err := NewCustomEventClient("http://example.com", "key", "vaultpulse")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if c == nil {
		t.Fatal("expected non-nil client")
	}
}

func TestCustomEventClient_Send_PostsCorrectPayload(t *testing.T) {
	var received map[string]string
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := json.NewDecoder(r.Body).Decode(&received); err != nil {
			t.Errorf("failed to decode body: %v", err)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	c, _ := NewCustomEventClient(ts.URL, "testkey", "vaultpulse")
	if err := c.Send("secret expiring soon"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if received["source"] != "vaultpulse" {
		t.Errorf("expected source vaultpulse, got %s", received["source"])
	}
	if received["message"] != "secret expiring soon" {
		t.Errorf("unexpected message: %s", received["message"])
	}
}

func TestCustomEventClient_Send_NonOKStatus_ReturnsError(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer ts.Close()

	c, _ := NewCustomEventClient(ts.URL, "testkey", "vaultpulse")
	if err := c.Send("test"); err == nil {
		t.Fatal("expected error on non-OK status")
	}
}

func TestCustomEventClient_Send_UnreachableServer_ReturnsError(t *testing.T) {
	c, _ := NewCustomEventClient("http://127.0.0.1:0", "testkey", "vaultpulse")
	if err := c.Send("test"); err == nil {
		t.Fatal("expected error for unreachable server")
	}
}
