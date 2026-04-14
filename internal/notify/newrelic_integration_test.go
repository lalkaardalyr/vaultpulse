package notify_test

import (
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"

	"github.com/yourusername/vaultpulse/internal/notify"
)

func TestNewRelicClient_MultiSender_Integration(t *testing.T) {
	var callCount int32

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&callCount, 1)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	nr1, err := notify.NewNewRelicClient("111", "key-a")
	if err != nil {
		t.Fatalf("client 1 error: %v", err)
	}
	nr1.(*notify.NewRelicClient) // type assertion not needed; use via Sender interface

	// Use MultiSender to combine two logical senders pointing at the test server.
	// We reuse the same client twice to simulate two independent senders.
	nr2, err := notify.NewNewRelicClient("222", "key-b")
	if err != nil {
		t.Fatalf("client 2 error: %v", err)
	}

	// Override endpoints to point at test server.
	patchEndpoint(nr1, server.URL)
	patchEndpoint(nr2, server.URL)

	multi, err := notify.NewMultiSender(nr1, nr2)
	if err != nil {
		t.Fatalf("multi sender error: %v", err)
	}

	if err := multi.Send("integration test alert"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if got := atomic.LoadInt32(&callCount); got != 2 {
		t.Errorf("expected 2 calls, got %d", got)
	}
}

// patchEndpoint uses the exported field via a helper to update the endpoint
// without breaking encapsulation — acceptable in integration tests.
func patchEndpoint(s notify.Sender, url string) {
	if nr, ok := s.(interface{ SetEndpoint(string) }); ok {
		nr.SetEndpoint(url)
	}
}
