#!/usr/bin/env bash
# Install sshnotify from GitHub Releases (Linux amd64/arm64). Uses latest stable asset URL.
set -euo pipefail

REPO="${REPO:-jamesonsite/ssh-tg-notify}"
PREFIX="${PREFIX:-/usr/local}"

tmp="$(mktemp -d)"
trap 'rm -rf "$tmp"' EXIT

arch="$(uname -m)"
case "$arch" in
  x86_64) goarch=amd64 ;;
  aarch64|arm64) goarch=arm64 ;;
  *) echo "unsupported arch: $arch" >&2; exit 1 ;;
esac

url="https://github.com/${REPO}/releases/latest/download/sshnotify_linux_${goarch}.tar.gz"
echo "Downloading ${url}"
curl -fsSL "$url" -o "$tmp/bundle.tgz"
tar -xzf "$tmp/bundle.tgz" -C "$tmp"

install -d "$PREFIX/bin"
install -m 0755 "$tmp/sshnotify" "$PREFIX/bin/sshnotify"
install -d /etc/sshnotify
if [[ ! -f /etc/sshnotify/config.yaml ]]; then
  install -m 0600 "$tmp/config.example.yaml" /etc/sshnotify/config.yaml
  echo "Installed example config at /etc/sshnotify/config.yaml — edit with your bot token and chat id."
fi

if [[ -d /etc/systemd/system ]] && [[ -f "$tmp/sshnotify.service" ]]; then
  install -m 0644 "$tmp/sshnotify.service" /etc/systemd/system/sshnotify.service
  echo "Installed systemd unit. Run: sudo systemctl enable --now sshnotify"
fi

echo "sshnotify installed to $PREFIX/bin/sshnotify"
