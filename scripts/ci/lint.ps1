Set-StrictMode -Version Latest
$ErrorActionPreference = 'Stop'

Push-Location (Split-Path $PSScriptRoot -Parent | Split-Path -Parent)
try {
  if (Get-Command golangci-lint -ErrorAction SilentlyContinue) {
    golangci-lint run ./...
    exit 0
  }

  Write-Host 'golangci-lint not installed; skipping lint.'
  exit 0
} finally {
  Pop-Location
}
