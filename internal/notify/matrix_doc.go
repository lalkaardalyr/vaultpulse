// Package notify provides integrations for sending alert messages to
// various external services.
//
// # Matrix
//
// MatrixClient delivers alerts to a Matrix room using the Matrix
// Client-Server API (v3). Authentication is performed via a Bearer
// access token obtained from the homeserver.
//
// Example:
//
//	client, err := notify.NewMatrixClient(
//		"https://matrix.org",
//		"syt_access_token",
//		"!roomid:matrix.org",
//	)
//	if err != nil {
//		log.Fatal(err)
//	}
//	err = client.Send("Vault secret expiring in 24 hours")
package notify
