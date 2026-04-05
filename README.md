# ssh-tg-notify

Small Linux agent that watches **successful SSH logins** and sends a **Telegram** message for each one. It is meant for homelab and small fleets: one static binary, `systemd` service, YAML config, no runtime beyond the OS.

## How it works

- **Primary input:** `journalctl` (follow mode), filtered to `sshd` by default (`_COMM=sshd`). Override `journal.args` if your distro logs differently.
- **Optional input:** plain-text auth log file (e.g. `/var/log/auth.log` or `/var/log/secure`).
- **Detection:** parses OpenSSH `Accepted …` lines (password, publickey, etc.).
- **Telegram:** outbound `sendMessage` only (no webhooks on the server).

## Publish to GitHub (first time)

GitHub requires a one-time login for the CLI (OAuth or a token). After that, creating the repo and pushing is automated:

```powershell
cd E:\GitHub\ssh-tg-notify
powershell -ExecutionPolicy Bypass -File .\scripts\publish-github.ps1
```

The script runs `gh auth login --web` if needed, creates `jamesonsite/ssh-tg-notify` if it does not exist, sets `origin`, and runs `git push -u origin main`. Alternatively, set `GH_TOKEN` (classic PAT with `repo` scope) and run the same script with no browser step.

## Quick setup

1. Create a bot with [@BotFather](https://t.me/BotFather), copy the **token**.
2. Start a chat with your bot and get your **chat id** (e.g. message [@userinfobot](https://t.me/userinfobot) or call `getUpdates` once).
3. On the server:

```bash
sudo mkdir -p /etc/sshnotify
sudo cp config.example.yaml /etc/sshnotify/config.yaml
sudo chmod 600 /etc/sshnotify/config.yaml
sudo nano /etc/sshnotify/config.yaml   # set bot_token and chat_id
```

4. Install the binary (pick one):

**From a release tarball** (after you publish a GitHub Release):

```bash
curl -fsSL -o /tmp/sshnotify.tgz "https://github.com/jamesonsite/ssh-tg-notify/releases/download/v0.1.0/sshnotify_v0.1.0_linux_amd64.tar.gz"
sudo tar -xzf /tmp/sshnotify.tgz -C /usr/local/bin sshnotify
sudo tar -xzf /tmp/sshnotify.tgz -C /tmp config.example.yaml sshnotify.service
sudo install -m 0644 /tmp/sshnotify.service /etc/systemd/system/sshnotify.service
sudo systemctl daemon-reload
sudo systemctl enable --now sshnotify
```

**From source** (needs Go 1.22+):

```bash
git clone https://github.com/jamesonsite/ssh-tg-notify.git
cd ssh-tg-notify
make build
sudo make install
sudo cp deploy/sshnotify.service /etc/systemd/system/sshnotify.service
sudo systemctl daemon-reload
sudo systemctl enable --now sshnotify
```

### RHEL / Rocky / Alma (optional auth log)

If you prefer file-based logs instead of or in addition to the journal:

```yaml
journal:
  enabled: false
authlog:
  enabled: true
  path: /var/log/secure
```

### Permissions

Reading the journal or `/var/log/auth.log` usually requires **root** (the shipped unit runs the binary as root). Tighten later with group membership (`systemd-journal`, `adm`) if you harden the service account.

## Configuration reference

| Field | Meaning |
| --- | --- |
| `telegram.bot_token` | Bot token from BotFather |
| `telegram.chat_id` | Destination chat (string; safe for large ids) |
| `server.label` | Optional label in messages (default: hostname) |
| `journal.enabled` | Follow journal (default **on** if omitted) |
| `journal.args` | Full `journalctl` argument list after the binary name |
| `authlog.enabled` | Tail a file (default **off**) |
| `authlog.path` | e.g. `/var/log/auth.log` or `/var/log/secure` |
| `notify.on_success` | Send on successful `Accepted` lines (default **on** if omitted) |
| `notify.dedupe_seconds` | Collapse duplicate user+source+method within this window |

## CLI

```text
sshnotify -config /etc/sshnotify/config.yaml
sshnotify -version
```

## Go module path

The module is [`github.com/jamesonsite/ssh-tg-notify`](https://github.com/jamesonsite/ssh-tg-notify). If you fork, run `go mod edit -module github.com/<you>/ssh-tg-notify` and replace import paths to match.

## License

MIT — see [LICENSE](LICENSE).
