package sshparse

import (
	"regexp"
	"strings"
)

// Event describes a successful SSH authentication we care about.
type Event struct {
	AuthMethod string
	User       string
	Source     string // IP or hostname as logged
	Raw        string
}

// acceptedRE matches OpenSSH "Accepted …" lines (password, publickey, etc.).
// IPv6 addresses may contain colons; we capture up to " port N".
var acceptedRE = regexp.MustCompile(`(?i)(?:sshd(?:-session)?\[\d+\]:\s*)?Accepted\s+(\S+)\s+for\s+(\S+)\s+from\s+(.+?)\s+port\s+(\d+)`)

// ParseLine returns a successful-login event if the line matches.
func ParseLine(line string) (*Event, bool) {
	s := strings.TrimSpace(line)
	if s == "" {
		return nil, false
	}
	m := acceptedRE.FindStringSubmatch(s)
	if m == nil {
		return nil, false
	}
	return &Event{
		AuthMethod: m[1],
		User:       m[2],
		Source:     strings.TrimSpace(m[3]),
		Raw:        s,
	}, true
}
