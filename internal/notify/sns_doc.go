// Package notify provides notification sender implementations for VaultPulse.
//
// # SNS
//
// SNSClient publishes alert messages to an AWS Simple Notification Service (SNS)
// topic. It resolves AWS credentials and region automatically from the standard
// AWS credential chain (environment variables, shared credentials file, IAM
// instance profile, etc.).
//
// Usage:
//
//	client, err := notify.NewSNSClient("arn:aws:sns:us-east-1:123456789012:vaultpulse-alerts")
//	if err != nil {
//		log.Fatal(err)
//	}
//	if err := client.Send("secret expiring soon"); err != nil {
//		log.Println("sns alert failed:", err)
//	}
package notify
