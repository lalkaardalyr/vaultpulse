package notify

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewClickatellClient_EmptyAPIKey_ReturnsError(t *testing.T) {
	_, err := NewClickatellClient("", "+12025550100")
	if err == nil {
		t.Fatal("expected error for empty API key")
	}
}

func TestNewClickatellClient_EmptyTo_ReturnsError(t *testing.T) {
	_, err := NewClickatellClient("key", "")
	if err == nil {
		t.Fatal("expected error for empty recipient")
	}
}

func TestNewClickatellClient_ValidConfig_ReturnsClient(t *testing.T) {
	c, err := NewClickatellClient("key", "+12025550100")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if c == nil {
		t.Fatal("expected non-nil client")
	}
}

func TestClickatellClient_Send_PostsCorrectPayload(t *testing.T) {
	var gotAuth, gotBody string
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotAuth = r.Header.Get("Authorization")
		body := make([]byte, r.ContentLength)
		r.Body.Read(body)
		gotBody = string(body)
		w.WriteHeader(http.StatusAccepted)
	}))
	defer ts.Close()

	c, _ := NewClickatellClient("testkey", "+12025550100")
	c.(*clickatellClient).endpoint = ts.URL

	err := c.Send(context.Background(), "test alert")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotAuth != "Bearer testkey" {
		t.Errorf("expected Bearer auth, got %q", gotAuth)
	}
	if gotBody == "" {
		t.Error("expected non-empty body")
	}
}

func TestClickatellClient_Send_NonOKStatus_ReturnsError(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	}))
	defer ts.Close()

	c, _ := NewClickatellClient("badkey", "+12025550100")
	c.(*clickatellClient).endpoint = ts.URL

	err := c.Send(context.Background(), "msg")
	if err == nil {
		t.Fatal("expected error on non-OK status")
	}
}

func TestClickatellClient_Send_UnreachableServer_ReturnsError(t *testing.T) {
	c, _ := NewClickatellClient("key", "+12025550100")
	c.(*clickatellClient).endpoint = "http://127.0.0.1:0"

	err := c.Send(context.Background(), "msg")
	if err == nil {
		t.Fatal("expected error for unreachable server")
	}
}
