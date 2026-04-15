package notify

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewZendutyClient_EmptyKey_ReturnsError(t *testing.T) {
	_, err := NewZendutyClient("")
	if err == nil {
		t.Fatal("expected error for empty integration key, got nil")
	}
}

func TestNewZendutyClient_ValidKey_ReturnsClient(t *testing.T) {
	client, err := NewZendutyClient("test-integration-key")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if client == nil {
		t.Fatal("expected non-nil client")
	}
}

func TestZendutyClient_Send_PostsCorrectPayload(t *testing.T) {
	var received zendutyPayload

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(body, &received)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client, _ := NewZendutyClient("my-key")
	client.endpoint = server.URL + "/api/events/"

	if err := client.Send("vault secret expiring soon"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if received.Message != "vault secret expiring soon" {
		t.Errorf("expected message %q, got %q", "vault secret expiring soon", received.Message)
	}
	if received.AlertType != "critical" {
		t.Errorf("expected alert_type %q, got %q", "critical", received.AlertType)
	}
}

func TestZendutyClient_Send_PostsIntegrationKey(t *testing.T) {
	var received zendutyPayload

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(body, &received)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client, _ := NewZendutyClient("expected-key")
	client.endpoint = server.URL + "/api/events/"

	if err := client.Send("test"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if received.IntegrationKey != "expected-key" {
		t.Errorf("expected integration_key %q, got %q", "expected-key", received.IntegrationKey)
	}
}

func TestZendutyClient_Send_NonOKStatus_ReturnsError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	client, _ := NewZendutyClient("my-key")
	client.endpoint = server.URL + "/api/events/"

	if err := client.Send("test"); err == nil {
		t.Fatal("expected error for non-OK status, got nil")
	}
}

func TestZendutyClient_Send_UnreachableServer_ReturnsError(t *testing.T) {
	client, _ := NewZendutyClient("my-key")
	client.endpoint = "http://127.0.0.1:0/api/events/"

	if err := client.Send("test"); err == nil {
		t.Fatal("expected error for unreachable server, got nil")
	}
}
