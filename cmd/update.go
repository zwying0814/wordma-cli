package cmd

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"wordma-cli/utils"
)

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update wordma CLI or project components",
	Long:  "Update wordma CLI to the latest version or update project components like themes",
	Run:   runSelfUpdate,
}

var updateThemesCmd = &cobra.Command{
	Use:   "themes",
	Short: "Update project themes",
	Long:  "Update various themes in the wordma project",
}

var updateThemeCmd = &cobra.Command{
	Use:   "theme <name>",
	Short: "Update a theme to the latest version",
	Long:  "Pull the latest code for the specified theme from its git repository",
	Args:  cobra.ExactArgs(1),
	Run:   runUpdateTheme,
}

func init() {
	updateThemesCmd.AddCommand(updateThemeCmd)
	updateCmd.AddCommand(updateThemesCmd)
}

// runSelfUpdate handles the self-update functionality
func runSelfUpdate(cmd *cobra.Command, args []string) {
	fmt.Printf("%s Checking for updates...\n", utils.ColorText("🔍", "blue"))

	// Get version information
	versionInfo, err := getVersionInfo()
	if err != nil {
		utils.PrintError(fmt.Sprintf("Failed to check for updates: %v", err))
		os.Exit(1)
	}

	// Display current version info
	fmt.Printf("%s Current version: %s\n", 
		utils.ColorText("📦", "blue"), 
		utils.ColorText(versionInfo.Current, "green"))
	fmt.Printf("%s Latest version: %s\n", 
		utils.ColorText("🚀", "blue"), 
		utils.ColorText(versionInfo.Latest, "green"))

	if !versionInfo.NeedsUpdate {
		fmt.Printf("\n%s You are already using the latest version!\n", 
			utils.ColorText("✅", "green"))
		return
	}

	// Check if this is a development version
	if versionInfo.Current == "dev" {
		fmt.Printf("\n%s You are using a development version.\n", 
			utils.ColorText("⚠️", "yellow"))
		fmt.Printf("%s Auto-update is not available for development builds.\n", 
			utils.ColorText("💡", "blue"))
		fmt.Printf("%s Please download the latest release from: %s\n", 
			utils.ColorText("🔗", "blue"),
			utils.ColorText("https://github.com/zwying0814/wordma-cli/releases", "cyan"))
		return
	}

	// Confirm update
	fmt.Printf("\n%s A new version (%s) is available!\n", 
		utils.ColorText("🎉", "yellow"), 
		utils.ColorText(versionInfo.Latest, "yellow"))
	fmt.Printf("%s Do you want to update? (y/N): ", 
		utils.ColorText("❓", "blue"))

	reader := bufio.NewReader(os.Stdin)
	response, _ := reader.ReadString('\n')
	response = strings.TrimSpace(strings.ToLower(response))

	if response != "y" && response != "yes" {
		fmt.Printf("%s Update cancelled.\n", utils.ColorText("❌", "red"))
		return
	}

	// Perform update
	if err := performSelfUpdate(versionInfo.Latest); err != nil {
		utils.PrintError(fmt.Sprintf("Update failed: %v", err))
		os.Exit(1)
	}

	fmt.Printf("\n%s Successfully updated to version %s!\n", 
		utils.ColorText("🎉", "green"), 
		utils.ColorText(versionInfo.Latest, "green"))
	fmt.Printf("%s Please restart wordma to use the new version.\n", 
		utils.ColorText("💡", "blue"))

	// Show changelog if available
	if versionInfo.Changelog != "" {
		fmt.Printf("\n%s Release Notes:\n", utils.ColorText("📝", "blue"))
		fmt.Printf("%s\n", versionInfo.Changelog)
	}
}

// performSelfUpdate downloads and installs the latest version
func performSelfUpdate(version string) error {
	fmt.Printf("%s Downloading version %s...\n", 
		utils.ColorText("⬇️", "blue"), version)

	// Get current executable path
	execPath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("failed to get executable path: %w", err)
	}

	// Determine download URL based on platform
	downloadURL := getDownloadURL(version)
	if downloadURL == "" {
		return fmt.Errorf("unsupported platform: %s/%s", runtime.GOOS, runtime.GOARCH)
	}

	// Download the new version
	tempFile, err := downloadFile(downloadURL)
	if err != nil {
		return fmt.Errorf("failed to download update: %w", err)
	}
	defer os.Remove(tempFile)

	fmt.Printf("%s Installing update...\n", utils.ColorText("🔧", "blue"))

	// Backup current executable
	backupPath := execPath + ".backup"
	if err := copyFileForUpdate(execPath, backupPath); err != nil {
		return fmt.Errorf("failed to backup current executable: %w", err)
	}

	// Replace executable
	if err := replaceExecutable(tempFile, execPath); err != nil {
		// Restore backup on failure
		copyFileForUpdate(backupPath, execPath)
		os.Remove(backupPath)
		return fmt.Errorf("failed to replace executable: %w", err)
	}

	// Clean up backup
	os.Remove(backupPath)
	return nil
}

