package notify_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	"github.com/yourusername/vaultpulse/internal/notify"
)

func TestGoogleChatClient_MultiSender_Integration(t *testing.T) {
	var mu sync.Mutex
	var received []string

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var payload map[string]string
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			t.Errorf("decode error: %v", err)
		}
		mu.Lock()
		received = append(received, payload["text"])
		mu.Unlock()
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	gc, err := notify.NewGoogleChatClient(server.URL)
	if err != nil {
		t.Fatalf("NewGoogleChatClient: %v", err)
	}

	multi, err := notify.NewMultiSender(gc)
	if err != nil {
		t.Fatalf("NewMultiSender: %v", err)
	}

	const msg = "[WARNING] secret/api/key expires in 5 days"
	if err := multi.Send(msg); err != nil {
		t.Fatalf("Send: %v", err)
	}

	mu.Lock()
	defer mu.Unlock()
	if len(received) != 1 || received[0] != msg {
		t.Errorf("expected one message %q, got %v", msg, received)
	}
}
