package notify

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewStatuspageClient_EmptyAPIKey_ReturnsError(t *testing.T) {
	_, err := NewStatuspageClient("", "page123")
	if err == nil {
		t.Fatal("expected error for empty api key")
	}
}

func TestNewStatuspageClient_EmptyPageID_ReturnsError(t *testing.T) {
	_, err := NewStatuspageClient("key", "")
	if err == nil {
		t.Fatal("expected error for empty page ID")
	}
}

func TestNewStatuspageClient_ValidConfig_ReturnsClient(t *testing.T) {
	c, err := NewStatuspageClient("key", "page123")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if c == nil {
		t.Fatal("expected non-nil client")
	}
}

func TestStatuspageClient_Send_PostsCorrectPayload(t *testing.T) {
	var received map[string]interface{}
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(body, &received)
		w.WriteHeader(http.StatusCreated)
	}))
	defer ts.Close()

	c, _ := NewStatuspageClient("key", "page123")
	c.endpoint = ts.URL

	if err := c.Send("vault secret expiring soon"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	inc, ok := received["incident"].(map[string]interface{})
	if !ok {
		t.Fatal("expected incident key in payload")
	}
	if inc["body"] != "vault secret expiring soon" {
		t.Errorf("unexpected body: %v", inc["body"])
	}
}

func TestStatuspageClient_Send_NonOKStatus_ReturnsError(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	}))
	defer ts.Close()

	c, _ := NewStatuspageClient("bad-key", "page123")
	c.endpoint = ts.URL

	if err := c.Send("test"); err == nil {
		t.Fatal("expected error on non-2xx status")
	}
}

func TestStatuspageClient_Send_UnreachableServer_ReturnsError(t *testing.T) {
	c, _ := NewStatuspageClient("key", "page123")
	c.endpoint = "http://127.0.0.1:0"
	if err := c.Send("test"); err == nil {
		t.Fatal("expected error for unreachable server")
	}
}
