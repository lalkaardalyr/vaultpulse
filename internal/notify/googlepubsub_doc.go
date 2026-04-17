// Package notify provides notification sender implementations for VaultPulse.
//
// # Google Cloud Pub/Sub
//
// GooglePubSubClient publishes alert messages to a Google Cloud Pub/Sub topic
// using the REST API with an API key for authentication.
//
// Usage:
//
//	client, err := notify.NewGooglePubSubClient("my-project", "my-topic", "AIza...")
//	if err != nil {
//		log.Fatal(err)
//	}
//	err = client.Send(ctx, "secret expiry warning")
package notify
