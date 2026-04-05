# ssh-tg-notify

Linux daemon that detects **successful SSH logins** and sends a **Telegram** message for each one. Ships as a small static binary plus a `systemd` unit and a YAML config—no language runtime on the server after install.

## Requirements

- Linux with **systemd** (typical for most VPS images).
- **journald** (default) or a readable auth log file (e.g. `/var/log/auth.log`, `/var/log/secure`).
- Outbound HTTPS to `api.telegram.org`.
- Usually runs as **root** so it can read the journal or auth logs (see [Permissions](#permissions)).

## How it works

- Follows **`journalctl`** for `sshd` by default, or optionally tails a log file.
- Parses OpenSSH **`Accepted …`** lines (password, publickey, etc.); not raw TCP connects.
- Sends alerts with the Telegram Bot API (`sendMessage` only).

## 1. Create the Telegram bot

1. Talk to [@BotFather](https://t.me/BotFather), create a bot, copy the **token** (format `123456789:AAH…`).
2. Start a chat with your bot, send any message, then obtain your **chat id** (e.g. [@userinfobot](https://t.me/userinfobot) or the `getUpdates` API). Use the numeric id; group ids are often negative (e.g. `-100…`).

## 2. Configuration

Copy the template from this repository ([`config.example.yaml`](config.example.yaml)) to `/etc/sshnotify/config.yaml`. It is **not** an empty file: it includes all sections and defaults.

**You normally only edit:**

| Field | Value |
| --- | --- |
| `telegram.bot_token` | Full token from BotFather (keep YAML quotes). |
| `telegram.chat_id` | Chat id as a **string**, e.g. `"123456789"`. |
| `server.label` | Optional display name in messages; leave `""` to use the machine hostname. |

Leave the rest unchanged unless you need file-based logs (see [RHEL / Rocky / Alma](#rhel--rocky--alma-optional-auth-log)).

## 3. Install

Use **either** a release tarball **or** build from source.

### Option A — Release binary (recommended)

No Go or Git required on the server—only `curl` and `tar`.

1. Open [Releases](https://github.com/jamesonsite/ssh-tg-notify/releases), pick a version, download **`sshnotify_<version>_linux_amd64.tar.gz`** (or `arm64` on ARM).

2. On the server (adjust URL and filename to match the release you chose):

```bash
curl -fsSL -o /tmp/sshnotify.tgz "https://github.com/jamesonsite/ssh-tg-notify/releases/download/v0.1.0/sshnotify_v0.1.0_linux_amd64.tar.gz"
tar -xzf /tmp/sshnotify.tgz -C /tmp
sudo install -m 0755 /tmp/sshnotify /usr/local/bin/sshnotify
sudo mkdir -p /etc/sshnotify
sudo install -m 0600 /tmp/config.example.yaml /etc/sshnotify/config.yaml
sudo nano /etc/sshnotify/config.yaml
sudo install -m 0644 /tmp/sshnotify.service /etc/systemd/system/sshnotify.service
sudo systemctl daemon-reload
sudo systemctl enable --now sshnotify
```

### Option B — Build from source on the server

Install **Git** to clone the repo and **Go** only to compile; the running service uses the compiled binary, not Go at runtime.

```bash
sudo apt update
sudo apt install -y git golang-go
git clone https://github.com/jamesonsite/ssh-tg-notify.git
cd ssh-tg-notify
CGO_ENABLED=0 go build -trimpath -ldflags "-s -w" -o sshnotify ./cmd/sshnotify
sudo install -m 0755 sshnotify /usr/local/bin/sshnotify
sudo mkdir -p /etc/sshnotify
sudo install -m 0600 config.example.yaml /etc/sshnotify/config.yaml
sudo nano /etc/sshnotify/config.yaml
sudo cp deploy/sshnotify.service /etc/systemd/system/sshnotify.service
sudo systemctl daemon-reload
sudo systemctl enable --now sshnotify
```

On non-Debian systems, install `git` and a recent Go toolchain with that distro’s packages or from [go.dev](https://go.dev/dl/).

## RHEL / Rocky / Alma (optional auth log)

If you rely on `/var/log/secure` instead of the journal for SSH lines:

```yaml
journal:
  enabled: false
authlog:
  enabled: true
  path: /var/log/secure
```

## Permissions

The packaged `systemd` unit runs `sshnotify` as root so it can read **journald** or **auth** logs. To harden later, you can run as a dedicated user with membership in `systemd-journal` or `adm`, depending on your distro.

## Configuration reference

| Field | Meaning |
| --- | --- |
| `telegram.bot_token` | Bot token from BotFather |
| `telegram.chat_id` | Destination chat (string) |
| `server.label` | Optional label (default: hostname) |
| `journal.enabled` | Follow journal (default **on** if omitted) |
| `journal.args` | Full `journalctl` argument list |
| `authlog.enabled` | Tail a file (default **off**) |
| `authlog.path` | e.g. `/var/log/auth.log` or `/var/log/secure` |
| `notify.on_success` | Notify on successful `Accepted` lines (default **on** if omitted) |
| `notify.dedupe_seconds` | Suppress duplicate user+source+method within this window |

## Command line

```text
sshnotify -config /etc/sshnotify/config.yaml
sshnotify -version
```

## Contributing & forks

Releases, tags, and developer notes: [CONTRIBUTING.md](CONTRIBUTING.md).

## License

MIT — see [LICENSE](LICENSE).
