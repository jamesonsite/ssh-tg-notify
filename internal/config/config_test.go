package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoad_Minimal(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, "cfg.yaml")
	content := `
telegram:
  bot_token: "1:token"
  chat_id: "99"
`
	if err := os.WriteFile(p, []byte(content), 0600); err != nil {
		t.Fatal(err)
	}
	c, err := Load(p)
	if err != nil {
		t.Fatal(err)
	}
	if !c.JournalEnabled() || c.AuthLogEnabled() {
		t.Fatalf("unexpected defaults: journal=%v authlog=%v", c.JournalEnabled(), c.AuthLogEnabled())
	}
	if !c.Notify.OnSuccess {
		t.Fatal("on_success should default true")
	}
}

func TestLoad_AuthLogOnly(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, "cfg.yaml")
	content := `
telegram:
  bot_token: "1:token"
  chat_id: "99"
journal:
  enabled: false
authlog:
  enabled: true
  path: /var/log/secure
`
	if err := os.WriteFile(p, []byte(content), 0600); err != nil {
		t.Fatal(err)
	}
	c, err := Load(p)
	if err != nil {
		t.Fatal(err)
	}
	if c.JournalEnabled() || !c.AuthLogEnabled() {
		t.Fatal("expected auth log only")
	}
}
