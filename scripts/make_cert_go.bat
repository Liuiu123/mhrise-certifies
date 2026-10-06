@echo off
setlocal
cd /d "%~dp0.."
go run server\certgen.go
