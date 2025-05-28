package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

const (
	defaultMaxSize = 1024 * 1024 // 1MB
	version        = "1.0.0"
)

var (
	// 内置的排除模式
	defaultIgnorePatterns = []string{
		".git",
		".svn",
		".hg",
		".DS_Store",
		"Thumbs.db",
		"*.swp",
		"*.swo",
		"*~",
		".idea",
		".vscode",
		"*.pyc",
		"__pycache__",
		"node_modules",
		".env",
		".env.local",
		"*.log",
		"dist",
		"build",
		"target",
		"*.exe",
		"*.dll",
		"*.so",
		"*.dylib",
	}

	// 二进制文件扩展名
	binaryExtensions = map[string]bool{
		".jpg": true, ".jpeg": true, ".png": true, ".gif": true, ".bmp": true,
		".ico": true, ".tiff": true, ".webp": true, ".svg": true,
		".mp3": true, ".mp4": true, ".avi": true, ".mov": true, ".wmv": true,
		".flv": true, ".webm": true, ".wav": true, ".flac": true,
		".zip": true, ".rar": true, ".7z": true, ".tar": true, ".gz": true,
		".pdf": true, ".doc": true, ".docx": true, ".xls": true, ".xlsx": true,
		".ppt": true, ".pptx": true,
		".exe": true, ".dll": true, ".so": true, ".dylib": true,
		".ttf": true, ".otf": true, ".woff": true, ".woff2": true,
		".db": true, ".sqlite": true,
		".jar": true, ".class": true,
		".o": true, ".a": true, ".lib": true,
	}
)

type Config struct {
	sourceDir   string
	outputFile  string
	maxFileSize int64
	ignoreRules []string
}

func main() {
	config := parseFlags()

	// 加载忽略规则
	ignorePatterns := loadIgnorePatterns(config.sourceDir)

	// 收集所有文件
	files, err := collectFiles(config.sourceDir, ignorePatterns)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error collecting files: %v\n", err)
		os.Exit(1)
	}

	// 生成Markdown
	if err := generateMarkdown(config, files); err != nil {
		fmt.Fprintf(os.Stderr, "Error generating markdown: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Successfully generated %s\n", config.outputFile)
}

func parseFlags() Config {
	var config Config

	flag.StringVar(&config.sourceDir, "s", ".", "Source directory (default: current directory)")
	flag.StringVar(&config.outputFile, "o", "project.md", "Output markdown file")

	var maxSizeStr string
	flag.StringVar(&maxSizeStr, "max-size", "1MB", "Maximum file size to include content")

	var showVersion bool
	flag.BoolVar(&showVersion, "version", false, "Show version")

	flag.Parse()

	if showVersion {
		fmt.Printf("code2md version %s\n", version)
		os.Exit(0)
	}

	// 解析文件大小
	config.maxFileSize = parseSize(maxSizeStr)

	// 转换为绝对路径
	absPath, err := filepath.Abs(config.sourceDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error resolving path: %v\n", err)
		os.Exit(1)
	}
	config.sourceDir = absPath

	return config
}

func parseSize(sizeStr string) int64 {
	sizeStr = strings.ToUpper(sizeStr)
	multiplier := int64(1)

	if strings.HasSuffix(sizeStr, "KB") {
		multiplier = 1024
		sizeStr = strings.TrimSuffix(sizeStr, "KB")
	} else if strings.HasSuffix(sizeStr, "MB") {
		multiplier = 1024 * 1024
		sizeStr = strings.TrimSuffix(sizeStr, "MB")
	} else if strings.HasSuffix(sizeStr, "GB") {
		multiplier = 1024 * 1024 * 1024
		sizeStr = strings.TrimSuffix(sizeStr, "GB")
	}

	var size int64
	fmt.Sscanf(sizeStr, "%d", &size)
	return size * multiplier
}

func loadIgnorePatterns(dir string) []string {
	patterns := make([]string, len(defaultIgnorePatterns))
	copy(patterns, defaultIgnorePatterns)

	// 读取.gitignore
	gitignorePath := filepath.Join(dir, ".gitignore")
	if patterns2 := readIgnoreFile(gitignorePath); patterns2 != nil {
		patterns = append(patterns, patterns2...)
	}

	// 读取.code2mdignore
	code2mdignorePath := filepath.Join(dir, ".code2mdignore")
	if patterns2 := readIgnoreFile(code2mdignorePath); patterns2 != nil {
		patterns = append(patterns, patterns2...)
	}

	return patterns
}

func readIgnoreFile(path string) []string {
	file, err := os.Open(path)
	if err != nil {
		return nil
	}
	defer file.Close()

	var patterns []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" && !strings.HasPrefix(line, "#") {
			patterns = append(patterns, line)
		}
	}
	return patterns
}

