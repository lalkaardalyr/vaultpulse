package notify

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewInfluxDBClient_EmptyURL_ReturnsError(t *testing.T) {
	_, err := NewInfluxDBClient("", "token", "org", "bucket")
	if err == nil {
		t.Fatal("expected error for empty URL")
	}
}

func TestNewInfluxDBClient_EmptyToken_ReturnsError(t *testing.T) {
	_, err := NewInfluxDBClient("http://localhost", "", "org", "bucket")
	if err == nil {
		t.Fatal("expected error for empty token")
	}
}

func TestNewInfluxDBClient_EmptyOrg_ReturnsError(t *testing.T) {
	_, err := NewInfluxDBClient("http://localhost", "token", "", "bucket")
	if err == nil {
		t.Fatal("expected error for empty org")
	}
}

func TestNewInfluxDBClient_EmptyBucket_ReturnsError(t *testing.T) {
	_, err := NewInfluxDBClient("http://localhost", "token", "org", "")
	if err == nil {
		t.Fatal("expected error for empty bucket")
	}
}

func TestNewInfluxDBClient_ValidConfig_ReturnsClient(t *testing.T) {
	c, err := NewInfluxDBClient("http://localhost", "token", "org", "bucket")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if c == nil {
		t.Fatal("expected non-nil client")
	}
}

func TestInfluxDBClient_Send_PostsCorrectPayload(t *testing.T) {
	var called bool
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
		if r.Header.Get("Authorization") == "" {
			t.Error("expected Authorization header")
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer ts.Close()

	c, _ := NewInfluxDBClient(ts.URL, "token", "org", "bucket")
	err := c.Send(context.Background(), "test alert")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !called {
		t.Fatal("expected server to be called")
	}
}

func TestInfluxDBClient_Send_NonOKStatus_ReturnsError(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer ts.Close()

	c, _ := NewInfluxDBClient(ts.URL, "token", "org", "bucket")
	err := c.Send(context.Background(), "test alert")
	if err == nil {
		t.Fatal("expected error for non-OK status")
	}
}
