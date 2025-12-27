Set-StrictMode -Version Latest
$ErrorActionPreference = 'Stop'

Push-Location (Split-Path $PSScriptRoot -Parent | Split-Path -Parent)
try {
  go test -cover ./...
} finally {
  Pop-Location
}
