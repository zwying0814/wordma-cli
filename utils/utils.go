package utils

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/fatih/color"
)

// CheckCommand 检查命令是否存在
func CheckCommand(command string) bool {
	_, err := exec.LookPath(command)
	return err == nil
}

// GetCommandVersion 获取命令版本
func GetCommandVersion(command string, versionFlag string) (string, error) {
	cmd := exec.Command(command, versionFlag)
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(output)), nil
}

// PrintSuccess 打印成功信息
func PrintSuccess(message string) {
	green := color.New(color.FgGreen, color.Bold)
	green.Printf("✓ %s\n", message)
}

// PrintError 打印错误信息
func PrintError(message string) {
	red := color.New(color.FgRed, color.Bold)
	red.Printf("✗ %s\n", message)
}

// PrintWarning 打印警告信息
func PrintWarning(message string) {
	yellow := color.New(color.FgYellow, color.Bold)
	yellow.Printf("⚠ %s\n", message)
}

// PrintInfo 打印信息
func PrintInfo(message string) {
	blue := color.New(color.FgBlue, color.Bold)
	blue.Printf("ℹ %s\n", message)
}

// ColorText 返回带颜色的文本
func ColorText(text, colorName string) string {
	var c *color.Color
	switch colorName {
	case "red":
		c = color.New(color.FgRed)
	case "green":
		c = color.New(color.FgGreen)
	case "yellow":
		c = color.New(color.FgYellow)
	case "blue":
		c = color.New(color.FgBlue)
	case "magenta":
		c = color.New(color.FgMagenta)
	case "cyan":
		c = color.New(color.FgCyan)
	case "white":
		c = color.New(color.FgWhite)
	default:
		return text
	}
	return c.Sprint(text)
}

