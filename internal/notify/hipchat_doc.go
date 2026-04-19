// Package notify provides notification sender implementations for vaultpulse.
//
// # HipChat
//
// HipChatClient delivers alert messages to a HipChat room using the
// HipChat v2 REST API.
//
// Usage:
//
//	client, err := notify.NewHipChatClient(
//		"https://api.hipchat.com/v2/room/123/notification",
//		"my-auth-token",
//	)
//	if err != nil {
//		log.Fatal(err)
//	}
//	err = client.Send("Vault secret expiring soon")
package notify
