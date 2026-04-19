package notify_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/yourusername/vaultpulse/internal/notify"
)

func TestMailgunClient_MultiSender_Integration(t *testing.T) {
	var received []string

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		received = append(received, r.FormValue("text"))
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	slack, err := notify.NewSlackClient(ts.URL)
	if err != nil {
		t.Fatalf("slack setup error: %v", err)
	}

	multi, err := notify.NewMultiSender(slack)
	if err != nil {
		t.Fatalf("multi setup error: %v", err)
	}

	if err := multi.Send("vault secret expiring soon"); err != nil {
		t.Fatalf("unexpected send error: %v", err)
	}
}
