package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"wordma-cli/utils"
)

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update components of the project",
	Long:  "Update various components like themes in the wordma project",
}

var updateThemeCmd = &cobra.Command{
	Use:   "theme <name>",
	Short: "Update a theme to the latest version",
	Long:  "Pull the latest code for the specified theme from its git repository",
	Args:  cobra.ExactArgs(1),
	Run:   runUpdateTheme,
}

func init() {
	updateCmd.AddCommand(updateThemeCmd)
}

func runUpdateTheme(cmd *cobra.Command, args []string) {
	themeName := args[0]

	// 检查 git 是否安装
	if !utils.CheckCommand("git") {
		utils.PrintError("Git is required for updating themes")
		fmt.Printf("  %s\n", utils.GetInstallInstructions("git"))
		os.Exit(1)
	}

	// 获取项目根目录
	projectRoot, err := utils.GetProjectRoot()
	if err != nil {
		utils.PrintError(fmt.Sprintf("Failed to find project root: %v", err))
		os.Exit(1)
	}

	// 检查主题目录是否存在
	themePath := filepath.Join(projectRoot, "themes", themeName)
	if !utils.FileExists(themePath) {
		utils.PrintError(fmt.Sprintf("Theme '%s' not found in themes directory", themeName))
		utils.PrintInfo("Available themes:")
		listAvailableThemes(projectRoot)
		os.Exit(1)
	}

	// 检查主题目录是否是git仓库
	gitPath := filepath.Join(themePath, ".git")
	if !utils.FileExists(gitPath) {
		utils.PrintError(fmt.Sprintf("Theme '%s' is not a git repository", themeName))
		utils.PrintInfo("This theme cannot be updated automatically")
		utils.PrintInfo("You may need to manually update it or re-add it using 'wordma add theme <git-url>'")
		os.Exit(1)
	}

	utils.PrintInfo(fmt.Sprintf("Updating theme '%s'...", themeName))

	// 检查是否有未提交的更改
	err = utils.RunCommandInDir(themePath, "git", "status", "--porcelain")
	if err != nil {
		utils.PrintError(fmt.Sprintf("Failed to check git status: %v", err))
		os.Exit(1)
	}

	// 获取当前分支
	utils.PrintInfo("Fetching latest changes...")
	err = utils.RunCommandInDir(themePath, "git", "fetch", "origin")
	if err != nil {
		utils.PrintError(fmt.Sprintf("Failed to fetch from remote: %v", err))
		os.Exit(1)
	}

	// 获取当前分支名
	currentBranch, err := getCurrentBranch(themePath)
	if err != nil {
		utils.PrintError(fmt.Sprintf("Failed to get current branch: %v", err))
		os.Exit(1)
	}

	// 检查是否有本地更改
	hasLocalChanges, err := hasUncommittedChanges(themePath)
	if err != nil {
		utils.PrintError(fmt.Sprintf("Failed to check for local changes: %v", err))
		os.Exit(1)
	}

	if hasLocalChanges {
		utils.PrintWarning("Theme has uncommitted local changes")
		utils.PrintInfo("Stashing local changes before update...")
		err = utils.RunCommandInDir(themePath, "git", "stash", "push", "-m", "wordma-cli auto stash before update")
		if err != nil {
			utils.PrintError(fmt.Sprintf("Failed to stash changes: %v", err))
			os.Exit(1)
		}
		utils.PrintInfo("Local changes stashed successfully")
	}

	// 拉取最新代码
	utils.PrintInfo(fmt.Sprintf("Pulling latest changes from %s...", currentBranch))
	err = utils.RunCommandInDir(themePath, "git", "pull", "origin", currentBranch)
	if err != nil {
		utils.PrintError(fmt.Sprintf("Failed to pull latest changes: %v", err))
		
		// 如果拉取失败且之前有stash，尝试恢复
		if hasLocalChanges {
			utils.PrintInfo("Attempting to restore stashed changes...")
			restoreErr := utils.RunCommandInDir(themePath, "git", "stash", "pop")
			if restoreErr != nil {
				utils.PrintWarning("Failed to restore stashed changes. You may need to manually run 'git stash pop' in the theme directory")
			}
		}
		os.Exit(1)
	}

	// 如果之前有stash，询问是否恢复
	if hasLocalChanges {
		utils.PrintInfo("Update completed successfully")
		utils.PrintWarning("Your local changes were stashed")
		utils.PrintInfo("To restore your local changes, run:")
		fmt.Printf("  cd %s\n", themePath)
		fmt.Printf("  git stash pop\n")
	} else {
		utils.PrintSuccess(fmt.Sprintf("Theme '%s' updated successfully!", themeName))
	}

	utils.PrintInfo("Next steps:")
	fmt.Printf("  1. wordma install (to update dependencies if needed)\n")
	fmt.Printf("  2. wordma dev %s (to test the updated theme)\n", themeName)
}

// getCurrentBranch 获取当前git分支名
func getCurrentBranch(repoPath string) (string, error) {
	cmd := utils.NewCommand("git", "rev-parse", "--abbrev-ref", "HEAD")
	cmd.Dir = repoPath
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(output)), nil
}

// hasUncommittedChanges 检查是否有未提交的更改
func hasUncommittedChanges(repoPath string) (bool, error) {
	cmd := utils.NewCommand("git", "status", "--porcelain")
	cmd.Dir = repoPath
	output, err := cmd.Output()
	if err != nil {
		return false, err
	}
	return len(output) > 0, nil
}

// listAvailableThemes 列出可用的主题
func listAvailableThemes(projectRoot string) {
	themesDir := filepath.Join(projectRoot, "themes")
	if !utils.FileExists(themesDir) {
		fmt.Println("  No themes directory found")
		return
	}

	entries, err := os.ReadDir(themesDir)
	if err != nil {
		fmt.Printf("  Failed to read themes directory: %v\n", err)
		return
	}

	if len(entries) == 0 {
		fmt.Println("  No themes found")
		return
	}

	for _, entry := range entries {
		if entry.IsDir() {
			themePath := filepath.Join(themesDir, entry.Name())
			gitPath := filepath.Join(themePath, ".git")
			if utils.FileExists(gitPath) {
				fmt.Printf("  - %s (git repository)\n", entry.Name())
			} else {
				fmt.Printf("  - %s (not a git repository)\n", entry.Name())
			}
		}
	}
}