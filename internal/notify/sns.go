package notify

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sns"
)

// SNSClient sends alert messages to an AWS SNS topic.
type SNSClient struct {
	topicARN string
	client   *sns.Client
}

type snsPublisher interface {
	Publish(ctx context.Context, input *sns.PublishInput, opts ...func(*sns.Options)) (*sns.PublishOutput, error)
}

// snsClientWrapper wraps the real SNS client for testability.
type snsClientWrapper struct {
	pub snsPublisher
	topicARN string
}

// NewSNSClient creates a new SNSClient for the given topic ARN.
// AWS credentials and region are resolved from the environment.
func NewSNSClient(topicARN string) (*SNSClient, error) {
	if topicARN == "" {
		return nil, fmt.Errorf("sns: topic ARN must not be empty")
	}
	cfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		return nil, fmt.Errorf("sns: failed to load AWS config: %w", err)
	}
	return &SNSClient{
		topicARN: topicARN,
		client:   sns.NewFromConfig(cfg),
	}, nil
}

// Send publishes the message to the configured SNS topic.
func (c *SNSClient) Send(msg string) error {
	payload, err := json.Marshal(map[string]string{"message": msg})
	if err != nil {
		return fmt.Errorf("sns: failed to marshal message: %w", err)
	}
	_, err = c.client.Publish(context.Background(), &sns.PublishInput{
		TopicArn: aws.String(c.topicARN),
		Message:  aws.String(string(payload)),
	})
	if err != nil {
		return fmt.Errorf("sns: publish failed: %w", err)
	}
	return nil
}
