package dedupe

import (
	"sync"
	"time"
)

// Window suppresses duplicate keys within a duration.
type Window struct {
	mu     sync.Mutex
	window time.Duration
	last   map[string]time.Time
}

// NewWindow creates a deduplicator. window <= 0 defaults to 3s.
func NewWindow(window time.Duration) *Window {
	if window <= 0 {
		window = 3 * time.Second
	}
	return &Window{
		window: window,
		last:   make(map[string]time.Time),
	}
}

// ShouldSend returns false if key was seen recently.
func (w *Window) ShouldSend(key string, now time.Time) bool {
	w.mu.Lock()
	defer w.mu.Unlock()
	if prev, ok := w.last[key]; ok && now.Sub(prev) < w.window {
		return false
	}
	w.last[key] = now
	// crude GC: drop stale entries occasionally
	if len(w.last) > 2048 {
		cutoff := now.Add(-w.window * 4)
		for k, t := range w.last {
			if t.Before(cutoff) {
				delete(w.last, k)
			}
		}
	}
	return true
}
