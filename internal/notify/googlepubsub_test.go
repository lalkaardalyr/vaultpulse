package notify

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewGooglePubSubClient_EmptyProjectID_ReturnsError(t *testing.T) {
	_, err := NewGooglePubSubClient("", "topic", "key")
	if err == nil {
		t.Fatal("expected error for empty project ID")
	}
}

func TestNewGooglePubSubClient_EmptyTopicID_ReturnsError(t *testing.T) {
	_, err := NewGooglePubSubClient("project", "", "key")
	if err == nil {
		t.Fatal("expected error for empty topic ID")
	}
}

func TestNewGooglePubSubClient_EmptyAPIKey_ReturnsError(t *testing.T) {
	_, err := NewGooglePubSubClient("project", "topic", "")
	if err == nil {
		t.Fatal("expected error for empty API key")
	}
}

func TestNewGooglePubSubClient_ValidConfig_ReturnsClient(t *testing.T) {
	c, err := NewGooglePubSubClient("project", "topic", "key")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if c == nil {
		t.Fatal("expected non-nil client")
	}
}

func TestGooglePubSubClient_Send_PostsCorrectPayload(t *testing.T) {
	var received pubsubMessage
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(body, &received)
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	c, _ := NewGooglePubSubClient("project", "topic", "key")
	c.httpClient = ts.Client()
	c.endpoint = ts.URL

	if err := c.Send(context.Background(), "hello vault"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(received.Messages) != 1 {
		t.Fatalf("expected 1 message, got %d", len(received.Messages))
	}
	decoded, _ := base64.StdEncoding.DecodeString(received.Messages[0].Data)
	if string(decoded) != "hello vault" {
		t.Errorf("unexpected message data: %s", decoded)
	}
}

func TestGooglePubSubClient_Send_NonOKStatus_ReturnsError(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
	}))
	defer ts.Close()

	c, _ := NewGooglePubSubClient("project", "topic", "key")
	c.httpClient = ts.Client()
	c.endpoint = ts.URL

	if err := c.Send(context.Background(), "msg"); err == nil {
		t.Fatal("expected error for non-OK status")
	}
}

func TestGooglePubSubClient_Send_UnreachableServer_ReturnsError(t *testing.T) {
	c, _ := NewGooglePubSubClient("project", "topic", "key")
	c.endpoint = "http://127.0.0.1:0"
	if err := c.Send(context.Background(), "msg"); err == nil {
		t.Fatal("expected error for unreachable server")
	}
}
