# ssh-tg-notify

Notifies Telegram when someone **successfully logs in over SSH** on a Linux host (systemd, `journalctl` or auth log).

**Needs:** Linux + systemd, outbound HTTPS, typically **root** (reads journal / auth logs).

---

### Telegram

1. [@BotFather](https://t.me/BotFather) → create bot → copy **token** into `telegram.bot_token`.
2. **Chat id** → put in `telegram.chat_id` as a quoted string, e.g. `"123456789"`.
   - DM: message [@userinfobot](https://t.me/userinfobot), use the **Id** it shows.
   - Or message your bot, then open `https://api.telegram.org/bot<TOKEN>/getUpdates` and read `"chat":{"id":…}` in the JSON (groups are often negative ids like `-100…`).

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

**From a release tarball** (no Go): only works after [Releases](https://github.com/jamesonsite/ssh-tg-notify/releases) lists a version with `sshnotify_*_linux_*.tar.gz`. If there is nothing there yet, use **From git** above, or wait until CI finishes after a tag is pushed.

```bash
mkdir -p /opt/ssh-tg-notify && cd /opt/ssh-tg-notify
curl -fsSL -o /tmp/t.tgz 'https://github.com/jamesonsite/ssh-tg-notify/releases/download/v0.1.0/sshnotify_v0.1.0_linux_amd64.tar.gz'
tar -xzf /tmp/t.tgz
sudo ./scripts/setup-server.sh && nano config.yaml
sudo systemctl enable --now sshnotify
```

On ARM use `…_linux_arm64.tar.gz`. For a newer version, change both path segments (`…/download/v0.2.0/sshnotify_v0.2.0_linux_amd64.tar.gz`).

Config changes: `sudo systemctl restart sshnotify`.

---

### `config.yaml`

| Key | Notes |
| --- | --- |
| `telegram.bot_token` | From BotFather |
| `telegram.chat_id` | String, e.g. `"123456789"` |
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
