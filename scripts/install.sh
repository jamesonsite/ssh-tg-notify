#!/usr/bin/env bash
# Install sshnotify binary from GitHub Releases (Linux amd64/arm64).
set -euo pipefail

REPO="${REPO:-jamesonsite/ssh-tg-notify}"
PREFIX="${PREFIX:-/usr/local}"
VERSION="${VERSION:-}"

tmp="$(mktemp -d)"
trap 'rm -rf "$tmp"' EXIT

arch="$(uname -m)"
case "$arch" in
  x86_64) goarch=amd64 ;;
  aarch64|arm64) goarch=arm64 ;;
  *) echo "unsupported arch: $arch" >&2; exit 1 ;;
esac

if [[ -z "$VERSION" || "$VERSION" == "latest" ]]; then
  VERSION="$(curl -fsSL "https://api.github.com/repos/${REPO}/releases/latest" | sed -n 's/.*"tag_name": *"\([^"]*\)".*/\1/p')"
fi
[[ -n "$VERSION" ]] || { echo "could not resolve version" >&2; exit 1; }

asset="sshnotify_${VERSION#v}_linux_${goarch}.tar.gz"
url="https://github.com/${REPO}/releases/download/${VERSION}/${asset}"

echo "Downloading ${url}"
curl -fsSL "$url" -o "$tmp/$asset"
tar -xzf "$tmp/$asset" -C "$tmp"

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
