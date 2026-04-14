package notify_test

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/yourusername/vaultpulse/internal/notify"
)

func TestNewVictorOpsClient_EmptyURL_ReturnsError(t *testing.T) {
	_, err := notify.NewVictorOpsClient("")
	if err == nil {
		t.Fatal("expected error for empty URL, got nil")
	}
}

func TestNewVictorOpsClient_ValidURL_ReturnsClient(t *testing.T) {
	client, err := notify.NewVictorOpsClient("https://alert.victorops.com/integrations/generic/webhook")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if client == nil {
		t.Fatal("expected non-nil client")
	}
}

func TestVictorOpsClient_Send_PostsCorrectPayload(t *testing.T) {
	var received map[string]interface{}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(body, &received)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client, err := notify.NewVictorOpsClient(server.URL)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	const msg = "secret/db expires in 2 days"
	if err := client.Send(msg); err != nil {
		t.Fatalf("Send returned error: %v", err)
	}

	if received["state_message"] != msg {
		t.Errorf("expected state_message %q, got %q", msg, received["state_message"])
	}
	if received["message_type"] != "CRITICAL" {
		t.Errorf("expected message_type CRITICAL, got %q", received["message_type"])
	}
	if received["monitoring_tool"] != "vaultpulse" {
		t.Errorf("expected monitoring_tool vaultpulse, got %q", received["monitoring_tool"])
	}
}

func TestVictorOpsClient_Send_NonOKStatus_ReturnsError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	client, _ := notify.NewVictorOpsClient(server.URL)
	if err := client.Send("test"); err == nil {
		t.Fatal("expected error on non-OK status, got nil")
	}
}

func TestVictorOpsClient_Send_UnreachableServer_ReturnsError(t *testing.T) {
	client, _ := notify.NewVictorOpsClient("http://127.0.0.1:0/webhook")
	if err := client.Send("test"); err == nil {
		t.Fatal("expected error for unreachable server, got nil")
	}
}
