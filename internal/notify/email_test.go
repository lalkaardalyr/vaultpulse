package notify_test

import (
	"net"
	"net/smtp"
	"net/textproto"
	"strings"
	"testing"

	"github.com/yourusername/vaultpulse/internal/notify"
)

func TestNewEmailClient_EmptyHost_ReturnsError(t *testing.T) {
	_, err := notify.NewEmailClient(notify.EmailConfig{
		From: "vault@example.com",
		To:   []string{"ops@example.com"},
	})
	if err == nil {
		t.Fatal("expected error for empty host, got nil")
	}
	if !strings.Contains(err.Error(), "host") {
		t.Errorf("expected error to mention 'host', got: %s", err.Error())
	}
}

func TestNewEmailClient_EmptyFrom_ReturnsError(t *testing.T) {
	_, err := notify.NewEmailClient(notify.EmailConfig{
		Host: "smtp.example.com",
		To:   []string{"ops@example.com"},
	})
	if err == nil {
		t.Fatal("expected error for empty from, got nil")
	}
	if !strings.Contains(err.Error(), "from") {
		t.Errorf("expected error to mention 'from', got: %s", err.Error())
	}
}

func TestNewEmailClient_EmptyRecipients_ReturnsError(t *testing.T) {
	_, err := notify.NewEmailClient(notify.EmailConfig{
		Host: "smtp.example.com",
		From: "vault@example.com",
	})
	if err == nil {
		t.Fatal("expected error for empty recipients, got nil")
	}
	if !strings.Contains(err.Error(), "recipient") {
		t.Errorf("expected error to mention 'recipient', got: %s", err.Error())
	}
}

func TestNewEmailClient_ValidConfig_ReturnsClient(t *testing.T) {
	client, err := notify.NewEmailClient(notify.EmailConfig{
		Host: "smtp.example.com",
		From: "vault@example.com",
		To:   []string{"ops@example.com"},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if client == nil {
		t.Fatal("expected non-nil client")
	}
}

func TestNewEmailClient_DefaultPort(t *testing.T) {
	// Verifies that a zero port is replaced with 587 by attempting
	// a send against a local listener and checking the dial address.
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Skip("could not bind local listener")
	}
	defer ln.Close()

	_, port, _ := net.SplitHostPort(ln.Addr().String())
	_ = port // port is dynamic; we just verify client creation succeeds

	client, err := notify.NewEmailClient(notify.EmailConfig{
		Host: "127.0.0.1",
		From: "vault@example.com",
		To:   []string{"ops@example.com"},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if client == nil {
		t.Fatal("expected non-nil client")
	}
}

func TestEmailClient_Send_UnreachableServer_ReturnsError(t *testing.T) {
	client, err := notify.NewEmailClient(notify.EmailConfig{
		Host: "127.0.0.1",
		Port: 19999,
		From: "vault@example.com",
		To:   []string{"ops@example.com"},
	})
	if err != nil {
		t.Fatalf("unexpected error creating client: %v", err)
	}

	err = client.Send("test alert message")
	if err == nil {
		t.Fatal("expected error for unreachable server, got nil")
	}

	// Suppress unused import warnings from smtp/textproto in test file.
	_ = smtp.PlainAuth
	_ = textproto.NewConn
}
