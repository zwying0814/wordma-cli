package utils

import (
	"fmt"
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