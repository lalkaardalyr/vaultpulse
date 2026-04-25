// Package notify provides integrations with external notification services.
//
// JumpCloud client sends alert events to the JumpCloud Insights API,
// allowing VaultPulse to emit secret-expiry events into JumpCloud's
// audit log stream for centralised visibility.
//
// Usage:
//
//	client, err := notify.NewJumpCloudClient(apiKey, orgID)
//	if err != nil {
//		log.Fatal(err)
//	}
//	err = client.Send("secret expiry alert: path/to/secret")
package notify
