package notify

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewPagerDutyClient_EmptyKey_ReturnsError(t *testing.T) {
	_, err := NewPagerDutyClient("")
	if err == nil {
		t.Fatal("expected error for empty integration key, got nil")
	}
}

func TestNewPagerDutyClient_ValidKey_ReturnsClient(t *testing.T) {
	client, err := NewPagerDutyClient("test-key-123")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if client == nil {
		t.Fatal("expected non-nil client")
	}
}

func TestPagerDutyClient_Send_PostsCorrectPayload(t *testing.T) {
	var received pagerDutyPayload

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(body, &received)
		w.WriteHeader(http.StatusAccepted)
	}))
	defer server.Close()

	client, _ := NewPagerDutyClient("routing-key-abc")
	client.eventsURL = server.URL

	details := map[string]string{"path": "secret/db", "days_left": "3"}
	err := client.Send("Vault secret expiring soon", "warning", details)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if received.RoutingKey != "routing-key-abc" {
		t.Errorf("expected routing key 'routing-key-abc', got %q", received.RoutingKey)
	}
	if received.EventAction != "trigger" {
		t.Errorf("expected event_action 'trigger', got %q", received.EventAction)
	}
	if received.Payload.Severity != "warning" {
		t.Errorf("expected severity 'warning', got %q", received.Payload.Severity)
	}
	if received.Payload.Source != "vaultpulse" {
		t.Errorf("expected source 'vaultpulse', got %q", received.Payload.Source)
	}
	if received.Payload.CustomDetails["path"] != "secret/db" {
		t.Errorf("expected custom_details path 'secret/db', got %q", received.Payload.CustomDetails["path"])
	}
}

func TestPagerDutyClient_Send_NonOKStatus_ReturnsError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
	}))
	defer server.Close()

	client, _ := NewPagerDutyClient("key")
	client.eventsURL = server.URL

	err := client.Send("test", "critical", nil)
	if err == nil {
		t.Fatal("expected error for non-2xx status, got nil")
	}
}

func TestPagerDutyClient_Send_UnreachableServer_ReturnsError(t *testing.T) {
	client, _ := NewPagerDutyClient("key")
	client.eventsURL = "http://127.0.0.1:1"

	err := client.Send("test", "info", nil)
	if err == nil {
		t.Fatal("expected error for unreachable server, got nil")
	}
}
