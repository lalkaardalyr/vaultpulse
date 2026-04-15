// Package notify provides notification sender implementations for VaultPulse.
//
// # Lark
//
// LarkClient sends alert messages to a Lark (Feishu) incoming webhook.
// It formats the message as a plain text card using the Lark Bot API.
//
// Usage:
//
//	client, err := notify.NewLarkClient("https://open.larksuite.com/open-apis/bot/v2/hook/<token>")
//	if err != nil {
//		log.Fatal(err)
//	}
//	if err := client.Send("Secret expiry warning: secret/db/password expires in 24h"); err != nil {
//		log.Println("lark notification failed:", err)
//	}
package notify
