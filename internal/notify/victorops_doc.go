// Package notify provides notification sender implementations for various
// alerting platforms used by VaultPulse.
//
// VictorOps (Splunk On-Call) integration
//
// The VictorOpsClient sends CRITICAL alert payloads to a VictorOps REST
// endpoint URL. Configure the webhook URL via the vaultpulse config file
// under notify.victorops.webhook_url.
//
// Example usage:
//
//	client, err := notify.NewVictorOpsClient(cfg.Notify.VictorOps.WebhookURL)
//	if err != nil {
//		log.Fatal(err)
//	}
//	_ = client.Send("secret/db expires in 24 hours")
package notify
