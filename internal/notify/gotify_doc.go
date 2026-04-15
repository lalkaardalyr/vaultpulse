// Package notify provides integrations with external notification services.
//
// # Gotify
//
// GotifyClient delivers alert messages to a self-hosted Gotify server
// (https://gotify.net). Each message is sent as a POST request to the
// /message endpoint authenticated via an application token.
//
// Example usage:
//
//	client, err := notify.NewGotifyClient(
//		"https://gotify.example.com",
//		"AppToken123",
//		5,
//	)
//	if err != nil {
//		log.Fatal(err)
//	}
//	if err := client.Send("secret expiry warning"); err != nil {
//		log.Println(err)
//	}
package notify
