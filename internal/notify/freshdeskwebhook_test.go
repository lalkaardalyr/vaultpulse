package notify_test

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/yourusername/vaultpulse/internal/notify"
)

func TestNewFreshdeskWebhookClient_EmptyURL_ReturnsError(t *testing.T) {
	_, err := notify.NewFreshdeskWebhookClient("")
	if err == nil {
		t.Fatal("expected error for empty URL, got nil")
	}
}

func TestNewFreshdeskWebhookClient_ValidURL_ReturnsClient(t *testing.T) {
	client, err := notify.NewFreshdeskWebhookClient("https://example.freshdesk.com/webhook")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if client == nil {
		t.Fatal("expected non-nil client")
	}
}

func TestFreshdeskWebhookClient_Send_PostsCorrectPayload(t *testing.T) {
	var received map[string]interface{}
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(body, &received)
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	client, err := notify.NewFreshdeskWebhookClient(ts.URL)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if err := client.Send("vault secret expiring soon"); err != nil {
		t.Fatalf("Send returned error: %v", err)
	}

	if received["text"] != "vault secret expiring soon" {
		t.Errorf("expected text payload, got %v", received)
	}
}

func TestFreshdeskWebhookClient_Send_NonOKStatus_ReturnsError(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer ts.Close()

	client, _ := notify.NewFreshdeskWebhookClient(ts.URL)
	if err := client.Send("test"); err == nil {
		t.Fatal("expected error for non-OK status, got nil")
	}
}

func TestFreshdeskWebhookClient_Send_UnreachableServer_ReturnsError(t *testing.T) {
	client, _ := notify.NewFreshdeskWebhookClient("http://127.0.0.1:19999")
	if err := client.Send("test"); err == nil {
		t.Fatal("expected error for unreachable server, got nil")
	}
}
