package notify

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewGoogleChatClient_EmptyURL_ReturnsError(t *testing.T) {
	_, err := NewGoogleChatClient("")
	if err == nil {
		t.Fatal("expected error for empty webhook URL, got nil")
	}
}

func TestNewGoogleChatClient_ValidURL_ReturnsClient(t *testing.T) {
	client, err := NewGoogleChatClient("https://chat.googleapis.com/v1/spaces/xxx")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if client == nil {
		t.Fatal("expected non-nil client")
	}
}

func TestGoogleChatClient_Send_PostsCorrectPayload(t *testing.T) {
	var received map[string]string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := json.NewDecoder(r.Body).Decode(&received); err != nil {
			t.Errorf("failed to decode request body: %v", err)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client, _ := NewGoogleChatClient(server.URL)
	if err := client.Send("vault alert: secret expiring soon"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if received["text"] != "vault alert: secret expiring soon" {
		t.Errorf("unexpected payload text: %q", received["text"])
	}
}

func TestGoogleChatClient_Send_NonOKStatus_ReturnsError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	client, _ := NewGoogleChatClient(server.URL)
	if err := client.Send("test"); err == nil {
		t.Fatal("expected error for non-2xx status, got nil")
	}
}

func TestGoogleChatClient_Send_UnreachableServer_ReturnsError(t *testing.T) {
	client, _ := NewGoogleChatClient("http://127.0.0.1:0/webhook")
	if err := client.Send("test"); err == nil {
		t.Fatal("expected error for unreachable server, got nil")
	}
}
