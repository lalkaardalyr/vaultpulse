// Package notify provides notification clients for various alerting platforms.
//
// ZulipClient sends VaultPulse secret expiry alerts to a Zulip stream using
// the Zulip REST API with HTTP Basic Auth (bot email + API key).
//
// Usage:
//
//	client, err := notify.NewZulipClient(
//		"https://yourorg.zulipchat.com",
//		"vaultpulse-bot@yourorg.zulipchat.com",
//		"<api-key>",
//		"alerts",
//		"vault-expiry",
//	)
//	if err != nil {
//		log.Fatal(err)
//	}
//	err = client.Send("Secret expiring in 24h: secret/prod/db")
package notify
