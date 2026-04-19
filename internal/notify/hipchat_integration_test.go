package notify_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/dkoosis/vaultpulse/internal/notify"
)

func TestHipChatClient_MultiSender_Integration(t *testing.T) {
	var received int
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		received++
		w.WriteHeader(http.StatusNoContent)
	}))
	defer ts.Close()

	c1, err := notify.NewHipChatClient(ts.URL, "token-a")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	c2, err := notify.NewHipChatClient(ts.URL, "token-b")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	multi, err := notify.NewMultiSender(c1, c2)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if err := multi.Send("integration test"); err != nil {
		t.Fatalf("Send returned error: %v", err)
	}

	if received != 2 {
		t.Errorf("expected 2 requests, got %d", received)
	}
}
