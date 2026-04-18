// Package notify provides notification clients for various alerting platforms.
//
// # Clickatell
//
// NewClickatellClient sends SMS alerts via the Clickatell HTTP API.
//
// Required fields:
//   - APIKey: your Clickatell API key
//   - To:     recipient phone number in international format (e.g. +12025550100)
//
// Example:
//
//	client, err := notify.NewClickatellClient("my-api-key", "+12025550100")
//	if err != nil {
//	    log.Fatal(err)
//	}
//	err = client.Send(ctx, "Vault secret expiring in 24h")
package notify
