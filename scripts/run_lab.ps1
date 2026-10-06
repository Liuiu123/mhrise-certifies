$ErrorActionPreference = "Stop"
$Root = Split-Path -Parent $PSScriptRoot
$Server = Join-Path $Root "server.exe"
if (-not (Test-Path $Server)) { throw "server.exe não encontrado: $Server" }
$env:NEXTENDO_SERVER_IP = "127.0.0.1"
$env:NEXTENDO_NAT_IP = "127.0.0.1"
& (Join-Path $PSScriptRoot "prepare_lab_ca.ps1")
Write-Host "[LAB] Iniciando servidor..."
Start-Process -FilePath $Server -WorkingDirectory $Root
Start-Sleep -Seconds 2
$Ryujinx = Join-Path $Root "Ryujinx-Nextendo-1.7.9\Ryujinx-Nextendo-1.7.9\src\Ryujinx\bin\Debug\net10.0\Ryujinx.exe"
if (-not (Test-Path $Ryujinx)) { throw "Ryujinx.exe não encontrado: $Ryujinx" }
Write-Host "[LAB] Iniciando Ryujinx-Nextendo..."
Start-Process -FilePath $Ryujinx -WorkingDirectory (Split-Path $Ryujinx)
Write-Host "[LAB] Pronto. Abra o MHRise e reproduza o teste NPLN."
