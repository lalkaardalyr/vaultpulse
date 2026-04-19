package notify_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/yourusername/vaultpulse/internal/notify"
)

func TestSpikeClient_MultiSender_Integration(t *testing.T) {
	var received []string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		received = append(received, r.URL.Path)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	spike, err := notify.NewSpikeClient(server.URL)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	multi, err := notify.NewMultiSender(spike)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if err := multi.Send("vault secret expiring soon"); err != nil {
		t.Fatalf("Send failed: %v", err)
	}

	if len(received) != 1 {
		t.Errorf("expected 1 request, got %d", len(received))
	}
}