// getDownloadURL returns the download URL for the current platform
func getDownloadURL(version string) string {
	baseURL := "https://github.com/zwying0814/wordma-cli/releases/download/v" + version
	
	var filename string
	switch runtime.GOOS {
	case "windows":
		filename = fmt.Sprintf("wordma-%s-%s.exe", runtime.GOOS, runtime.GOARCH)
	case "linux", "darwin":
		filename = fmt.Sprintf("wordma-%s-%s", runtime.GOOS, runtime.GOARCH)
	default:
		return ""
	}
	
	return fmt.Sprintf("%s/%s", baseURL, filename)
}

// downloadFile downloads a file from URL and returns the temp file path
func downloadFile(url string) (string, error) {
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	resp, err := client.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("download failed with status %d", resp.StatusCode)
	}

	// Create temp file
	tempFile, err := os.CreateTemp("", "wordma-update-*")
	if err != nil {
		return "", err
	}
	defer tempFile.Close()

	// Copy response body to temp file
	_, err = io.Copy(tempFile, resp.Body)
	if err != nil {
		os.Remove(tempFile.Name())
		return "", err
	}

	return tempFile.Name(), nil
}

// copyFileForUpdate copies a file from src to dst
func copyFileForUpdate(src, dst string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	dstFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	_, err = io.Copy(dstFile, srcFile)
	return err
}

// replaceExecutable replaces the current executable with the new one
func replaceExecutable(newPath, execPath string) error {
	// On Windows, we might need to handle file locking differently
	if runtime.GOOS == "windows" {
		return replaceExecutableWindows(newPath, execPath)
	}
	return copyFileForUpdate(newPath, execPath)
}

