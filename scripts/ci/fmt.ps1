Set-StrictMode -Version Latest
$ErrorActionPreference = 'Stop'

Push-Location (Split-Path $PSScriptRoot -Parent | Split-Path -Parent)
try {
  gofmt -w .
} finally {
  Pop-Location
}
