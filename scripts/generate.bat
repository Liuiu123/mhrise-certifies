@echo off
setlocal
cd /d "%~dp0.."
where buf >nul 2>nul
if errorlevel 1 (
  echo [ERROR] buf is not in PATH.
  echo Run scripts\install_codegen.bat and open a new CMD.
  exit /b 1
)
echo [+] Generating Go protobuf/gRPC files...
buf generate
if errorlevel 1 exit /b 1
echo [+] Generation complete.
dir /s /b generated\*.go
