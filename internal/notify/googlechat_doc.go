// Package notify provides notification sender implementations for various
// alerting channels used by vaultpulse.
//
// GoogleChatClient sends secret expiry alerts to a Google Chat space via
// an incoming webhook URL. Configure the webhook in your Google Chat space
// settings and supply the URL in your vaultpulse configuration.
//
// Example usage:
//
//	client, err := notify.NewGoogleChatClient("https://chat.googleapis.com/v1/spaces/.../messages?key=...")
//	if err != nil {
//		log.Fatal(err)
//	}
//	client.Send("[CRITICAL] secret/db/password expires in 2 days")
package notify