func shouldIgnore(path string, patterns []string) bool {
	// 检查每个忽略模式
	for _, pattern := range patterns {
		// 简单的模式匹配（可以后续改进为更复杂的glob匹配）
		if matched, _ := filepath.Match(pattern, filepath.Base(path)); matched {
			return true
		}
		// 检查路径中是否包含该模式
		if strings.Contains(path, pattern) {
			return true
		}
	}
	return false
}

func collectFiles(rootDir string, ignorePatterns []string) ([]string, error) {
	var files []string

	err := filepath.Walk(rootDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// 获取相对路径
		relPath, err := filepath.Rel(rootDir, path)
		if err != nil {
			return err
		}

		// 统一使用正斜杠
		relPath = filepath.ToSlash(relPath)

		// 检查是否应该忽略
		if shouldIgnore(relPath, ignorePatterns) {
			if info.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}

		// 只收集文件，不收集目录
		if !info.IsDir() {
			files = append(files, relPath)
		}

		return nil
	})

	return files, err
}

func generateMarkdown(config Config, files []string) error {
	output, err := os.Create(config.outputFile)
	if err != nil {
		return err
	}
	defer output.Close()

	writer := bufio.NewWriter(output)
	defer writer.Flush()

	// 写入标题
	fmt.Fprintf(writer, "# Project Structure\n\n")

	// 生成目录树
	fmt.Fprintf(writer, "```\n")
	generateTree(writer, config.sourceDir, files)
	fmt.Fprintf(writer, "```\n\n")

	// 写入文件内容
	fmt.Fprintf(writer, "# Files\n\n")

	for _, file := range files {
		fullPath := filepath.Join(config.sourceDir, file)
		if err := writeFileContent(writer, file, fullPath, config.maxFileSize); err != nil {
			fmt.Fprintf(os.Stderr, "Warning: Error processing %s: %v\n", file, err)
		}
	}

	return nil
}

func generateTree(w io.Writer, rootDir string, files []string) {
	// 构建目录结构
	tree := make(map[string][]string)

	for _, file := range files {
		// 使用正斜杠分割路径
		dir := filepath.Dir(file)
		dir = filepath.ToSlash(dir)
		if dir == "." {
			dir = ""
		}
		tree[dir] = append(tree[dir], filepath.Base(file))
	}

	// 打印树形结构
	fmt.Fprintf(w, ".\n")
	printTreeRecursive(w, "", tree, "", "")
}

func printTreeRecursive(w io.Writer, currentPath string, tree map[string][]string, prefix string, childPrefix string) {
	// 获取当前目录下的所有项
	items := make(map[string]bool) // 用于去重

	// 收集子目录
	for dir := range tree {
		if dir == currentPath {
			continue
		}
		if strings.HasPrefix(dir, currentPath) {
			remaining := strings.TrimPrefix(dir, currentPath)
			if currentPath != "" {
				remaining = strings.TrimPrefix(remaining, "/")
			}
			if remaining != "" {
				parts := strings.Split(remaining, "/")
				if len(parts) > 0 && parts[0] != "" {
					items[parts[0]+"/"] = true
				}
			}
		}
	}

	// 收集文件
	if files, ok := tree[currentPath]; ok {
		for _, file := range files {
			items[file] = true
		}
	}

	// 转换为排序的列表
	var itemList []string
	for item := range items {
		itemList = append(itemList, item)
	}

	// 打印项目
	for i, item := range itemList {
		isLast := i == len(itemList)-1

		if isLast {
			fmt.Fprintf(w, "%s└── %s\n", prefix, item)
		} else {
			fmt.Fprintf(w, "%s├── %s\n", prefix, item)
		}

		// 如果是目录，递归打印
		if strings.HasSuffix(item, "/") {
			newPath := currentPath
			if newPath != "" {
				newPath += "/"
			}
			newPath += strings.TrimSuffix(item, "/")

			newPrefix := childPrefix
			if isLast {
				newPrefix += "    "
			} else {
				newPrefix += "│   "
			}
			printTreeRecursive(w, newPath, tree, newPrefix, newPrefix)
		}
	}
}

