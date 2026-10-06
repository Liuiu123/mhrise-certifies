@echo off
setlocal
echo [+] Installing protoc-gen-go...
go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.36.9
if errorlevel 1 exit /b 1
echo [+] Installing protoc-gen-go-grpc...
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.5.1
if errorlevel 1 exit /b 1
echo [+] Installing Buf...
go install github.com/bufbuild/buf/cmd/buf@v1.56.0
if errorlevel 1 exit /b 1
echo.
echo [+] Done. If buf is not recognized, add %%USERPROFILE%%\goin to PATH.
pause
