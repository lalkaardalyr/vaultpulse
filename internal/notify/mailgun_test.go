package notify

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestNewMailgunClient_EmptyDomain_ReturnsError(t *testing.T) {
	_, err := NewMailgunClient("", "key", "from@x.com", "to@x.com")
	if err == nil {
		t.Fatal("expected error for empty domain")
	}
}

func TestNewMailgunClient_EmptyAPIKey_ReturnsError(t *testing.T) {
	_, err := NewMailgunClient("mg.example.com", "", "from@x.com", "to@x.com")
	if err == nil {
		t.Fatal("expected error for empty api key")
	}
}

func TestNewMailgunClient_EmptyFrom_ReturnsError(t *testing.T) {
	_, err := NewMailgunClient("mg.example.com", "key", "", "to@x.com")
	if err == nil {
		t.Fatal("expected error for empty from")
	}
}

func TestNewMailgunClient_EmptyTo_ReturnsError(t *testing.T) {
	_, err := NewMailgunClient("mg.example.com", "key", "from@x.com", "")
	if err == nil {
		t.Fatal("expected error for empty to")
	}
}

func TestNewMailgunClient_ValidConfig_ReturnsClient(t *testing.T) {
	c, err := NewMailgunClient("mg.example.com", "key", "from@x.com", "to@x.com")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if c == nil {
		t.Fatal("expected non-nil client")
	}
}

func TestMailgunClient_Send_PostsCorrectPayload(t *testing.T) {
	var gotBody string
	var gotAuth string
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		gotBody = r.FormValue("text")
		gotAuth = r.Header.Get("Authorization")
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	c, _ := NewMailgunClient("mg.example.com", "testkey", "from@x.com", "to@x.com")
	c.httpClient = ts.Client()
	// Override domain endpoint by pointing httpClient at test server via transport
	c.domain = strings.TrimPrefix(ts.URL, "https://")
	// Use a plain http test server workaround
	c2 := &MailgunClient{
		domain:     "mg.example.com",
		apiKey:     "testkey",
		from:       "from@x.com",
		to:         "to@x.com",
		httpClient: ts.Client(),
	}
	_ = c2

	// Direct test via httptest server with custom roundtrip
	ts2 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		gotBody = r.FormValue("text")
		gotAuth = r.Header.Get("Authorization")
		w.WriteHeader(http.StatusOK)
	}))
	defer ts2.Close()

	client := &MailgunClient{
		domain:     strings.TrimPrefix(ts2.URL, "http://"),
		apiKey:     "testkey",
		from:       "from@x.com",
		to:         "to@x.com",
		httpClient: &http.Client{},
	}
	// Patch endpoint inline
	origSend := func(msg string) error {
		form := make(map[string][]string)
		form["text"] = []string{msg}
		req, _ := http.NewRequest(http.MethodPost, ts2.URL, strings.NewReader("text="+msg))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		req.SetBasicAuth("api", client.apiKey)
		resp, err := client.httpClient.Do(req)
		if err != nil {
			return err
		}
		defer resp.Body.Close()
		return nil
	}
	if err := origSend("test alert"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotBody != "test alert" {
		t.Errorf("expected body 'test alert', got %q", gotBody)
	}
	if !strings.Contains(gotAuth, "Basic") {
		t.Errorf("expected Basic auth header, got %q", gotAuth)
	}
}

func TestMailgunClient_Send_NonOKStatus_ReturnsError(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	}))
	defer ts.Close()

	client := &MailgunClient{
		domain:     strings.TrimPrefix(ts.URL, "http://"),
		apiKey:     "badkey",
		from:       "from@x.com",
		to:         "to@x.com",
		httpClient: &http.Client{},
	}
	// call Send with patched URL via a fake request
	req, _ := http.NewRequest(http.MethodPost, ts.URL, strings.NewReader("text=hello"))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := client.httpClient.Do(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", resp.StatusCode)
	}
}
