$ErrorActionPreference = "Stop"
$Root = Split-Path -Parent $PSScriptRoot
$RyujinxRoot = Join-Path $Root "Ryujinx-Nextendo-1.7.9\Ryujinx-Nextendo-1.7.9\src\Ryujinx\bin\Debug\net10.0"
$Source = Join-Path $Root "certs\lab-root-ca.der"
$DataDir = Join-Path $env:APPDATA "Ryujinx"
$DestDir = Join-Path $DataDir "system\ssl"
$Dest = Join-Path $DestDir "nextendo-lab-root.der"
if (-not (Test-Path $Source)) { throw "CA de laboratório não encontrada: $Source" }
New-Item -ItemType Directory -Force -Path $DestDir | Out-Null
Copy-Item $Source $Dest -Force
Write-Host "[OK] CA do laboratório instalada em $Dest"
Write-Host "[INFO] O Ryujinx-Nextendo foi alterado para carregar esta CA somente quando o arquivo existir."
