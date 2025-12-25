Set-StrictMode -Version Latest
$ErrorActionPreference = 'Stop'

Push-Location (Split-Path $PSScriptRoot -Parent | Split-Path -Parent)
try {
  if (Get-Command govulncheck -ErrorAction SilentlyContinue) {
    govulncheck ./...
    exit 0
  }

  Write-Host 'govulncheck not installed; skipping vuln scan.'
  exit 0
} finally {
  Pop-Location
}
