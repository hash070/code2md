#!/bin/bash

# code2md 快速安装脚本

set -e

# 检测操作系统和架构
OS=""
ARCH=""

# 检测操作系统
if [[ "$OSTYPE" == "linux-gnu"* ]]; then
    OS="linux"
elif [[ "$OSTYPE" == "darwin"* ]]; then
    OS="darwin"
else
    echo "Unsupported operating system: $OSTYPE"
    exit 1
fi

# 检测架构
MACHINE_TYPE=$(uname -m)
if [[ "$MACHINE_TYPE" == "x86_64" ]]; then
    ARCH="amd64"
elif [[ "$MACHINE_TYPE" == "aarch64" ]] || [[ "$MACHINE_TYPE" == "arm64" ]]; then
    ARCH="arm64"
elif [[ "$MACHINE_TYPE" == "i386" ]] || [[ "$MACHINE_TYPE" == "i686" ]]; then
    ARCH="386"
else
    echo "Unsupported architecture: $MACHINE_TYPE"
    exit 1
fi

# 构建二进制名称
BINARY_NAME="code2md-${OS}-${ARCH}"

echo "Detected system: $OS $ARCH"
echo "Building $BINARY_NAME..."

# 编译
go build -ldflags "-s -w" -o code2md main.go

# 安装到系统路径
INSTALL_PATH="/usr/local/bin"

# 检查是否有权限
if [[ -w "$INSTALL_PATH" ]]; then
    mv code2md "$INSTALL_PATH/"
    echo "Installed to $INSTALL_PATH/code2md"
else
    echo "No write permission to $INSTALL_PATH"
    echo "Trying with sudo..."
    sudo mv code2md "$INSTALL_PATH/"
    echo "Installed to $INSTALL_PATH/code2md (with sudo)"
fi

# 验证安装
if command -v code2md &> /dev/null; then
    echo ""
    echo "Installation successful!"
    echo "Run 'code2md -version' to verify."
else
    echo ""
    echo "Installation may have failed."
    echo "Please check if $INSTALL_PATH is in your PATH."
fi