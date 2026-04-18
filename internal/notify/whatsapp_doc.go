// Package notify provides notification clients for various alerting platforms.
//
// WhatsAppClient sends alert messages via the WhatsApp Business Cloud API.
// It requires a valid access token, a phone number ID registered with the
// WhatsApp Business platform, and a recipient phone number in E.164 format.
//
// Example usage:
//
//	client, err := notify.NewWhatsAppClient(token, phoneID, "+15550001234")
//	if err != nil {
//		log.Fatal(err)
//	}
//	err = client.Send("Vault secret expiring in 24h: secret/myapp/db")
package notify
