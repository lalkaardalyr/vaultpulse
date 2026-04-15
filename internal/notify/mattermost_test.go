package notify_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/yourusername/vaultpulse/internal/notify"
)

func TestNewMattermostClient_EmptyURL_ReturnsError(t *testing.T) {
	_, err := notify.NewMattermostClient("")
	if err == nil {
		t.Fatal("expected error for empty webhook URL, got nil")
	}
}

func TestNewMattermostClient_ValidURL_ReturnsClient(t *testing.T) {
	client, err := notify.NewMattermostClient("https://mattermost.example.com/hooks/abc123")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if client == nil {
		t.Fatal("expected non-nil client")
	}
}

func TestMattermostClient_Send_PostsCorrectPayload(t *testing.T) {
	var received map[string]interface{}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := json.NewDecoder(r.Body).Decode(&received); err != nil {
			t.Errorf("failed to decode request body: %v", err)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client, err := notify.NewMattermostClient(server.URL)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if err := client.Send("secret expiring soon"); err != nil {
		t.Fatalf("Send returned error: %v", err)
	}

	if received["text"] != "secret expiring soon" {
		t.Errorf("expected text %q, got %q", "secret expiring soon", received["text"])
	}
	if received["username"] != "VaultPulse" {
		t.Errorf("expected username %q, got %q", "VaultPulse", received["username"])
	}
}

func TestMattermostClient_Send_NonOKStatus_ReturnsError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	client, err := notify.NewMattermostClient(server.URL)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if err := client.Send("test"); err == nil {
		t.Fatal("expected error for non-OK status, got nil")
	}
}

func TestMattermostClient_Send_UnreachableServer_ReturnsError(t *testing.T) {
	client, err := notify.NewMattermostClient("http://127.0.0.1:0/hooks/invalid")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if err := client.Send("test"); err == nil {
		t.Fatal("expected error for unreachable server, got nil")
	}
}
