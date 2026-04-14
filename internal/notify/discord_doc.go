// Package notify provides notification clients for sending alert messages
// to external services such as Slack, PagerDuty, OpsGenie, Teams, Discord,
// generic webhooks, and email.
//
// # Discord
//
// The DiscordClient sends messages to a Discord channel via an incoming
// webhook URL. Create the webhook in your Discord server's channel settings
// and pass the URL to NewDiscordClient.
//
// Example:
//
//	client, err := notify.NewDiscordClient("https://discord.com/api/webhooks/...")
//	if err != nil {
//		log.Fatal(err)
//	}
//	_ = client.Send("Secret expiry alert: my/secret expires in 24h")
package notify
