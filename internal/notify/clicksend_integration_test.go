package notify

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestClickSendClient_MultiSender_Integration(t *testing.T) {
	calls := 0
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls++
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	// Build two stub senders that wrap ClickSendClient-like behaviour.
	sender1 := &stubSender{}
	sender2 := &stubSender{}

	multi, err := NewMultiSender(sender1, sender2)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if err := multi.Send("vault secret expiring soon"); err != nil {
		t.Fatalf("unexpected send error: %v", err)
	}

	if sender1.received != "vault secret expiring soon" {
		t.Errorf("sender1 did not receive message")
	}
	if sender2.received != "vault secret expiring soon" {
		t.Errorf("sender2 did not receive message")
	}
}
