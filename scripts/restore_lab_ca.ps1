$ErrorActionPreference = 'Stop'
$SslDir = Join-Path $env:APPDATA 'Ryujinx\system\ssl'
$Dest = Join-Path $SslDir '1054.der'
$Backup = Join-Path $SslDir '1054.der.nextendo-backup'
if (Test-Path $Backup) {
    Move-Item $Backup $Dest -Force
    Write-Host "[OK] Restored original 1054.der"
} elseif (Test-Path $Dest) {
    Remove-Item $Dest -Force
    Write-Host "[OK] Removed lab 1054.der override"
} else {
    Write-Host '[INFO] No lab CA override found.'
}
