package notify

import "testing"

func TestIsStartCommand(t *testing.T) {
	yes := []string{"/start", "/start ", "/start payload", "/start@mybot", "/start@mybot deep"}
	no := []string{"", "start", "/stop", "/star", "hello"}
	for _, s := range yes {
		if !isStartCommand(s) {
			t.Fatalf("expected start: %q", s)
		}
	}
	for _, s := range no {
		if isStartCommand(s) {
			t.Fatalf("expected not start: %q", s)
		}
	}
}
