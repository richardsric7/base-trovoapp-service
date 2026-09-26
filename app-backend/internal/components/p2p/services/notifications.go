package p2p

import (
	usersDB "trovo-wallet-api/internal/components/users/db"
	"trovo-wallet-api/internal/sharedconfig"
)

// NotifyUsername sends a direct push notification to a username, matching
// app-backend's existing direct-call notification pattern (there is no
// event-bus/rules-engine anywhere in app-backend - Plan Section 67).
// Failures are swallowed by User.SendPushMessage itself (it already logs);
// a notification failure must never block an order-lifecycle transition.
func NotifyUsername(gc *sharedconfig.GlobalConfig, username, title, body string, dataPayload map[string]string) {
	if username == "" {
		return
	}
	user, err := usersDB.GetUser(username, gc.DB, gc)
	if err != nil || user.Username == "" {
		return
	}
	user.SendPushMessage(title, body, "", dataPayload, gc)
}
