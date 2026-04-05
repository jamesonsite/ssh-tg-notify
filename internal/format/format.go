package format

import (
	"fmt"
	"time"

	"github.com/jamesonsite/ssh-tg-notify/internal/sshparse"
)

// LoginMessage builds a plain-text Telegram body for a successful SSH login.
func LoginMessage(hostLabel, hostname string, ev sshparse.Event) string {
	server := hostLabel
	if server == "" {
		server = hostname
	}
	return fmt.Sprintf(
		"SSH login\nServer: %s\nHost: %s\nUser: %s\nSource: %s\nAuth: %s\nTime: %s",
		server,
		hostname,
		ev.User,
		ev.Source,
		ev.AuthMethod,
		time.Now().UTC().Format(time.RFC3339),
	)
}