// RunCommand 执行命令
func RunCommand(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// RunCommandInDir 在指定目录执行命令
func RunCommandInDir(dir, name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Dir = dir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// NewCommand 创建一个新的命令，用于获取输出
func NewCommand(name string, args ...string) *exec.Cmd {
	return exec.Command(name, args...)
}

// FileExists 检查文件是否存在
func FileExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

// IsDirEmpty 检查目录是否为空（忽略指定的文件）
func IsDirEmpty(dir string, ignoreFiles ...string) (bool, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return false, err
	}
	
	// 创建忽略文件的映射
	ignoreMap := make(map[string]bool)
	for _, file := range ignoreFiles {
		ignoreMap[file] = true
	}
	
	// 检查是否有非忽略的文件
	for _, entry := range entries {
		if !ignoreMap[entry.Name()] {
			return false, nil
		}
	}
	
	return true, nil
}

// CreateDir 创建目录
func CreateDir(path string) error {
	return os.MkdirAll(path, 0755)
}

// GetInstallInstructions 获取安装说明
func GetInstallInstructions(tool string) string {
	switch tool {
	case "nodejs":
		switch runtime.GOOS {
		case "windows":
			return "请访问 https://nodejs.org/ 下载并安装 Node.js"
		case "darwin":
			return "使用 Homebrew: brew install node 或访问 https://nodejs.org/"
		default:
			return "使用包管理器安装: sudo apt install nodejs npm 或访问 https://nodejs.org/"
		}
	case "pnpm":
		return "安装 pnpm: npm install -g pnpm 或访问 https://pnpm.io/"
	case "git":
		switch runtime.GOOS {
		case "windows":
			return "请访问 https://git-scm.com/ 下载并安装 Git"
		case "darwin":
			return "使用 Homebrew: brew install git 或访问 https://git-scm.com/"
		default:
			return "使用包管理器安装: sudo apt install git 或访问 https://git-scm.com/"
		}
	default:
		return fmt.Sprintf("请安装 %s", tool)
	}
}

// GetProjectRoot 获取项目根目录
func GetProjectRoot() (string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	
	// 查找包含 package.json 的目录
	for {
		if FileExists(filepath.Join(wd, "package.json")) {
			return wd, nil
		}
		
		parent := filepath.Dir(wd)
		if parent == wd {
			break
		}
		wd = parent
	}
	
	// 如果没找到，返回当前工作目录
	return os.Getwd()
}

// CopyDirectory 复制整个目录
func CopyDirectory(src, dst string) error {
	// 获取源目录信息
	srcInfo, err := os.Stat(src)
	if err != nil {
		return err
	}

	// 创建目标目录
	err = os.MkdirAll(dst, srcInfo.Mode())
	if err != nil {
		return err
	}

	// 读取源目录内容
	entries, err := os.ReadDir(src)
	if err != nil {
		return err
	}

	// 复制每个条目
	for _, entry := range entries {
		srcPath := filepath.Join(src, entry.Name())
		dstPath := filepath.Join(dst, entry.Name())

		if entry.IsDir() {
			// 递归复制子目录
			err = CopyDirectory(srcPath, dstPath)
			if err != nil {
				return err
			}
		} else {
			// 复制文件
			err = CopyFile(srcPath, dstPath)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// CopyFile 复制单个文件
func CopyFile(src, dst string) error {
	// 打开源文件
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	// 获取源文件信息
	srcInfo, err := srcFile.Stat()
	if err != nil {
		return err
	}

	// 创建目标文件
	dstFile, err := os.OpenFile(dst, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, srcInfo.Mode())
	if err != nil {
		return err
	}
	defer dstFile.Close()

	// 复制文件内容
	_, err = io.Copy(dstFile, srcFile)
	return err
}

// CreateDeployReadme 创建 .deploy 目录的 README.md 文件
func CreateDeployReadme(deployPath, projectName string) error {
	readmeContent := fmt.Sprintf(`# %s - Deploy Directory

This directory contains the built static files for the %s wordma blog.

## About

This directory is automatically managed by the wordma CLI tool:
- Built files from themes are placed here
- Each theme build creates a subdirectory named after the theme
- This directory is initialized as a git repository for deployment purposes

## Usage

- Run ` + "`wordma build <theme-name>`" + ` to build a theme
- The built files will be placed in ` + "`./`" + `<theme-name>/` + "`" + ` subdirectory
- You can deploy these files to any static hosting service

## Deployment

You can deploy the contents of this directory to:
- GitHub Pages
- Netlify
- Vercel
- Any static hosting service

## Note

Do not manually edit files in this directory as they will be overwritten on the next build.
`, projectName, projectName)

	readmePath := filepath.Join(deployPath, "README.md")
	return os.WriteFile(readmePath, []byte(readmeContent), 0644)
}

// CreateDeployGitignore 创建 .deploy 目录的 .gitignore 文件
func CreateDeployGitignore(deployPath string) error {
	gitignoreContent := `# Temporary build directory
.temp/

# Node.js dependencies and cache
node_modules/
npm-debug.log*
yarn-debug.log*
yarn-error.log*
pnpm-debug.log*

# Build cache and temporary files
.cache/
.parcel-cache/
.next/
.nuxt/
dist/
build/
out/

# Environment files
.env
.env.local
.env.development.local
.env.test.local
.env.production.local

# IDE and editor files
.vscode/
.idea/
*.swp
*.swo
*~

# OS generated files
.DS_Store
.DS_Store?
._*
.Spotlight-V100
.Trashes
ehthumbs.db
Thumbs.db

# Logs
*.log
logs/

# Runtime data
pids/
*.pid
*.seed
*.pid.lock

# Coverage directory used by tools like istanbul
coverage/
*.lcov

# Dependency directories
jspm_packages/

# Optional npm cache directory
.npm

# Optional eslint cache
.eslintcache

# Microbundle cache
.rpt2_cache/
.rts2_cache_cjs/
.rts2_cache_es/
.rts2_cache_umd/

# Optional REPL history
.node_repl_history

# Output of 'npm pack'
*.tgz

# Yarn Integrity file
.yarn-integrity

# dotenv environment variables file
.env.test

# Stores VSCode versions used for testing VSCode extensions
.vscode-test

# Temporary folders
tmp/
temp/
`

	gitignorePath := filepath.Join(deployPath, ".gitignore")
	return os.WriteFile(gitignorePath, []byte(gitignoreContent), 0644)
}