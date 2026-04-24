// Package notify provides notification clients for various alerting
// platforms used by vaultpulse to deliver secret-expiry alerts.
//
// # PagerTree
//
// PagerTreeClient posts alert payloads to the PagerTree Incoming
// Integration endpoint (https://api.pagertree.com/integration/<id>).
//
// Usage:
//
//	c, err := notify.NewPagerTreeClient("your-integration-id")
//	if err != nil {
//		log.Fatal(err)
//	}
//	if err := c.Send("secret/db/password expires in 24h"); err != nil {
//		log.Println("pagertree send:", err)
//	}
package notify
