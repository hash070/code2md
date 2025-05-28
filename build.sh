#!/bin/bash

# 构建脚本 - 编译多平台版本

# 设置版本号
VERSION="1.0.0"
BUILD_DATE=$(date +%Y%m%d)

# 创建dist目录
echo "Creating dist directory..."
mkdir -p dist

# 设置编译参数
LDFLAGS="-s -w -X main.version=$VERSION"

echo "Building code2md v$VERSION..."

# Windows AMD64
echo "Building Windows AMD64..."
GOOS=windows GOARCH=amd64 go build -ldflags "$LDFLAGS" -o dist/code2md-windows-amd64.exe

# Windows 386
echo "Building Windows 386..."
GOOS=windows GOARCH=386 go build -ldflags "$LDFLAGS" -o dist/code2md-windows-386.exe

# macOS AMD64
echo "Building macOS AMD64..."
GOOS=darwin GOARCH=amd64 go build -ldflags "$LDFLAGS" -o dist/code2md-darwin-amd64

# macOS ARM64 (M1/M2)
echo "Building macOS ARM64..."
GOOS=darwin GOARCH=arm64 go build -ldflags "$LDFLAGS" -o dist/code2md-darwin-arm64

# Linux AMD64
echo "Building Linux AMD64..."
GOOS=linux GOARCH=amd64 go build -ldflags "$LDFLAGS" -o dist/code2md-linux-amd64

# Linux ARM64
echo "Building Linux ARM64..."
GOOS=linux GOARCH=arm64 go build -ldflags "$LDFLAGS" -o dist/code2md-linux-arm64

# Linux 386
echo "Building Linux 386..."
GOOS=linux GOARCH=386 go build -ldflags "$LDFLAGS" -o dist/code2md-linux-386

echo "Build complete! Binaries are in dist/"
echo ""
echo "File sizes:"
ls -lh dist/

# 可选：创建压缩包
if command -v zip &> /dev/null; then
    echo ""
    echo "Creating release archives..."
    cd dist
    for file in *; do
        if [[ $file == *.exe ]]; then
            zip "${file%.exe}-v$VERSION.zip" "$file"
        else
            tar -czf "$file-v$VERSION.tar.gz" "$file"
        fi
    done
    cd ..
    echo "Archives created!"
fi