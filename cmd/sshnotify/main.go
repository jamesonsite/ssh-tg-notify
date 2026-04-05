package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/james/ssh-tg-notify/internal/config"
	"github.com/james/ssh-tg-notify/internal/dedupe"
	"github.com/james/ssh-tg-notify/internal/follow"
	"github.com/james/ssh-tg-notify/internal/format"
	"github.com/james/ssh-tg-notify/internal/notify"
	"github.com/james/ssh-tg-notify/internal/sshparse"
)

// Set via -ldflags "-X main.version=..." at release build time.
var version = "dev"

func main() {
	configPath := flag.String("config", "/etc/sshnotify/config.yaml", "path to YAML config")
	showVer := flag.Bool("version", false, "print version and exit")
	flag.Parse()

	if *showVer {
		fmt.Println(version)
		return
	}

	cfg, err := config.Load(*configPath)
	if err != nil {
		log.Fatalf("config: %v", err)
	}

	hostname, err := os.Hostname()
	if err != nil {
		hostname = "unknown"
	}
	label := strings.TrimSpace(cfg.Server.Label)

	tg := &notify.Telegram{
		BotToken: cfg.Telegram.BotToken,
		ChatID:   strings.TrimSpace(cfg.Telegram.ChatID),
	}

	deduper := dedupe.NewWindow(time.Duration(cfg.Notify.DedupeSec) * time.Second)

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	lines := make(chan string, 256)

	if cfg.JournalEnabled() {
		args := append([]string{}, cfg.Journal.Args...)
		go func() {
			for {
				if ctx.Err() != nil {
					return
				}
				log.Printf("following journal via %s %v", cfg.Journal.Binary, args)
				err := follow.Journal(ctx, cfg.Journal.Binary, args, lines)
				if ctx.Err() != nil {
					return
				}
				if err != nil {
					log.Printf("journal follower exited: %v (retry in 5s)", err)
				} else {
					log.Printf("journal follower ended (retry in 5s)")
				}
				time.Sleep(5 * time.Second)
			}
		}()
	}
	if cfg.AuthLogEnabled() {
		go func() {
			log.Printf("tailing auth log %s", cfg.AuthLog.Path)
			if err := follow.AuthLog(ctx, cfg.AuthLog.Path, lines); err != nil && ctx.Err() == nil {
				log.Printf("auth log follower exited: %v", err)
			}
		}()
	}

	for {
		select {
		case <-ctx.Done():
			log.Println("shutting down")
			return
		case line := <-lines:
			if !cfg.Notify.OnSuccess {
				continue
			}
			ev, ok := sshparse.ParseLine(line)
			if !ok {
				continue
			}
			key := ev.User + "|" + ev.Source + "|" + ev.AuthMethod
			if !deduper.ShouldSend(key, time.Now()) {
				continue
			}
			msg := format.LoginMessage(label, hostname, *ev)
			sendCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
			err := tg.SendMessage(sendCtx, msg)
			cancel()
			if err != nil {
				log.Printf("telegram: %v", err)
			}
		}
	}
}
