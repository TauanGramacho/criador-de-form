param(
  [string]$AppName = "criador-de-form",
  [string]$ProjectRef = "gptvpmdacbhgnyisvzsl"
)

$dbPassword = Read-Host -AsSecureString "Digite a senha do banco Supabase"
$ptr = [Runtime.InteropServices.Marshal]::SecureStringToBSTR($dbPassword)

try {
  $plain = [Runtime.InteropServices.Marshal]::PtrToStringBSTR($ptr)
  $escaped = [System.Uri]::EscapeDataString($plain)
  $databaseUrl = "postgresql://postgres:$escaped@db.$ProjectRef.supabase.co:5432/postgres?sslmode=require"

  fly secrets set --app $AppName DATABASE_URL="$databaseUrl" --stage
}
finally {
  if ($ptr -ne [IntPtr]::Zero) {
    [Runtime.InteropServices.Marshal]::ZeroFreeBSTR($ptr)
  }
}
