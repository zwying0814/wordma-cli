package cmd

import (
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"wordma-cli/utils"
)

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add components to the project",
	Long:  "Add various components like themes to the wordma project",
}

var addThemeCmd = &cobra.Command{
	Use:   "theme <git-url>",
	Short: "Add a theme from a git repository",
	Long:  "Clone a theme from the specified git URL and add it to the themes directory",
	Args:  cobra.ExactArgs(1),
	Run:   runAddTheme,
}

func init() {
	addCmd.AddCommand(addThemeCmd)
}

func runAddTheme(cmd *cobra.Command, args []string) {
	gitURL := args[0]

	// 检查 git 是否安装
	if !utils.CheckCommand("git") {
		utils.PrintError("Git is required for adding themes")
		fmt.Printf("  %s\n", utils.GetInstallInstructions("git"))
		os.Exit(1)
	}

	// 验证 URL 格式
	if !isValidGitURL(gitURL) {
		utils.PrintError("Invalid git URL format")
		utils.PrintInfo("Examples of valid URLs:")
		fmt.Println("  - https://github.com/user/theme.git")
		fmt.Println("  - git@github.com:user/theme.git")
		os.Exit(1)
	}

	// 获取项目根目录
	projectRoot, err := utils.GetProjectRoot()
	if err != nil {
		utils.PrintError(fmt.Sprintf("Failed to find project root: %v", err))
		os.Exit(1)
	}

	// 确保 themes 目录存在
	themesDir := filepath.Join(projectRoot, "themes")
	if !utils.FileExists(themesDir) {
		err = utils.CreateDir(themesDir)
		if err != nil {
			utils.PrintError(fmt.Sprintf("Failed to create themes directory: %v", err))
			os.Exit(1)
		}
		utils.PrintInfo("Created themes directory")
	}

	// 从 URL 提取主题名称
	themeName := extractThemeNameFromURL(gitURL)
	themePath := filepath.Join(themesDir, themeName)

	// 检查主题是否已存在
	if utils.FileExists(themePath) {
		utils.PrintError(fmt.Sprintf("Theme '%s' already exists in themes directory", themeName))
		os.Exit(1)
	}

	utils.PrintInfo(fmt.Sprintf("Adding theme '%s' from %s...", themeName, gitURL))

	// 克隆主题仓库
	err = utils.RunCommand("git", "clone", gitURL, themePath)
	if err != nil {
		utils.PrintError(fmt.Sprintf("Failed to clone theme repository: %v", err))
		os.Exit(1)
	}

	// 删除主题的 .git 目录（可选）
	themeGitPath := filepath.Join(themePath, ".git")
	if utils.FileExists(themeGitPath) {
		err = os.RemoveAll(themeGitPath)
		if err != nil {
			utils.PrintWarning(fmt.Sprintf("Failed to remove theme .git directory: %v", err))
		} else {
			utils.PrintInfo("Removed theme .git directory")
		}
	}

	utils.PrintSuccess(fmt.Sprintf("Theme '%s' added successfully!", themeName))
	utils.PrintInfo("Next steps:")
	fmt.Printf("  1. wordma install (if not already done)\n")
	fmt.Printf("  2. wordma dev %s\n", themeName)
}

// isValidGitURL 检查是否为有效的 git URL
func isValidGitURL(gitURL string) bool {
	// 检查 HTTPS URL
	if strings.HasPrefix(gitURL, "https://") {
		_, err := url.Parse(gitURL)
		return err == nil
	}

	// 检查 SSH URL (git@...)
	if strings.HasPrefix(gitURL, "git@") {
		return strings.Contains(gitURL, ":")
	}

	// 检查其他 git:// 协议
	if strings.HasPrefix(gitURL, "git://") {
		_, err := url.Parse(gitURL)
		return err == nil
	}

	return false
}

// extractThemeNameFromURL 从 git URL 提取主题名称
func extractThemeNameFromURL(gitURL string) string {
	// 移除协议部分
	url := gitURL
	if strings.HasPrefix(url, "https://") {
		url = strings.TrimPrefix(url, "https://")
	} else if strings.HasPrefix(url, "git@") {
		url = strings.TrimPrefix(url, "git@")
		url = strings.Replace(url, ":", "/", 1)
	} else if strings.HasPrefix(url, "git://") {
		url = strings.TrimPrefix(url, "git://")
	}

	// 获取路径的最后一部分
	parts := strings.Split(url, "/")
	if len(parts) > 0 {
		name := parts[len(parts)-1]
		// 移除 .git 后缀
		name = strings.TrimSuffix(name, ".git")
		return name
	}

	return "theme"
}
