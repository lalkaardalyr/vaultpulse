// Package notify — awseventbridge.go
//
// EventBridgeClient sends alert notifications to AWS EventBridge by posting
// a structured JSON event to a configured endpoint URL. This is useful for
// routing VaultPulse secret-expiry alerts into AWS-native automation workflows
// such as Lambda functions, Step Functions, or cross-account event buses.
//
// Usage:
//
//	client, err := notify.NewEventBridgeClient(
//		"https://events.us-east-1.amazonaws.com/",
//		"vaultpulse",
//		"SecretExpiry",
//	)
//	if err != nil {
//		log.Fatal(err)
//	}
//	_ = client.Send("secret/db/password expires in 3 days")
package notify
