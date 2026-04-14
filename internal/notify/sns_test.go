package notify

import (
	"context"
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/sns"
)

// mockSNSPublisher is a test double for the SNS Publish call.
type mockSNSPublisher struct {
	called  bool
	input   *sns.PublishInput
	retErr  error
}

func (m *mockSNSPublisher) Publish(_ context.Context, input *sns.PublishInput, _ ...func(*sns.Options)) (*sns.PublishOutput, error) {
	m.called = true
	m.input = input
	return &sns.PublishOutput{}, m.retErr
}

func newTestSNSClient(topicARN string, pub snsPublisher) *snsClientWrapper {
	return &snsClientWrapper{pub: pub, topicARN: topicARN}
}

func (w *snsClientWrapper) Send(msg string) error {
	_, err := w.pub.Publish(context.Background(), &sns.PublishInput{
		TopicArn: strPtr(w.topicARN),
		Message:  strPtr(msg),
	})
	if err != nil {
		return err
	}
	return nil
}

func strPtr(s string) *string { return &s }

func TestNewSNSClient_EmptyARN_ReturnsError(t *testing.T) {
	_, err := NewSNSClient("")
	if err == nil {
		t.Fatal("expected error for empty topic ARN")
	}
}

func TestSNSClient_Send_CallsPublish(t *testing.T) {
	mock := &mockSNSPublisher{}
	client := newTestSNSClient("arn:aws:sns:us-east-1:123456789012:alerts", mock)

	if err := client.Send("test alert"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !mock.called {
		t.Error("expected Publish to be called")
	}
	if mock.input == nil || *mock.input.Message != "test alert" {
		t.Errorf("unexpected message payload: %v", mock.input)
	}
}

func TestSNSClient_Send_PublishError_ReturnsError(t *testing.T) {
	mock := &mockSNSPublisher{retErr: errors.New("publish error")}
	client := newTestSNSClient("arn:aws:sns:us-east-1:123456789012:alerts", mock)

	if err := client.Send("msg"); err == nil {
		t.Fatal("expected error from failed publish")
	}
}
