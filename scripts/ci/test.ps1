Set-StrictMode -Version Latest
$ErrorActionPreference = 'Stop'

Push-Location (Split-Path $PSScriptRoot -Parent | Split-Path -Parent)
try {
  go test ./...
} finally {
  Pop-Location
}
