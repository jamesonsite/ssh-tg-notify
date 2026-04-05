# One-shot: ensure gh auth, create public repo if missing, push main.
# Run from repo root:  powershell -ExecutionPolicy Bypass -File scripts/publish-github.ps1
$ErrorActionPreference = "Stop"
$env:Path = [System.Environment]::GetEnvironmentVariable("Path", "Machine") + ";" +
  [System.Environment]::GetEnvironmentVariable("Path", "User")

$repoRoot = Split-Path $PSScriptRoot -Parent
Set-Location $repoRoot

gh auth status 2>$null | Out-Null
if ($LASTEXITCODE -ne 0) {
  Write-Host "GitHub CLI is not logged in. Complete login in the browser when prompted."
  gh auth login --hostname github.com --git-protocol https --web
}

$ownerRepo = "jamesonsite/ssh-tg-notify"
gh repo view $ownerRepo 2>$null | Out-Null
if ($LASTEXITCODE -ne 0) {
  Write-Host "Creating public repo $ownerRepo ..."
  gh repo create $ownerRepo --public --description "SSH login -> Telegram notifier for Linux (systemd agent)"
}

git remote remove origin 2>$null
git remote add origin "https://github.com/$ownerRepo.git"
git push -u origin main
Write-Host "Done. Remote: https://github.com/$ownerRepo"
