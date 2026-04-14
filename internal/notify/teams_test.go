package notify_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/yourusername/vaultpulse/internal/notify"
)

func TestNewTeamsClient_EmptyURL_ReturnsError(t *testing.T) {
	_, err := notify.NewTeamsClient("")
	if err == nil {
		t.Fatal("expected error for empty URL, got nil")
	}
}

func TestNewTeamsClient_ValidURL_ReturnsClient(t *testing.T) {
	client, err := notify.NewTeamsClient("https://outlook.office.com/webhook/test")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if client == nil {
		t.Fatal("expected non-nil client")
	}
}

func TestTeamsClient_Send_PostsCorrectPayload(t *testing.T) {
	var received map[string]string

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := json.NewDecoder(r.Body).Decode(&received); err != nil {
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client, _ := notify.NewTeamsClient(server.URL)
	if err := client.Send("secret expiring soon"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if received["text"] != "secret expiring soon" {
		t.Errorf("expected text %q, got %q", "secret expiring soon", received["text"])
	}
}

func TestTeamsClient_Send_NonOKStatus_ReturnsError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	client, _ := notify.NewTeamsClient(server.URL)
	if err := client.Send("alert"); err == nil {
		t.Fatal("expected error for non-2xx status, got nil")
	}
}

func TestTeamsClient_Send_UnreachableServer_ReturnsError(t *testing.T) {
	client, _ := notify.NewTeamsClient("http://127.0.0.1:0/webhook")
	if err := client.Send("alert"); err == nil {
		t.Fatal("expected error for unreachable server, got nil")
	}
}
