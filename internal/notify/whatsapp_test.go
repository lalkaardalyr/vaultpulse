package notify

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewWhatsAppClient_EmptyToken_ReturnsError(t *testing.T) {
	_, err := NewWhatsAppClient("", "phone123", "+15550001234")
	if err == nil {
		t.Fatal("expected error for empty token")
	}
}

func TestNewWhatsAppClient_EmptyPhoneID_ReturnsError(t *testing.T) {
	_, err := NewWhatsAppClient("token", "", "+15550001234")
	if err == nil {
		t.Fatal("expected error for empty phone ID")
	}
}

func TestNewWhatsAppClient_EmptyRecipient_ReturnsError(t *testing.T) {
	_, err := NewWhatsAppClient("token", "phone123", "")
	if err == nil {
		t.Fatal("expected error for empty recipient")
	}
}

func TestNewWhatsAppClient_ValidConfig_ReturnsClient(t *testing.T) {
	client, err := NewWhatsAppClient("token", "phone123", "+15550001234")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if client == nil {
		t.Fatal("expected non-nil client")
	}
}

func TestWhatsAppClient_Send_PostsCorrectPayload(t *testing.T) {
	var gotAuth, gotContentType string
	var gotBody []byte

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotAuth = r.Header.Get("Authorization")
		gotContentType = r.Header.Get("Content-Type")
		gotBody = make([]byte, r.ContentLength)
		r.Body.Read(gotBody)
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	client, _ := NewWhatsAppClient("mytoken", "phone123", "+15550001234")
	client.endpoint = ts.URL

	if err := client.Send("test alert"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotAuth != "Bearer mytoken" {
		t.Errorf("expected Bearer mytoken, got %s", gotAuth)
	}
	if gotContentType != "application/json" {
		t.Errorf("expected application/json, got %s", gotContentType)
	}
	if len(gotBody) == 0 {
		t.Error("expected non-empty body")
	}
}

func TestWhatsAppClient_Send_NonOKStatus_ReturnsError(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	}))
	defer ts.Close()

	client, _ := NewWhatsAppClient("token", "phone123", "+15550001234")
	client.endpoint = ts.URL

	if err := client.Send("msg"); err == nil {
		t.Fatal("expected error for non-OK status")
	}
}

func TestWhatsAppClient_Send_UnreachableServer_ReturnsError(t *testing.T) {
	client, _ := NewWhatsAppClient("token", "phone123", "+15550001234")
	client.endpoint = "http://127.0.0.1:0"

	if err := client.Send("msg"); err == nil {
		t.Fatal("expected error for unreachable server")
	}
}
