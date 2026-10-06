$ErrorActionPreference = 'Stop'
$Root = Split-Path -Parent $PSScriptRoot
$Server = Join-Path $Root 'server.exe'
$Ryujinx = Join-Path $Root 'Ryujinx\Ryujinx.exe'
$LabData = Join-Path $Root 'lab-runtime'
$UserData = Join-Path $env:APPDATA 'Ryujinx'

if (-not (Test-Path $Server)) { throw "server.exe não encontrado: $Server" }
if (-not (Test-Path $Ryujinx)) { throw "Ryujinx.exe não encontrado: $Ryujinx" }
if (-not (Test-Path $UserData)) { throw "Pasta do Ryujinx não encontrada: $UserData" }

function Ensure-Junction([string]$Path, [string]$Target) {
    if (Test-Path $Path) { return }
    New-Item -ItemType Junction -Path $Path -Target $Target | Out-Null
}

New-Item -ItemType Directory -Force -Path $LabData | Out-Null
Ensure-Junction (Join-Path $LabData 'bis') (Join-Path $UserData 'bis')
Ensure-Junction (Join-Path $LabData 'games') (Join-Path $UserData 'games')
Ensure-Junction (Join-Path $LabData 'profiles') (Join-Path $UserData 'profiles')
Ensure-Junction (Join-Path $LabData 'system') (Join-Path $UserData 'system')

$SslDir = Join-Path $UserData 'system\ssl'
New-Item -ItemType Directory -Force -Path $SslDir | Out-Null
$DestCa = Join-Path $SslDir '1054.der'
$BackupCa = Join-Path $SslDir '1054.der.nextendo-backup'
if ((Test-Path $DestCa) -and -not (Test-Path $BackupCa)) { Copy-Item $DestCa $BackupCa -Force }
Copy-Item (Join-Path $Root 'certs\lab-root-ca.der') $DestCa -Force

$env:NEXTENDO_SERVER_IP = '127.0.0.1'
$env:NEXTENDO_NAT_IP = '127.0.0.1'

Write-Host ''
Write-Host '=============================================='
Write-Host ' MHRise NPLN Lab - READY TEST'
Write-Host '=============================================='
Write-Host "[LAB] Server : $Server"
Write-Host "[LAB] Ryujinx: $Ryujinx"
Write-Host "[LAB] Data   : $LabData"
Write-Host '[LAB] NPLN   : 127.0.0.1:443'
Write-Host '[LAB] CA     : system\ssl\1054.der'
Write-Host '=============================================='

Write-Host '[1/2] Starting NPLN server...'
Start-Process -FilePath $Server -WorkingDirectory $Root
Start-Sleep -Seconds 2

Write-Host '[2/2] Starting Ryujinx-Nextendo...'
Start-Process -FilePath $Ryujinx -ArgumentList @('--root-data-dir', $LabData) -WorkingDirectory (Split-Path $Ryujinx)
Write-Host '[OK] Lab started. Reproduce the MHRise NPLN flow.'
