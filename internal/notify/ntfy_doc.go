// Package notify provides notification sender implementations for VaultPulse.
//
// # Ntfy
//
// NtfyClient sends alert messages to an ntfy topic.
// ntfy (https://ntfy.sh) is a simple, open-source push notification service
// that supports self-hosting.
//
// Example usage:
//
//	client, err := notify.NewNtfyClient("https://ntfy.sh", "vaultpulse-alerts")
//	if err != nil {
//		log.Fatal(err)
//	}
//	if err := client.Send("Secret expiring in 24h"); err != nil {
//		log.Println(err)
//	}
package notify