// replaceExecutableWindows handles executable replacement on Windows
func replaceExecutableWindows(newPath, execPath string) error {
	// Try to copy directly first
	if err := copyFileForUpdate(newPath, execPath); err == nil {
		return nil
	}

	// If direct copy fails, try moving current exe and then copying
	tempPath := execPath + ".old"
	if err := os.Rename(execPath, tempPath); err != nil {
		return err
	}

	if err := copyFileForUpdate(newPath, execPath); err != nil {
		// Restore original if copy fails
		os.Rename(tempPath, execPath)
		return err
	}

	// Clean up old executable
	os.Remove(tempPath)
	return nil
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

	// 备份配置文件
	configBackupPath, hasConfig, err := backupConfigIfExists(themePath)
	if err != nil {
		utils.PrintError(fmt.Sprintf("Failed to backup config: %v", err))
		os.Exit(1)
	}
	if hasConfig {
		utils.PrintInfo("Configuration files backed up")
	}

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

	// 检查是否有配置文件更改
	hasConfigChanges, err := hasConfigFileChanges(themePath)
	if err != nil {
		utils.PrintError(fmt.Sprintf("Failed to check config changes: %v", err))
		os.Exit(1)
	}

	// 检查是否有非配置文件的更改
	hasNonConfigChanges, err := hasNonConfigFileChanges(themePath)
	if err != nil {
		utils.PrintError(fmt.Sprintf("Failed to check non-config changes: %v", err))
		os.Exit(1)
	}

	if hasLocalChanges {
		if hasConfigChanges && hasNonConfigChanges {
			// 既有配置文件更改，也有其他文件更改
			utils.PrintWarning("Theme has uncommitted local changes (including config files)")
			utils.PrintInfo("Stashing non-config changes before update...")
			err = stashNonConfigChanges(themePath)
			if err != nil {
				utils.PrintError(fmt.Sprintf("Failed to stash non-config changes: %v", err))
				os.Exit(1)
			}
			utils.PrintInfo("Non-config changes stashed successfully")
		} else if hasNonConfigChanges {
			// 只有非配置文件更改
			utils.PrintWarning("Theme has uncommitted local changes")
			utils.PrintInfo("Stashing local changes before update...")
			err = utils.RunCommandInDir(themePath, "git", "stash", "push", "-m", "wordma-cli auto stash before update")
			if err != nil {
				utils.PrintError(fmt.Sprintf("Failed to stash changes: %v", err))
				os.Exit(1)
			}
			utils.PrintInfo("Local changes stashed successfully")
		} else if hasConfigChanges {
			// 只有配置文件更改
			utils.PrintInfo("Detected config file changes - these will be protected during update")
		}
	}

	// 拉取最新代码
	utils.PrintInfo(fmt.Sprintf("Pulling latest changes from %s...", currentBranch))
	err = utils.RunCommandInDir(themePath, "git", "pull", "origin", currentBranch)
	if err != nil {
		utils.PrintError(fmt.Sprintf("Failed to pull latest changes: %v", err))
		
		// 如果拉取失败且之前有stash，尝试恢复
		if hasNonConfigChanges {
			utils.PrintInfo("Attempting to restore stashed changes...")
			restoreErr := utils.RunCommandInDir(themePath, "git", "stash", "pop")
			if restoreErr != nil {
				utils.PrintWarning("Failed to restore stashed changes. You may need to manually run 'git stash pop' in the theme directory")
			}
		}
		os.Exit(1)
	}

	// 处理配置文件恢复
	if hasConfig {
		err = handleConfigRestore(themePath, configBackupPath)
		if err != nil {
			utils.PrintError(fmt.Sprintf("Failed to handle config restore: %v", err))
		}
	}

	// 如果之前有stash，询问是否恢复
	if hasNonConfigChanges {
		utils.PrintInfo("Update completed successfully")
		utils.PrintWarning("Your non-config changes were stashed")
		utils.PrintInfo("To restore your non-config changes, run:")
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

// backupConfigIfExists 备份配置目录（如果存在）
func backupConfigIfExists(themePath string) (string, bool, error) {
	configPath := filepath.Join(themePath, "config")
	if !utils.FileExists(configPath) {
		return "", false, nil
	}

	// 创建备份目录
	backupPath := filepath.Join(themePath, ".wordma-config-backup")
	
	// 如果备份目录已存在，先删除
	if utils.FileExists(backupPath) {
		err := os.RemoveAll(backupPath)
		if err != nil {
			return "", false, fmt.Errorf("failed to remove existing backup: %v", err)
		}
	}

	// 复制配置目录到备份位置
	err := utils.CopyDirectory(configPath, backupPath)
	if err != nil {
		return "", false, fmt.Errorf("failed to backup config directory: %v", err)
	}

	return backupPath, true, nil
}

// handleConfigRestore 处理配置文件恢复
func handleConfigRestore(themePath, backupPath string) error {
	configPath := filepath.Join(themePath, "config")
	
	// 检查更新后是否有新的配置文件
	hasNewConfig := utils.FileExists(configPath)
	
	if !hasNewConfig {
		// 如果更新后没有配置目录，直接恢复备份
		err := utils.CopyDirectory(backupPath, configPath)
		if err != nil {
			return fmt.Errorf("failed to restore config: %v", err)
		}
		utils.PrintSuccess("Configuration files restored")
		return cleanupBackup(backupPath)
	}

	// 检查配置是否有变化
	hasChanges, err := hasConfigChanges(configPath, backupPath)
	if err != nil {
		return fmt.Errorf("failed to check config changes: %v", err)
	}

	if !hasChanges {
		// 配置没有变化，直接清理备份
		utils.PrintInfo("Configuration files unchanged")
		return cleanupBackup(backupPath)
	}

	// 配置有变化，询问用户选择
	utils.PrintWarning("Configuration files have been updated in the new theme version")
	utils.PrintInfo("Your options:")
	fmt.Println("  1. Keep your current configuration (recommended)")
	fmt.Println("  2. Use the new default configuration")
	fmt.Println("  3. Keep backup for manual merge")
	
	choice := getUserChoice()
	
	switch choice {
	case "1":
		// 恢复用户配置
		err := os.RemoveAll(configPath)
		if err != nil {
			return fmt.Errorf("failed to remove new config: %v", err)
		}
		err = utils.CopyDirectory(backupPath, configPath)
		if err != nil {
			return fmt.Errorf("failed to restore config: %v", err)
		}
		utils.PrintSuccess("Your configuration has been restored")
		return cleanupBackup(backupPath)
		
	case "2":
		// 使用新配置
		utils.PrintInfo("Using new default configuration")
		utils.PrintInfo(fmt.Sprintf("Your old configuration is backed up at: %s", backupPath))
		return nil
		
	case "3":
		// 保留备份供手动合并
		utils.PrintInfo("Configuration backup preserved for manual merge:")
		fmt.Printf("  Old config backup: %s\n", backupPath)
		fmt.Printf("  New config: %s\n", configPath)
		utils.PrintInfo("You can manually compare and merge the configurations")
		return nil
		
	default:
		// 默认选择1
		err := os.RemoveAll(configPath)
		if err != nil {
			return fmt.Errorf("failed to remove new config: %v", err)
		}
		err = utils.CopyDirectory(backupPath, configPath)
		if err != nil {
			return fmt.Errorf("failed to restore config: %v", err)
		}
		utils.PrintSuccess("Your configuration has been restored")
		return cleanupBackup(backupPath)
	}
}

// hasConfigChanges 检查配置是否有变化
func hasConfigChanges(newConfigPath, backupPath string) (bool, error) {
	// 简单的文件数量和名称比较
	newFiles, err := getConfigFiles(newConfigPath)
	if err != nil {
		return false, err
	}
	
	oldFiles, err := getConfigFiles(backupPath)
	if err != nil {
		return false, err
	}
	
	// 如果文件数量不同，说明有变化
	if len(newFiles) != len(oldFiles) {
		return true, nil
	}
	
	// 检查文件名是否相同
	for _, newFile := range newFiles {
		found := false
		for _, oldFile := range oldFiles {
			if newFile == oldFile {
				found = true
				break
			}
		}
		if !found {
			return true, nil
		}
	}
	
	return false, nil
}

// getConfigFiles 获取配置目录中的文件列表
func getConfigFiles(configPath string) ([]string, error) {
	var files []string
	
	err := filepath.Walk(configPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			relPath, err := filepath.Rel(configPath, path)
			if err != nil {
				return err
			}
			files = append(files, relPath)
		}
		return nil
	})
	
	return files, err
}

