// Package notify provides notification sender implementations for VaultPulse.
//
// ChatworkClient sends alert messages to a Chatwork room using the
// Chatwork REST API v2. Authentication is performed via the
// X-ChatWorkToken header.
//
// Usage:
//
//	client, err := notify.NewChatworkClient("your-api-token", "12345678")
//	if err != nil {
//		log.Fatal(err)
//	}
//	if err := client.Send("Secret expiry warning"); err != nil {
//		log.Println(err)
//	}
package notify
