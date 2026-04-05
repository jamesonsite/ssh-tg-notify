# Contributing

Issues and pull requests are welcome.

## Development

- Go **1.22+**
- `go test ./...`
- `CGO_ENABLED=0 go build -trimpath -ldflags "-s -w" -o sshnotify ./cmd/sshnotify`

## Releasing

1. Tag a version: `git tag v0.1.0` then `git push origin v0.1.0`.
2. The [release workflow](.github/workflows/release.yml) builds Linux **amd64** and **arm64** tarballs (binary, `config.example.yaml`, `scripts/setup-server.sh`, reference unit) and attaches them to the GitHub Release.

End users extract a tarball and run `sudo ./scripts/setup-server.sh` — no separate “copy template to `/etc`” step.

## Publishing from Windows (maintainers)

If the GitHub repository already exists, a normal `git push` is enough.

To automate “create repo + push” once authenticated with the GitHub CLI:

```powershell
powershell -ExecutionPolicy Bypass -File .\scripts\publish-github.ps1
```

Or set `GH_TOKEN` (classic PAT with `repo` scope) for non-interactive use.

## Forking

After forking, point the Go module at your copy:

```bash
go mod edit -module github.com/<you>/ssh-tg-notify
```

Replace import paths from `github.com/jamesonsite/ssh-tg-notify/...` to match.
