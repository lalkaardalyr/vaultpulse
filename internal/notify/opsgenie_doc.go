// Package notify provides integrations for sending secret-expiry
// alerts to external notification services.
//
// # OpsGenie
//
// OpsGenieClient sends alerts to OpsGenie using the Alerts v2 REST API.
// All messages are dispatched as P1 priority incidents.
//
// Usage:
//
//	client, err := notify.NewOpsGenieClient("<api-key>")
//	if err != nil {
//		log.Fatal(err)
//	}
//	if err := client.Send("secret /kv/db/password expires in 24h"); err != nil {
//		log.Printf("opsgenie alert failed: %v", err)
//	}
package notify
