package notify

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewJiraClient_EmptyBaseURL_ReturnsError(t *testing.T) {
	_, err := NewJiraClient("", "user@example.com", "token", "OPS")
	if err == nil {
		t.Fatal("expected error for empty baseURL")
	}
}

func TestNewJiraClient_EmptyEmail_ReturnsError(t *testing.T) {
	_, err := NewJiraClient("https://example.atlassian.net", "", "token", "OPS")
	if err == nil {
		t.Fatal("expected error for empty email")
	}
}

func TestNewJiraClient_EmptyAPIToken_ReturnsError(t *testing.T) {
	_, err := NewJiraClient("https://example.atlassian.net", "user@example.com", "", "OPS")
	if err == nil {
		t.Fatal("expected error for empty apiToken")
	}
}

func TestNewJiraClient_EmptyProject_ReturnsError(t *testing.T) {
	_, err := NewJiraClient("https://example.atlassian.net", "user@example.com", "token", "")
	if err == nil {
		t.Fatal("expected error for empty project")
	}
}

func TestNewJiraClient_ValidConfig_ReturnsClient(t *testing.T) {
	c, err := NewJiraClient("https://example.atlassian.net", "user@example.com", "token", "OPS")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if c == nil {
		t.Fatal("expected non-nil client")
	}
}

func TestJiraClient_Send_PostsCorrectPayload(t *testing.T) {
	var received jiraIssuePayload
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := json.NewDecoder(r.Body).Decode(&received); err != nil {
			t.Errorf("decode body: %v", err)
		}
		w.WriteHeader(http.StatusCreated)
	}))
	defer ts.Close()

	c, _ := NewJiraClient(ts.URL, "user@example.com", "token", "OPS")
	if err := c.Send("secret expiring soon"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if received.Fields.Summary != "secret expiring soon" {
		t.Errorf("expected summary %q, got %q", "secret expiring soon", received.Fields.Summary)
	}
	if received.Fields.Project.Key != "OPS" {
		t.Errorf("expected project key OPS, got %q", received.Fields.Project.Key)
	}
}

func TestJiraClient_Send_NonOKStatus_ReturnsError(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
	}))
	defer ts.Close()

	c, _ := NewJiraClient(ts.URL, "user@example.com", "token", "OPS")
	if err := c.Send("alert"); err == nil {
		t.Fatal("expected error on non-2xx status")
	}
}

func TestJiraClient_Send_UnreachableServer_ReturnsError(t *testing.T) {
	c, _ := NewJiraClient("http://127.0.0.1:19999", "user@example.com", "token", "OPS")
	if err := c.Send("alert"); err == nil {
		t.Fatal("expected error on unreachable server")
	}
}