// getUserChoice 获取用户选择
func getUserChoice() string {
	fmt.Print("Please choose an option (1-3) [1]: ")
	
	scanner := bufio.NewScanner(os.Stdin)
	if scanner.Scan() {
		choice := strings.TrimSpace(scanner.Text())
		if choice == "" {
			return "1" // 默认选择
		}
		return choice
	}
	
	return "1" // 默认选择
}

// cleanupBackup 清理备份目录
func cleanupBackup(backupPath string) error {
	return os.RemoveAll(backupPath)
}

// hasConfigFileChanges 检查是否有配置文件更改
func hasConfigFileChanges(repoPath string) (bool, error) {
	cmd := utils.NewCommand("git", "status", "--porcelain", "config/")
	cmd.Dir = repoPath
	output, err := cmd.Output()
	if err != nil {
		return false, err
	}
	return len(output) > 0, nil
}

// hasNonConfigFileChanges 检查是否有非配置文件的更改
func hasNonConfigFileChanges(repoPath string) (bool, error) {
	// 获取所有更改
	cmd := utils.NewCommand("git", "status", "--porcelain")
	cmd.Dir = repoPath
	output, err := cmd.Output()
	if err != nil {
		return false, err
	}
	
	if len(output) == 0 {
		return false, nil
	}
	
	// 检查是否有非config目录的更改
	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		
		// 跳过状态标记，获取文件路径
		if len(line) > 3 {
			filePath := line[3:]
			// 如果不是config目录下的文件，说明有非配置文件更改
			if !strings.HasPrefix(filePath, "config/") && !strings.HasPrefix(filePath, ".wordma-config-backup/") {
				return true, nil
			}
		}
	}
	
	return false, nil
}

// stashNonConfigChanges 只stash非配置文件的更改
func stashNonConfigChanges(repoPath string) error {
	// 先添加配置文件到暂存区（保护它们不被stash）
	err := utils.RunCommandInDir(repoPath, "git", "add", "config/")
	if err != nil {
		return fmt.Errorf("failed to add config files: %v", err)
	}
	
	// stash所有其他更改
	err = utils.RunCommandInDir(repoPath, "git", "stash", "push", "-m", "wordma-cli auto stash non-config changes", "--keep-index")
	if err != nil {
		return fmt.Errorf("failed to stash non-config changes: %v", err)
	}
	
	// 将配置文件从暂存区移除（恢复到工作区）
	err = utils.RunCommandInDir(repoPath, "git", "reset", "HEAD", "config/")
	if err != nil {
		return fmt.Errorf("failed to unstage config files: %v", err)
	}
	
	return nil
}