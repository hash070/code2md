﻿# code2md

将项目代码转换为单个Markdown文件的命令行工具，便于AI阅读整个项目。

## 功能特性

- 递归扫描项目目录，生成树形结构
- 自动识别并排除二进制文件
- 支持.gitignore和自定义排除规则
- 可配置的文件大小限制
- 自动检测文件类型，提供语法高亮
- 跨平台支持（Windows/macOS/Linux）

## 安装

### 从源码编译

```bash
# 克隆或下载源码后
cd code2md
go build -o code2md
```

### 编译为不同平台

```bash
# Windows
GOOS=windows GOARCH=amd64 go build -o code2md.exe

# macOS
GOOS=darwin GOARCH=amd64 go build -o code2md-mac

# Linux
GOOS=linux GOARCH=amd64 go build -o code2md-linux
```

## 使用方法

### 基本用法

```bash
# 在当前目录生成project.md
code2md

# 指定输出文件
code2md -o output.md

# 指定源目录
code2md -s /path/to/project -o output.md

# 设置文件大小限制
code2md -max-size 2MB -o output.md
```

### 命令行参数

- `-s` : 源目录路径（默认：当前目录）
- `-o` : 输出文件路径（默认：project.md）
- `-max-size` : 最大文件大小限制（默认：1MB）
- `-version` : 显示版本信息

## 排除规则

### 默认排除

工具会自动排除以下内容：
- 版本控制目录：`.git`, `.svn`, `.hg`
- IDE配置：`.idea`, `.vscode`
- 依赖目录：`node_modules`, `vendor`
- 编译产物：`dist`, `build`, `target`
- 临时文件：`*.swp`, `*~`, `.DS_Store`
- 二进制文件：自动检测

### 自定义排除规则

1. 使用`.gitignore`文件（自动读取）
2. 创建`.code2mdignore`文件，格式同`.gitignore`

## 输出格式示例

生成的Markdown文件格式如下：

```markdown
# Project Structure

​```
.
├── src/
│   ├── main.go
│   └── utils/
│       └── helper.go
├── go.mod
└── README.md
​```

# Files

## src/main.go
​```go
package main
// 文件内容...
​```

## [Binary] assets/logo.png
*Binary file (245KB)*
```

## 开发

### 项目结构

```
code2md/
├── src/
│   └── main.go      # 主程序
├── go.mod           # Go模块文件
├── README.md        # 本文档
└── .gitignore       # Git忽略规则
```

### 构建脚本

创建`build.sh`（Unix）或`build.bat`（Windows）来自动化构建：

```bash
#!/bin/bash
# build.sh
mkdir -p dist
GOOS=windows GOARCH=amd64 go build -o dist/code2md-windows.exe
GOOS=darwin GOARCH=amd64 go build -o dist/code2md-mac
GOOS=linux GOARCH=amd64 go build -o dist/code2md-linux
```

## 许可证

MIT License
