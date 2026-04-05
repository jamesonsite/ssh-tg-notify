#!/usr/bin/env bash
# Run from the cloned repo (use sudo): creates config.yaml, builds or uses ./sshnotify, systemd unit.
set -euo pipefail

REPO_ROOT="$(cd "$(dirname "$0")/.." && pwd)"
cd "$REPO_ROOT"

if [[ "${EUID:-0}" -ne 0 ]]; then
  echo "Run with sudo, e.g.:  sudo ./scripts/setup-server.sh" >&2
  exit 1
fi

if [[ ! -f config.example.yaml ]]; then
  echo "config.example.yaml missing — run this from a clone of the repository." >&2
  exit 1
fi

if [[ ! -f config.yaml ]]; then
  cp -a config.example.yaml config.yaml
  chmod 600 config.yaml
  echo "Created ${REPO_ROOT}/config.yaml (from example). Edit telegram.bot_token and telegram.chat_id next."
else
  echo "Using existing ${REPO_ROOT}/config.yaml"
fi

BIN="/usr/local/bin/sshnotify"
# Prefer a fresh Go build when source exists (avoids re-installing a stale ./sshnotify from an old tarball).
if command -v go >/dev/null 2>&1 && [[ -f "${REPO_ROOT}/cmd/sshnotify/main.go" ]]; then
  echo "Building sshnotify with Go..."
  ( cd "$REPO_ROOT" && CGO_ENABLED=0 go build -trimpath -ldflags "-s -w" -o sshnotify ./cmd/sshnotify )
  install -m 0755 "${REPO_ROOT}/sshnotify" "${BIN}"
elif [[ -x "${REPO_ROOT}/sshnotify" ]]; then
  echo "Installing ${REPO_ROOT}/sshnotify -> ${BIN}"
  echo "(Tarball-only tree: to upgrade, curl the latest release tarball here again and re-run this script, or git clone the repo and use Go.)"
  install -m 0755 "${REPO_ROOT}/sshnotify" "${BIN}"
elif command -v go >/dev/null 2>&1; then
  echo "go is installed but ${REPO_ROOT}/cmd/sshnotify/main.go is missing (not a full git checkout)." >&2
  echo "Download a release tarball with sshnotify inside, or: git clone https://github.com/jamesonsite/ssh-tg-notify.git" >&2
  exit 1
else
  echo "No ./sshnotify binary and no 'go' in PATH." >&2
  echo "Either: apt install golang-go and git clone this repo, then re-run," >&2
  echo "Or: curl the latest sshnotify_linux_*.tar.gz from GitHub Releases, extract here, re-run." >&2
  exit 1
fi

UNIT=/etc/systemd/system/sshnotify.service
cat >"$UNIT" <<EOF
[Unit]
Description=SSH login Telegram notifier
Documentation=https://github.com/jamesonsite/ssh-tg-notify
After=network-online.target
Wants=network-online.target

[Service]
Type=simple
ExecStart=${BIN} -config ${REPO_ROOT}/config.yaml
Restart=on-failure
RestartSec=5

[Install]
WantedBy=multi-user.target
EOF
chmod 644 "$UNIT"
systemctl daemon-reload

echo ""
echo "Done. Next:"
echo "  1. Edit:  ${REPO_ROOT}/config.yaml"
echo "  2. Start:  systemctl enable --now sshnotify"
echo "After config changes:  systemctl restart sshnotify"
