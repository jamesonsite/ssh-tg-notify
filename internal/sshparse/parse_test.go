package sshparse

import "testing"

func TestParseLine_Accepted(t *testing.T) {
	tests := []struct {
		line string
		want Event
	}{
		{
			line: "Accepted publickey for alice from 203.0.113.10 port 22 ssh2",
			want: Event{AuthMethod: "publickey", User: "alice", Source: "203.0.113.10"},
		},
		{
			line: "Apr  5 10:00:00 web sshd[12345]: Accepted password for bob from 10.0.0.5 port 54321 ssh2",
			want: Event{AuthMethod: "password", User: "bob", Source: "10.0.0.5"},
		},
		{
			line: "sshd[999]: Accepted keyboard-interactive for root from 192.0.2.1 port 22 ssh2",
			want: Event{AuthMethod: "keyboard-interactive", User: "root", Source: "192.0.2.1"},
		},
		{
			line: "Accepted publickey for deploy from 2001:db8::1 port 22 ssh2",
			want: Event{AuthMethod: "publickey", User: "deploy", Source: "2001:db8::1"},
		},
		{
			line: "sshd-session[42]: Accepted publickey for git from 172.16.0.3 port 22 ssh2",
			want: Event{AuthMethod: "publickey", User: "git", Source: "172.16.0.3"},
		},
	}
	for _, tt := range tests {
		ev, ok := ParseLine(tt.line)
		if !ok {
			t.Fatalf("expected match: %q", tt.line)
		}
		if ev.AuthMethod != tt.want.AuthMethod || ev.User != tt.want.User || ev.Source != tt.want.Source {
			t.Fatalf("line %q\ngot  %+v\nwant %+v", tt.line, *ev, tt.want)
		}
	}
}

func TestParseLine_NoMatch(t *testing.T) {
	lines := []string{
		"",
		"Failed password for root from 1.2.3.4 port 22 ssh2",
		"Connection closed by 1.2.3.4",
		"Invalid user admin from 1.2.3.4",
	}
	for _, line := range lines {
		if _, ok := ParseLine(line); ok {
			t.Fatalf("should not match: %q", line)
		}
	}
}
