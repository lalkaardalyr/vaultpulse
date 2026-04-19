package notify

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestNewLinearClient_EmptyAPIKey_ReturnsError(t *testing.T) {
	_, err := NewLinearClient("", "team123")
	if err == nil {
		t.Fatal("expected error for empty API key")
	}
}

func TestNewLinearClient_EmptyTeamID_ReturnsError(t *testing.T) {
	_, err := NewLinearClient("key", "")
	if err == nil {
		t.Fatal("expected error for empty team ID")
	}
}

func TestNewLinearClient_ValidConfig_ReturnsClient(t *testing.T) {
	c, err := NewLinearClient("key", "team123")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if c == nil {
		t.Fatal("expected non-nil client")
	}
}

func TestLinearClient_Send_PostsCorrectPayload(t *testing.T) {
	var gotAuth, gotBody string
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotAuth = r.Header.Get("Authorization")
		var buf strings.Builder
		j := json.NewDecoder(r.Body)
		var raw map[string]string
		_ = j.Decode(&raw)
		gotBody = raw["query"]
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"data":{"issueCreate":{"success":true}}}`))
	}))
	defer ts.Close()

	c, _ := NewLinearClient("mykey", "team1")
	c.endpoint = ts.URL

	if err := c.Send("vault secret expiring"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotAuth != "mykey" {
		t.Errorf("expected auth header 'mykey', got %q", gotAuth)
	}
	if !strings.Contains(gotBody, "vault secret expiring") {
		t.Errorf("expected body to contain message, got %q", gotBody)
	}
}

func TestLinearClient_Send_NonOKStatus_ReturnsError(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"data":{"issueCreate":{"success":false}}}`))
	}))
	defer ts.Close()

	c, _ := NewLinearClient("key", "team1")
	c.endpoint = ts.URL

	if err := c.Send("alert"); err == nil {
		t.Fatal("expected error on failed issue creation")
	}
}
