package notify

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

const welcomeText = `This bot is used by ssh-tg-notify on your Linux servers to send SSH login alerts.

To set up your own bot token (for alerts on your machines):
• Tap “Open BotFather” below
• Send /newbot and follow the prompts
• Put that token in ssh-tg-notify config as telegram.bot_token
• Your numeric chat id goes in telegram.chat_id (message @userinfobot for your id)

Then install ssh-tg-notify on each server and run the systemd service.`

// RunWelcomePoller long-polls getUpdates and replies to /start with setup help + BotFather link.
// Blocks until ctx is cancelled.
func (t *Telegram) RunWelcomePoller(ctx context.Context) {
	longClient := &http.Client{Timeout: 55 * time.Second}
	offset := 0
	for {
		if err := ctx.Err(); err != nil {
			return
		}
		updates, err := t.getUpdates(ctx, longClient, offset)
		if err != nil {
			if ctx.Err() != nil {
				return
			}
			time.Sleep(3 * time.Second)
			continue
		}
		for _, u := range updates {
			if u.UpdateID >= offset {
				offset = u.UpdateID + 1
			}
			if u.Message == nil || !isStartCommand(u.Message.Text) {
				continue
			}
			chatID := strconv.FormatInt(u.Message.Chat.ID, 10)
			sendCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
			err := t.sendWelcome(sendCtx, chatID)
			cancel()
			if err != nil {
				log.Printf("welcome reply: %v", err)
			}
		}
	}
}

func isStartCommand(text string) bool {
	text = strings.TrimSpace(text)
	if text == "/start" {
		return true
	}
	return strings.HasPrefix(text, "/start ") || strings.HasPrefix(text, "/start@")
}

type update struct {
	UpdateID int `json:"update_id"`
	Message  *struct {
		MessageID int64 `json:"message_id"`
		Text      string `json:"text"`
		Chat      struct {
			ID int64 `json:"id"`
		} `json:"chat"`
	} `json:"message"`
}

func (t *Telegram) getUpdates(ctx context.Context, client *http.Client, offset int) ([]update, error) {
	q := url.Values{}
	if offset > 0 {
		q.Set("offset", strconv.Itoa(offset))
	}
	q.Set("timeout", "50")
	u := t.endpoint("getUpdates") + "?" + q.Encode()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
	if err != nil {
		return nil, err
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	b, err := io.ReadAll(io.LimitReader(resp.Body, 4<<20))
	if err != nil {
		return nil, err
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("getUpdates HTTP %d: %s", resp.StatusCode, strings.TrimSpace(string(b)))
	}
	var out struct {
		OK     bool     `json:"ok"`
		Result []update `json:"result"`
	}
	if err := json.Unmarshal(b, &out); err != nil {
		return nil, err
	}
	if !out.OK {
		return nil, fmt.Errorf("getUpdates not ok")
	}
	return out.Result, nil
}

func (t *Telegram) sendWelcome(ctx context.Context, chatID string) error {
	payload := map[string]interface{}{
		"chat_id": chatID,
		"text":    welcomeText,
		"reply_markup": map[string]interface{}{
			"inline_keyboard": [][]map[string]string{
				{{"text": "Open BotFather", "url": "https://t.me/BotFather"}},
			},
		},
	}
	body, err := json.Marshal(payload)
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
		return fmt.Errorf("sendMessage HTTP %d: %s", resp.StatusCode, strings.TrimSpace(string(b)))
	}
	var api struct {
		OK          bool   `json:"ok"`
		Description string `json:"description"`
	}
	if err := json.Unmarshal(b, &api); err != nil {
		return err
	}
	if !api.OK {
		return fmt.Errorf("telegram api: %s", api.Description)
	}
	return nil
}
