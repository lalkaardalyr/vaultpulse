package notify

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewCircleCIClient_EmptyToken_ReturnsError(t *testing.T) {
	_, err := NewCircleCIClient("", "gh/org/repo")
	if err == nil {
		t.Fatal("expected error for empty token")
	}
}

func TestNewCircleCIClient_EmptyProject_ReturnsError(t *testing.T) {
	_, err := NewCircleCIClient("token123", "")
	if err == nil {
		t.Fatal("expected error for empty project")
	}
}

func TestNewCircleCIClient_ValidConfig_ReturnsClient(t *testing.T) {
	c, err := NewCircleCIClient("token123", "gh/org/repo")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if c == nil {
		t.Fatal("expected non-nil client")
	}
}

func TestCircleCIClient_Send_PostsCorrectPayload(t *testing.T) {
	var received map[string]interface{}
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(body, &received)
		w.WriteHeader(http.StatusCreated)
	}))
	defer ts.Close()

	c, _ := NewCircleCIClient("tok", "gh/org/repo")
	c.endpoint = ts.URL

	if err := c.Send("secret expiring soon"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	params, ok := received["parameters"].(map[string]interface{})
	if !ok {
		t.Fatal("expected parameters field")
	}
	if params["alert_message"] != "secret expiring soon" {
		t.Errorf("unexpected message: %v", params["alert_message"])
	}
}

func TestCircleCIClient_Send_NonOKStatus_ReturnsError(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
	}))
	defer ts.Close()

	c, _ := NewCircleCIClient("tok", "gh/org/repo")
	c.endpoint = ts.URL

	if err := c.Send("msg"); err == nil {
		t.Fatal("expected error for non-OK status")
	}
}

func TestCircleCIClient_Send_UnreachableServer_ReturnsError(t *testing.T) {
	c, _ := NewCircleCIClient("tok", "gh/org/repo")
	c.endpoint = "http://127.0.0.1:0"
	if err := c.Send("msg"); err == nil {
		t.Fatal("expected error for unreachable server")
	}
}
