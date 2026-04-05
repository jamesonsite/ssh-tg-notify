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

## 1. Create the Telegram bot and get a **chat id**

### Bot token (from BotFather)

1. In Telegram, open [@BotFather](https://t.me/BotFather).
2. Send `/newbot` (or use an existing bot), follow the prompts, then copy the **token** BotFather gives you. It looks like `123456789:AAHxxxxxxxxxxxxxxxxxxxxxxxxxxx` — that string goes in `telegram.bot_token` in your config.

### Chat id — what it is

Telegram needs to know **where to deliver** each alert. That destination is a numeric **chat id**:

- If alerts should go to **you personally** (DM), the chat id is **your Telegram user id** (a positive number).
- If alerts should go to a **group**, the chat id is usually a **negative** number (often starting with `-100`).

Your config field `telegram.chat_id` must be that number in quotes, e.g. `"123456789"` or `"-1001234567890"`.

### Option A — Easiest for DMs only: [@userinfobot](https://t.me/userinfobot)

This bot tells you **your own** user id. For “notify me in private,” that id is exactly what you use as `chat_id`.

1. In Telegram, search for **`@userinfobot`** and open the chat.
2. Tap **Start** (or send `/start`).
3. It replies with your **Id** — a number like `123456789`.
4. Put that in config as `chat_id: "123456789"` (same digits, with quotes).

You should **still** open your **new** bot once and tap **Start** / send any message so Telegram has an active chat with it (some setups expect that).

### Option B — Always works: `getUpdates` (DM or group)

Use this if you want to double-check the id or you are sending to a **group**.

1. Open a chat with **your** bot (the one you created with BotFather) and send any message, e.g. `hi`.
2. In a browser, open this URL (paste your **real** token where shown):

   `https://api.telegram.org/bot<YOUR_BOT_TOKEN>/getUpdates`

   Example: if your token is `123:ABC`, the URL is `https://api.telegram.org/bot123:ABC/getUpdates`

3. You will see JSON. Find the first `"chat"` block and look for **`"id"`** next to it, for example:

   `"chat":{"id":123456789,...}`

4. That number is your **chat id**. For a **group**, repeat after posting a message **in the group** (you may need to add the bot to the group first); the id is often negative, e.g. `-1001234567890`.

If `getUpdates` shows `"ok":true` but `"result":[]` (empty), you did not message your bot yet — go back to step 1.

### Groups (short version)

1. Add your bot to the group (group settings → add members → your bot).
2. Send a normal message in the group (some bots need `/start@YourBotName` once).
3. Call `getUpdates` again and read the **`"id"`** inside `"chat"` for that group message — use that (including the minus sign) as `chat_id`.

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
