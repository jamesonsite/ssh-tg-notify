package config

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

// Config is loaded from YAML (e.g. /etc/sshnotify/config.yaml).
type Config struct {
	Telegram Telegram `yaml:"telegram"`
	Server   Server   `yaml:"server"`
	Journal  Journal  `yaml:"journal"`
	AuthLog  AuthLog  `yaml:"authlog"`
	Notify   Notify   `yaml:"notify"`
}

type Telegram struct {
	BotToken string `yaml:"bot_token"`
	ChatID   string `yaml:"chat_id"`
}

type Server struct {
	Label string `yaml:"label"` // optional; defaults to hostname
}

type Journal struct {
	Enabled *bool  `yaml:"enabled"` // default true when omitted
	Binary  string `yaml:"binary"`  // default: journalctl
	// Args are passed to journalctl after the binary name. If empty, a safe default is used.
	Args []string `yaml:"args"`
}

type AuthLog struct {
	Enabled *bool `yaml:"enabled"` // default false when omitted
	Path    string `yaml:"path"`
}

type Notify struct {
	OnSuccess bool `yaml:"on_success"` // default true when omitted (see Defaults)
	DedupeSec int  `yaml:"dedupe_seconds"`
}

// JournalEnabled reports whether journal following is on (default true).
func (c *Config) JournalEnabled() bool {
	if c.Journal.Enabled == nil {
		return true
	}
	return *c.Journal.Enabled
}

// AuthLogEnabled reports whether auth.log tailing is on (default false).
func (c *Config) AuthLogEnabled() bool {
	if c.AuthLog.Enabled == nil {
		return false
	}
	return *c.AuthLog.Enabled
}

// Defaults fills derived defaults after YAML load.
func (c *Config) Defaults() {
	if c.Journal.Binary == "" {
		c.Journal.Binary = "journalctl"
	}
	if len(c.Journal.Args) == 0 {
		c.Journal.Args = []string{"-f", "-n", "0", "-o", "cat", "_COMM=sshd"}
	}
	if c.AuthLog.Path == "" {
		c.AuthLog.Path = "/var/log/auth.log"
	}
	if c.Notify.DedupeSec == 0 {
		c.Notify.DedupeSec = 3
	}
}

// Validate checks required fields and source configuration.
func (c *Config) Validate() error {
	c.Defaults()
	if strings.TrimSpace(c.Telegram.BotToken) == "" {
		return errors.New("telegram.bot_token is required")
	}
	if strings.TrimSpace(c.Telegram.ChatID) == "" {
		return errors.New("telegram.chat_id is required")
	}
	if !c.JournalEnabled() && !c.AuthLogEnabled() {
		return errors.New("enable journal and/or authlog (both off)")
	}
	return nil
}

// Load reads and parses a YAML config file.
func Load(path string) (*Config, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read config: %w", err)
	}
	var raw struct {
		Telegram Telegram `yaml:"telegram"`
		Server   Server   `yaml:"server"`
		Journal  Journal  `yaml:"journal"`
		AuthLog  AuthLog  `yaml:"authlog"`
		Notify   *struct {
			OnSuccess *bool `yaml:"on_success"`
			DedupeSec int   `yaml:"dedupe_seconds"`
		} `yaml:"notify"`
	}
	if err := yaml.Unmarshal(b, &raw); err != nil {
		return nil, fmt.Errorf("parse yaml: %w", err)
	}
	c := Config{
		Telegram: raw.Telegram,
		Server:   raw.Server,
		Journal:  raw.Journal,
		AuthLog:  raw.AuthLog,
	}
	if raw.Notify != nil {
		if raw.Notify.OnSuccess != nil {
			c.Notify.OnSuccess = *raw.Notify.OnSuccess
		} else {
			c.Notify.OnSuccess = true
		}
		c.Notify.DedupeSec = raw.Notify.DedupeSec
	} else {
		c.Notify.OnSuccess = true
	}
	if err := c.Validate(); err != nil {
		return nil, err
	}
	return &c, nil
}
