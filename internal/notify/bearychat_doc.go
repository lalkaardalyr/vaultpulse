// Package notify provides notification sender implementations for vaultpulse.
//
// BearyChatClient sends alert messages to a BearyChat incoming webhook.
// BearyChat is a team collaboration tool popular in Asia.
//
// Usage:
//
//	client, err := notify.NewBearyChatClient("https://hook.bearychat.com/...")
//	if err != nil {
//	    log.Fatal(err)
//	}
//	if err := client.Send("secret expiry alert"); err != nil {
//	    log.Println(err)
//	}
package notify
