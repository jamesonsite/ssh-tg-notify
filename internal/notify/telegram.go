package notify

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// Telegram sends outbound Bot API messages.
type Telegram struct {
	HTTP     *http.Client
	BotToken string
	ChatID   string
	// BaseURL is optional; default https://api.telegram.org (set in tests).
	BaseURL string
}

func (t *Telegram) endpoint(method string) string {
	base := strings.TrimSuffix(strings.TrimSpace(t.BaseURL), "/")
	if base == "" {
		base = "https://api.telegram.org"
	}
	// Token is placed in the path per Telegram docs; do not URL-encode it (":" must stay literal).
	return base + "/bot" + t.botToken() + "/" + method
}

func (t *Telegram) botToken() string {
	return strings.TrimSpace(t.BotToken)
}

func (t *Telegram) client() *http.Client {
	if t.HTTP != nil {
		return t.HTTP
	}
	return &http.Client{Timeout: 25 * time.Second}
}

// SendMessage posts text to the configured chat.
func (t *Telegram) SendMessage(ctx context.Context, text string) error {
	body, err := json.Marshal(map[string]string{
		"chat_id": t.ChatID,
		"text":    text,
	})
	if err != nil {
		return err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, t.endpoint("sendMessage"), bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := t.client().Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	b, _ := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("telegram HTTP %d: %s", resp.StatusCode, strings.TrimSpace(string(b)))
	}
	var api struct {
		OK          bool   `json:"ok"`
		Description string `json:"description"`
	}
	if err := json.Unmarshal(b, &api); err != nil {
		return fmt.Errorf("telegram response: %w", err)
	}
	if !api.OK {
		return fmt.Errorf("telegram api: %s", api.Description)
	}
	return nil
}
