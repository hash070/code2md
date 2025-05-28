@echo off

REM 构建脚本 - Windows版本

REM 设置版本号
set VERSION=1.0.0

REM 创建dist目录
echo Creating dist directory...
if not exist dist mkdir dist

REM 设置编译参数
set LDFLAGS=-s -w -X main.version=%VERSION%

echo Building code2md v%VERSION%...

REM Windows AMD64
echo Building Windows AMD64...
set GOOS=windows
set GOARCH=amd64
go build -ldflags "%LDFLAGS%" -o dist\code2md-windows-amd64.exe

REM Windows 386
echo Building Windows 386...
set GOOS=windows
set GOARCH=386
go build -ldflags "%LDFLAGS%" -o dist\code2md-windows-386.exe

REM macOS AMD64
echo Building macOS AMD64...
set GOOS=darwin
set GOARCH=amd64
go build -ldflags "%LDFLAGS%" -o dist\code2md-darwin-amd64

REM macOS ARM64
echo Building macOS ARM64...
set GOOS=darwin
set GOARCH=arm64
go build -ldflags "%LDFLAGS%" -o dist\code2md-darwin-arm64

REM Linux AMD64
echo Building Linux AMD64...
set GOOS=linux
set GOARCH=amd64
go build -ldflags "%LDFLAGS%" -o dist\code2md-linux-amd64

REM Linux ARM64
echo Building Linux ARM64...
set GOOS=linux
set GOARCH=arm64
go build -ldflags "%LDFLAGS%" -o dist\code2md-linux-arm64

echo.
echo Build complete! Binaries are in dist\
echo.
dir dist\

pause