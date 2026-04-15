// Package notify provides notification sender implementations for vaultpulse.
//
// # Rocket.Chat
//
// RocketChatClient delivers alert messages via a Rocket.Chat incoming webhook.
// Create an incoming webhook integration in your Rocket.Chat administration
// panel and supply the generated URL to NewRocketChatClient.
//
// Example:
//
//	client, err := notify.NewRocketChatClient("https://chat.example.com/hooks/TOKEN")
//	if err != nil {
//		log.Fatal(err)
//	}
//	err = client.Send("[CRITICAL] secret/db/password expires in 2 days")
package notify
