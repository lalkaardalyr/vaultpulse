package notify

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

// SQSClient sends alert messages to an Amazon SQS queue.
type SQSClient struct {
	queueURL string
	client   sqsSender
}

type sqsSender interface {
	SendMessage(ctx context.Context, params *sqs.SendMessageInput, optFns ...func(*sqs.Options)) (*sqs.SendMessageOutput, error)
}

// NewSQSClient creates a new SQSClient. queueURL must be a valid SQS queue URL.
func NewSQSClient(queueURL string) (*SQSClient, error) {
	if queueURL == "" {
		return nil, errors.New("sqs: queue URL must not be empty")
	}

	cfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		return nil, fmt.Errorf("sqs: failed to load AWS config: %w", err)
	}

	return &SQSClient{
		queueURL: queueURL,
		client:   sqs.NewFromConfig(cfg),
	}, nil
}

func newSQSClientWithSender(queueURL string, sender sqsSender) (*SQSClient, error) {
	if queueURL == "" {
		return nil, errors.New("sqs: queue URL must not be empty")
	}
	return &SQSClient{queueURL: queueURL, client: sender}, nil
}

// Send publishes the message to the configured SQS queue as a JSON body.
func (c *SQSClient) Send(message string) error {
	payload, err := json.Marshal(map[string]string{"message": message})
	if err != nil {
		return fmt.Errorf("sqs: failed to marshal message: %w", err)
	}

	_, err = c.client.SendMessage(context.Background(), &sqs.SendMessageInput{
		QueueUrl:    aws.String(c.queueURL),
		MessageBody: aws.String(string(payload)),
	})
	if err != nil {
		return fmt.Errorf("sqs: failed to send message: %w", err)
	}
	return nil
}
