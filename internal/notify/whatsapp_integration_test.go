package notify

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestWhatsAppClient_MultiSender_Integration(t *testing.T) {
	var calls int
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls++
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	wa, err := NewWhatsAppClient("token", "phone123", "+15550001234")
	if err != nil {
		t.Fatalf("setup error: %v", err)
	}
	wa.endpoint = ts.URL

	slack, err := NewSlackClient(ts.URL)
	if err != nil {
		t.Fatalf("setup slack error: %v", err)
	}

	multi, err := NewMultiSender(wa, slack)
	if err != nil {
		t.Fatalf("multi sender error: %v", err)
	}

	if err := multi.Send("vault secret expiring"); err != nil {
		t.Fatalf("unexpected send error: %v", err)
	}

	if calls != 2 {
		t.Errorf("expected 2 HTTP calls, got %d", calls)
	}
}
