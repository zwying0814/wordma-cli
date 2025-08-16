package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"wordma-cli/utils"
)

var deployCmd = &cobra.Command{
	Use:   "deploy",
	Short: "Deploy related commands",
	Long:  "Commands for managing deployment directory and deployment operations",
}

var deployInitCmd = &cobra.Command{
	Use:   "init [git-url]",
	Short: "Initialize or recreate the .deploy directory",
	Long:  "Initialize or recreate the .deploy directory as a git repository with proper configuration files. Optionally specify a git URL to set as remote origin.",
	Args:  cobra.MaximumNArgs(1),
	Run:   runDeployInit,
}

func runDeployInit(cmd *cobra.Command, args []string) {
	// 获取当前工作目录
	currentDir, err := os.Getwd()
	if err != nil {
		utils.PrintError(fmt.Sprintf("Failed to get current directory: %v", err))
		os.Exit(1)
	}

	// 检查是否在 wordma 项目根目录
	if !isWordmaProject(currentDir) {
		utils.PrintError("Not in a wordma project directory")
		utils.PrintInfo("Please run this command in the root directory of a wordma project")
		utils.PrintInfo("(Directory should contain themes/ folder or package.json)")
		os.Exit(1)
	}

	// 检查必要的依赖
	if !utils.CheckCommand("git") {
		utils.PrintError("Git is required for deploy directory initialization")
		fmt.Printf("  %s\n", utils.GetInstallInstructions("git"))
		os.Exit(1)
	}

	deployPath := filepath.Join(currentDir, ".deploy")
	
	// 检查 .deploy 目录是否已存在
	if utils.FileExists(deployPath) {
		utils.PrintWarning(".deploy directory already exists")
		utils.PrintInfo("Do you want to reinitialize it? This will remove all existing content.")
		
		// 简单的确认提示
		fmt.Print("Continue? (y/N): ")
		var response string
		fmt.Scanln(&response)
		
		if response != "y" && response != "Y" && response != "yes" && response != "Yes" {
			utils.PrintInfo("Operation cancelled")
			return
		}
		
		// 删除现有目录
		utils.PrintInfo("Removing existing .deploy directory...")
		err = os.RemoveAll(deployPath)
		if err != nil {
			utils.PrintError(fmt.Sprintf("Failed to remove existing .deploy directory: %v", err))
			os.Exit(1)
		}
	}

	utils.PrintInfo("Initializing .deploy directory...")

	// 创建 .deploy 目录
	err = utils.CreateDir(deployPath)
	if err != nil {
		utils.PrintError(fmt.Sprintf("Failed to create .deploy directory: %v", err))
		os.Exit(1)
	}

	// 初始化为 git 仓库
	err = utils.RunCommandInDir(deployPath, "git", "init")
	if err != nil {
		utils.PrintError(fmt.Sprintf("Failed to initialize .deploy git repository: %v", err))
		os.Exit(1)
	}

	// 如果提供了 Git URL，设置远程仓库
	if len(args) > 0 {
		gitURL := args[0]
		utils.PrintInfo(fmt.Sprintf("Adding remote origin: %s", gitURL))
		
		err = utils.RunCommandInDir(deployPath, "git", "remote", "add", "origin", gitURL)
		if err != nil {
			utils.PrintError(fmt.Sprintf("Failed to add remote origin: %v", err))
			os.Exit(1)
		}
		
		utils.PrintSuccess("Remote origin added successfully!")
	}

	// 获取项目名称（当前目录名）
	projectName := filepath.Base(currentDir)

	// 创建 README.md 文件
	err = utils.CreateDeployReadme(deployPath, projectName)
	if err != nil {
		utils.PrintError(fmt.Sprintf("Failed to create README.md: %v", err))
		os.Exit(1)
	}

	// 创建 .gitignore 文件
	err = utils.CreateDeployGitignore(deployPath)
	if err != nil {
		utils.PrintError(fmt.Sprintf("Failed to create .gitignore: %v", err))
		os.Exit(1)
	}

	utils.PrintSuccess(".deploy directory initialized successfully!")
	utils.PrintInfo("Next steps:")
	fmt.Printf("  1. wordma build <theme-name>  # Build a theme\n")
	fmt.Printf("  2. cd .deploy                 # Navigate to deploy directory\n")
	fmt.Printf("  3. git add .                  # Stage files for deployment\n")
	fmt.Printf("  4. git commit -m \"Deploy\"     # Commit changes\n")
	
	if len(args) > 0 {
		fmt.Printf("  5. git push -u origin main    # Push to remote repository (first time)\n")
		fmt.Printf("     git push                   # Push to remote repository (subsequent times)\n")
	} else {
		fmt.Printf("  5. git remote add origin <url> # Add remote repository\n")
		fmt.Printf("  6. git push -u origin main    # Push to remote repository\n")
	}
}

// isWordmaProject 检查当前目录是否是 wordma 项目
func isWordmaProject(dir string) bool {
	// 检查是否存在 themes 目录
	themesDir := filepath.Join(dir, "themes")
	if utils.FileExists(themesDir) {
		return true
	}
	
	// 检查是否存在 package.json（可能是单主题项目）
	packageJson := filepath.Join(dir, "package.json")
	if utils.FileExists(packageJson) {
		return true
	}
	
	// 检查是否存在 wordma 相关的配置文件
	configFiles := []string{"wordma.config.js", "wordma.config.json", ".wordmarc"}
	for _, configFile := range configFiles {
		if utils.FileExists(filepath.Join(dir, configFile)) {
			return true
		}
	}
	
	return false
}



func init() {
	deployCmd.AddCommand(deployInitCmd)
}