func writeFileContent(w io.Writer, relPath, fullPath string, maxSize int64) error {
	info, err := os.Stat(fullPath)
	if err != nil {
		return err
	}

	// 统一使用正斜杠作为路径分隔符
	displayPath := filepath.ToSlash(relPath)

	// 检查是否是二进制文件
	ext := strings.ToLower(filepath.Ext(fullPath))
	isBinary := binaryExtensions[ext]

	// 如果不确定，检查文件内容
	if !isBinary && info.Size() > 0 {
		isBinary = isBinaryFile(fullPath)
	}

	fmt.Fprintf(w, "## %s\n", displayPath)

	// 二进制文件只显示信息
	if isBinary {
		fmt.Fprintf(w, "*Binary file (%s)*\n\n", formatSize(info.Size()))
		return nil
	}

	// 大文件只显示信息
	if info.Size() > maxSize {
		fmt.Fprintf(w, "*File too large (%s) - content omitted*\n\n", formatSize(info.Size()))
		return nil
	}

	// 读取并写入文件内容
	content, err := os.ReadFile(fullPath)
	if err != nil {
		return err
	}

	// 获取语言类型（用于语法高亮）
	lang := getLanguageFromExt(ext)

	fmt.Fprintf(w, "```%s\n", lang)
	w.Write(content)
	if !bytes.HasSuffix(content, []byte("\n")) {
		fmt.Fprintf(w, "\n")
	}
	fmt.Fprintf(w, "```\n\n")

	return nil
}

func isBinaryFile(path string) bool {
	file, err := os.Open(path)
	if err != nil {
		return false // 改为false，打开失败不代表是二进制文件
	}
	defer file.Close()

	// 读取前8192字节来判断
	buf := make([]byte, 8192)
	n, err := file.Read(buf)
	if err != nil && err != io.EOF {
		return false
	}

	if n == 0 {
		return false // 空文件当作文本文件处理
	}

	buf = buf[:n]

	// 检查是否包含NULL字符（\0）
	if bytes.Contains(buf, []byte{0}) {
		return true
	}

	// 计算非ASCII字符的比例
	nonASCII := 0
	for _, b := range buf {
		if b > 127 || (b < 32 && b != '\n' && b != '\r' && b != '\t') {
			nonASCII++
		}
	}

	// 如果非ASCII字符超过30%，可能是二进制文件
	if float64(nonASCII)/float64(n) > 0.3 {
		return true
	}

	// 对于常见的文本文件扩展名，直接返回false
	ext := strings.ToLower(filepath.Ext(path))
	textExtensions := map[string]bool{
		".txt": true, ".md": true, ".go": true, ".py": true,
		".js": true, ".ts": true, ".java": true, ".c": true,
		".cpp": true, ".h": true, ".hpp": true, ".cs": true,
		".php": true, ".rb": true, ".rs": true, ".swift": true,
		".kt": true, ".scala": true, ".r": true, ".R": true,
		".sh": true, ".bash": true, ".zsh": true, ".fish": true,
		".ps1": true, ".bat": true, ".cmd": true,
		".sql": true, ".html": true, ".htm": true, ".xml": true,
		".css": true, ".scss": true, ".sass": true, ".less": true,
		".json": true, ".yaml": true, ".yml": true, ".toml": true,
		".ini": true, ".cfg": true, ".conf": true, ".properties": true,
		".vue": true, ".jsx": true, ".tsx": true,
		".tex": true, ".rst": true, ".asciidoc": true,
		".gitignore": true, ".dockerignore": true,
		".editorconfig": true, ".eslintrc": true,
	}

	if textExtensions[ext] {
		return false
	}

	return false // 默认当作文本文件
}

func getLanguageFromExt(ext string) string {
	langMap := map[string]string{
		".go":    "go",
		".py":    "python",
		".js":    "javascript",
		".ts":    "typescript",
		".jsx":   "jsx",
		".tsx":   "tsx",
		".java":  "java",
		".c":     "c",
		".cpp":   "cpp",
		".cc":    "cpp",
		".cxx":   "cpp",
		".h":     "c",
		".hpp":   "cpp",
		".cs":    "csharp",
		".php":   "php",
		".rb":    "ruby",
		".swift": "swift",
		".kt":    "kotlin",
		".rs":    "rust",
		".r":     "r",
		".R":     "r",
		".scala": "scala",
		".sh":    "bash",
		".bash":  "bash",
		".zsh":   "zsh",
		".fish":  "fish",
		".ps1":   "powershell",
		".sql":   "sql",
		".html":  "html",
		".htm":   "html",
		".xml":   "xml",
		".css":   "css",
		".scss":  "scss",
		".sass":  "sass",
		".less":  "less",
		".json":  "json",
		".yaml":  "yaml",
		".yml":   "yaml",
		".toml":  "toml",
		".ini":   "ini",
		".cfg":   "ini",
		".conf":  "conf",
		".md":    "markdown",
		".rst":   "rst",
		".tex":   "latex",
	}

	if lang, ok := langMap[ext]; ok {
		return lang
	}
	return ""
}

func formatSize(size int64) string {
	const unit = 1024
	if size < unit {
		return fmt.Sprintf("%d B", size)
	}
	div, exp := int64(unit), 0
	for n := size / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(size)/float64(div), "KMGTPE"[exp])
}
