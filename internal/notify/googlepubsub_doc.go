// Package notify provides notification clients for various platforms.
//
// # Google Pub/Sub
//
// GooglePubSubClient publishes alert messages to a Google Cloud Pub/Sub topic
// using the REST API with an API key for authentication.
//
// Usage:
//
//	client, err := notify.NewGooglePubSubClient("my-project", "my-topic", "api-key")
//	if err != nil {
//	    log.Fatal(err)
//	}
//	err = client.Send(ctx, "secret expiring soon")
package notify
