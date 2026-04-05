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

1. Create a bot with [@BotFather](https://t.me/BotFather), copy the **token** (looks like `123456789:AAH...`).
2. Open a chat with your bot, send any message, then get your **chat id** (e.g. [@userinfobot](https://t.me/userinfobot) or Telegram `getUpdates` once). Use the **numeric id** (and for groups, often a negative id like `-100...`).

### Config file (not blank)

You never start from an empty file. Use the **template** [`config.example.yaml`](config.example.yaml) from this repo: it already lists every section with sane defaults. After copying it to `/etc/sshnotify/config.yaml`, you normally change **only**:

| Field | What to put |
| --- | --- |
| `telegram.bot_token` | Paste the full token from BotFather (keep the quotes). |
| `telegram.chat_id` | Your id as a **string** in quotes, e.g. `"123456789"` (avoids large-number issues). |
| `server.label` | Optional friendly name in Telegram messages; leave `""` to use the server hostname. |

Everything else can stay as in the example unless you are on RHEL-style logs (see below).

Pick **one** install path on the server.

### A — Prebuilt binary (no Go on the server)

The running agent is a **single executable**; end servers do **not** need Go installed if you use a [release](https://github.com/jamesonsite/ssh-tg-notify/releases) tarball (build the release once via GitHub Actions, or build on your PC and upload the binary).

Replace `v0.1.0` with the tag you published:

```bash
curl -fsSL -o /tmp/sshnotify.tgz "https://github.com/jamesonsite/ssh-tg-notify/releases/download/v0.1.0/sshnotify_v0.1.0_linux_amd64.tar.gz"
tar -xzf /tmp/sshnotify.tgz -C /tmp
sudo install -m 0755 /tmp/sshnotify /usr/local/bin/sshnotify
sudo mkdir -p /etc/sshnotify
sudo install -m 0600 /tmp/config.example.yaml /etc/sshnotify/config.yaml
sudo nano /etc/sshnotify/config.yaml   # set bot_token and chat_id
sudo install -m 0644 /tmp/sshnotify.service /etc/systemd/system/sshnotify.service
sudo systemctl daemon-reload
sudo systemctl enable --now sshnotify
```

`curl` / `tar` are standard; **no** `git` or `golang-go`.

### B — Build on the server (needs Go once)

`git` downloads the source. **`golang-go`** is Debian/Ubuntu’s **Go compiler package**—it is only used to **compile** `sshnotify` on that machine. After `go build`, the binary does not need Go at runtime; you could remove `golang-go` afterward if you wanted (uncommon).

```bash
sudo apt update
sudo apt install -y git golang-go
git clone https://github.com/jamesonsite/ssh-tg-notify.git
cd ssh-tg-notify
CGO_ENABLED=0 go build -trimpath -ldflags "-s -w" -o sshnotify ./cmd/sshnotify
sudo install -m 0755 sshnotify /usr/local/bin/sshnotify
sudo mkdir -p /etc/sshnotify
sudo install -m 0600 config.example.yaml /etc/sshnotify/config.yaml
sudo nano /etc/sshnotify/config.yaml   # set bot_token and chat_id
sudo cp deploy/sshnotify.service /etc/systemd/system/sshnotify.service
sudo systemctl daemon-reload
sudo systemctl enable --now sshnotify
```

If you do not want Go on the VPS, use path **A** or build the Linux binary on Windows (`GOOS=linux GOARCH=amd64 go build ...`) and copy `sshnotify` over with `scp`.

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
