package notify

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestNewGrafanaClient_EmptyURL_ReturnsError(t *testing.T) {
	_, err := NewGrafanaClient("", "some-key")
	if err == nil {
		t.Fatal("expected error for empty URL, got nil")
	}
}

func TestNewGrafanaClient_EmptyKey_ReturnsError(t *testing.T) {
	_, err := NewGrafanaClient("https://grafana.example.com", "")
	if err == nil {
		t.Fatal("expected error for empty API key, got nil")
	}
}

func TestNewGrafanaClient_ValidConfig_ReturnsClient(t *testing.T) {
	c, err := NewGrafanaClient("https://grafana.example.com", "token123")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if c == nil {
		t.Fatal("expected non-nil client")
	}
}

func TestGrafanaClient_Send_PostsCorrectPayload(t *testing.T) {
	var received grafanaAnnotation
	var authHeader string

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader = r.Header.Get("Authorization")
		body, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(body, &received)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	c, _ := NewGrafanaClient(server.URL, "test-token")
	if err := c.Send("secret expiring soon"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if received.Text != "secret expiring soon" {
		t.Errorf("expected text %q, got %q", "secret expiring soon", received.Text)
	}
	if !strings.Contains(authHeader, "test-token") {
		t.Errorf("expected Authorization header to contain token, got %q", authHeader)
	}
}

func TestGrafanaClient_Send_NonOKStatus_ReturnsError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
	}))
	defer server.Close()

	c, _ := NewGrafanaClient(server.URL, "bad-token")
	if err := c.Send("test"); err == nil {
		t.Fatal("expected error on non-OK status, got nil")
	}
}

func TestGrafanaClient_Send_UnreachableServer_ReturnsError(t *testing.T) {
	c, _ := NewGrafanaClient("http://127.0.0.1:19999", "token")
	if err := c.Send("test"); err == nil {
		t.Fatal("expected error for unreachable server, got nil")
	}
}
