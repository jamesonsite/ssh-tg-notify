package notify

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestTelegram_SendMessage(t *testing.T) {
	var got map[string]string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := json.NewDecoder(r.Body).Decode(&got); err != nil {
			t.Fatal(err)
		}
		_ = json.NewEncoder(w).Encode(map[string]any{"ok": true, "result": map[string]any{}})
	}))
	t.Cleanup(srv.Close)

	tg := &Telegram{
		HTTP:     srv.Client(),
		BaseURL:  srv.URL,
		BotToken: "TEST:TOKEN",
		ChatID:   "42",
	}
	if err := tg.SendMessage(context.Background(), "hello"); err != nil {
		t.Fatal(err)
	}
	if got["chat_id"] != "42" || got["text"] != "hello" {
		t.Fatalf("unexpected payload: %v", got)
	}
}
