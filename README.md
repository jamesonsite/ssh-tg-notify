# ssh-tg-notify

Notifies Telegram when someone **successfully logs in over SSH** on a Linux host (systemd, `journalctl` or auth log).

**Needs:** Linux + systemd, outbound HTTPS, typically **root** (reads journal / auth logs).

---

### Telegram

1. [@BotFather](https://t.me/BotFather) → create bot → copy **token** into `telegram.bot_token`.
2. **Chat id** → put in `telegram.chat_id` as a quoted string, e.g. `"123456789"`.
   - DM: message [@userinfobot](https://t.me/userinfobot), use the **Id** it shows.
   - Or message your bot, then open `https://api.telegram.org/bot<TOKEN>/getUpdates` and read `"chat":{"id":…}` in the JSON (groups are often negative ids like `-100…`).
3. With ssh-tg-notify running, anyone who taps **Start** in Telegram with your bot gets setup instructions and an **Open BotFather** button (on by default). Set `telegram.welcome_on_start: false` to turn that off. **Same bot token on many servers:** only one host should leave this on—otherwise several daemons would share `getUpdates` and it is unreliable.

---

### Install

**From git** (needs Go to build on the server):

```bash
sudo apt update && sudo apt install -y git golang-go
git clone https://github.com/jamesonsite/ssh-tg-notify.git && cd ssh-tg-notify
sudo ./scripts/setup-server.sh
nano config.yaml   # token + chat_id
sudo systemctl enable --now sshnotify
```

**From a release tarball** (no Go): `latest` always points at the newest release asset (no version in the URL).

```bash
mkdir -p /opt/ssh-tg-notify && cd /opt/ssh-tg-notify
curl -fsSL -o /tmp/t.tgz 'https://github.com/jamesonsite/ssh-tg-notify/releases/latest/download/sshnotify_linux_amd64.tar.gz'
tar -xzf /tmp/t.tgz
sudo ./scripts/setup-server.sh && nano config.yaml
sudo systemctl enable --now sshnotify
```

ARM: `…/sshnotify_linux_arm64.tar.gz`. If Releases is empty, use **From git** until a maintainer pushes a version tag.

Config changes: `sudo systemctl restart sshnotify`.

---

### `config.yaml`

| Key | Notes |
| --- | --- |
| `telegram.bot_token` | From BotFather |
| `telegram.chat_id` | String, e.g. `"123456789"` |
| `telegram.welcome_on_start` | `true` (default): answer `/start` with help + BotFather button; `false` to skip `getUpdates` polling |
| `server.label` | Optional; empty = hostname |
| `journal.*` | Default: follow `sshd` in journal |
| `authlog.*` | Set `enabled: true` and `path` for file-based logs (e.g. `/var/log/secure`) |
| `notify.*` | `on_success`, `dedupe_seconds` |

`setup-server.sh` creates `config.yaml` from `config.example.yaml` (gitignored).

---

### CLI

`sshnotify -config /path/to/config.yaml` · `sshnotify -version`

---

[CONTRIBUTING.md](CONTRIBUTING.md) · MIT [LICENSE](LICENSE)
