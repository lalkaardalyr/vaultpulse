package notify

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewSpikeClient_EmptyURL_ReturnsError(t *testing.T) {
	_, err := NewSpikeClient("")
	if err == nil {
		t.Fatal("expected error for empty URL")
	}
}

func TestNewSpikeClient_ValidURL_ReturnsClient(t *testing.T) {
	c, err := NewSpikeClient("https://hooks.spike.sh/abc123")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if c == nil {
		t.Fatal("expected non-nil client")
	}
}

func TestSpikeClient_Send_PostsCorrectPayload(t *testing.T) {
	var received map[string]string
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(body, &received)
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	c, _ := NewSpikeClient(ts.URL)
	if err := c.Send("vault secret expiring soon"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if received["message"] != "vault secret expiring soon" {
		t.Errorf("unexpected message: %s", received["message"])
	}
	if received["severity"] != "critical" {
		t.Errorf("unexpected severity: %s", received["severity"])
	}
}

func TestSpikeClient_Send_NonOKStatus_ReturnsError(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer ts.Close()

	c, _ := NewSpikeClient(ts.URL)
	if err := c.Send("test"); err == nil {
		t.Fatal("expected error for non-OK status")
	}
}

func TestSpikeClient_Send_UnreachableServer_ReturnsError(t *testing.T) {
	c, _ := NewSpikeClient("http://127.0.0.1:0/spike")
	if err := c.Send("test"); err == nil {
		t.Fatal("expected error for unreachable server")
	}
}
