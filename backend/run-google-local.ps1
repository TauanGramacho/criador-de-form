param(
  [Parameter(Mandatory = $true)]
  [string]$ClientSecretJson
)

$json = Get-Content -LiteralPath $ClientSecretJson -Raw | ConvertFrom-Json

if (-not $json.web.client_id -or -not $json.web.client_secret) {
  throw "Invalid Google OAuth client JSON. Expected web.client_id and web.client_secret."
}

$env:GOOGLE_CLIENT_ID = $json.web.client_id
$env:GOOGLE_CLIENT_SECRET = $json.web.client_secret
$env:GOOGLE_REDIRECT_URL = "http://localhost:8080/api/auth/google/callback"

go run ./cmd/server run
