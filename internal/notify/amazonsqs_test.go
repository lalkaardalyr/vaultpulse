package notify

import (
	"context"
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

type mockSQSSender struct {
	called  bool
	body    string
	returnErr error
}

func (m *mockSQSSender) SendMessage(_ context.Context, params *sqs.SendMessageInput, _ ...func(*sqs.Options)) (*sqs.SendMessageOutput, error) {
	m.called = true
	if params.MessageBody != nil {
		m.body = *params.MessageBody
	}
	return &sqs.SendMessageOutput{}, m.returnErr
}

func TestNewSQSClient_EmptyQueueURL_ReturnsError(t *testing.T) {
	_, err := newSQSClientWithSender("", &mockSQSSender{})
	if err == nil {
		t.Fatal("expected error for empty queue URL")
	}
}

func TestNewSQSClient_ValidQueueURL_ReturnsClient(t *testing.T) {
	c, err := newSQSClientWithSender("https://sqs.us-east-1.amazonaws.com/123456789012/my-queue", &mockSQSSender{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if c == nil {
		t.Fatal("expected non-nil client")
	}
}

func TestSQSClient_Send_PostsCorrectPayload(t *testing.T) {
	mock := &mockSQSSender{}
	c, _ := newSQSClientWithSender("https://sqs.us-east-1.amazonaws.com/123456789012/my-queue", mock)

	if err := c.Send("vault secret expiring soon"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !mock.called {
		t.Fatal("expected SendMessage to be called")
	}
	if mock.body == "" {
		t.Fatal("expected non-empty message body")
	}
}

func TestSQSClient_Send_ReturnsError_OnFailure(t *testing.T) {
	mock := &mockSQSSender{returnErr: errors.New("sqs unavailable")}
	c, _ := newSQSClientWithSender("https://sqs.us-east-1.amazonaws.com/123456789012/my-queue", mock)

	if err := c.Send("test"); err == nil {
		t.Fatal("expected error from failed SendMessage")
	}
}
