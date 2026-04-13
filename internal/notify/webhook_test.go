package notify_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/yourusername/vaultpulse/internal/notify"
)

func TestNewWebhookClient_EmptyURL_ReturnsError(t *testing.T) {
	_, err := notify.NewWebhookClient("")
	if err == nil {
		t.Fatal("expected error for empty URL, got nil")
	}
}

func TestNewWebhookClient_ValidURL_ReturnsClient(t *testing.T) {
	client, err := notify.NewWebhookClient("https://example.com/hook")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if client == nil {
		t.Fatal("expected non-nil client")
	}
}

func TestWebhookClient_Send_PostsCorrectPayload(t *testing.T) {
	var received map[string]string

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if ct := r.Header.Get("Content-Type"); ct != "application/json" {
			t.Errorf("expected application/json, got %s", ct)
		}
		if err := json.NewDecoder(r.Body).Decode(&received); err != nil {
			t.Errorf("failed to decode body: %v", err)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	client, err := notify.NewWebhookClient(ts.URL)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	msg := notify.Message{Body: "secret expires soon", Severity: "critical"}
	if err := client.Send(msg); err != nil {
		t.Fatalf("Send returned error: %v", err)
	}

	if received["message"] != "secret expires soon" {
		t.Errorf("unexpected message: %s", received["message"])
	}
	if received["severity"] != "critical" {
		t.Errorf("unexpected severity: %s", received["severity"])
	}
}

func TestWebhookClient_Send_NonOKStatus_ReturnsError(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer ts.Close()

	client, _ := notify.NewWebhookClient(ts.URL)
	msg := notify.Message{Body: "test", Severity: "warning"}
	if err := client.Send(msg); err == nil {
		t.Fatal("expected error for non-2xx status, got nil")
	}
}

func TestWebhookClient_Send_UnreachableServer_ReturnsError(t *testing.T) {
	client, _ := notify.NewWebhookClient("http://127.0.0.1:19999/hook")
	msg := notify.Message{Body: "test", Severity: "info"}
	if err := client.Send(msg); err == nil {
		t.Fatal("expected error for unreachable server, got nil")
	}
}
