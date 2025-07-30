package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"wordma-cli/utils"
)

var initCmd = &cobra.Command{
	Use:   "init <name>",
	Short: "Initialize a new wordma project",
	Long:  "Initialize a new wordma static blog project by cloning the template repository",
	Args:  cobra.ExactArgs(1),
	Run:   runInit,
}

func runInit(cmd *cobra.Command, args []string) {
	projectName := args[0]
	
	// 检查必要的依赖
	if !utils.CheckCommand("git") {
		utils.PrintError("Git is required for project initialization")
		fmt.Printf("  %s\n", utils.GetInstallInstructions("git"))
		os.Exit(1)
	}

	// 确定项目路径
	var projectPath string
	if filepath.IsAbs(projectName) {
		projectPath = projectName
	} else {
		wd, err := os.Getwd()
		if err != nil {
			utils.PrintError(fmt.Sprintf("Failed to get current directory: %v", err))
			os.Exit(1)
		}
		projectPath = filepath.Join(wd, projectName)
	}

	// 检查目录是否已存在
	if utils.FileExists(projectPath) {
		utils.PrintError(fmt.Sprintf("Directory '%s' already exists", projectPath))
		os.Exit(1)
	}

	utils.PrintInfo(fmt.Sprintf("Initializing wordma project in '%s'...", projectPath))

	// 第一步：克隆远程仓库
	utils.PrintInfo("Cloning wordma template repository...")
	repoURL := "https://github.com/zwying0814/wordma.git"
	
	err := utils.RunCommand("git", "clone", "-b", "main", repoURL, projectPath)
	if err != nil {
		utils.PrintError(fmt.Sprintf("Failed to clone repository: %v", err))
		os.Exit(1)
	}
	utils.PrintSuccess("Repository cloned successfully")

	// 第二步：初始化 .deploy 目录为 git 项目
	deployPath := filepath.Join(projectPath, ".deploy")
	utils.PrintInfo("Initializing .deploy directory as git repository...")
	
	err = utils.CreateDir(deployPath)
	if err != nil {
		utils.PrintError(fmt.Sprintf("Failed to create .deploy directory: %v", err))
		os.Exit(1)
	}

	err = utils.RunCommandInDir(deployPath, "git", "init")
	if err != nil {
		utils.PrintError(fmt.Sprintf("Failed to initialize .deploy git repository: %v", err))
		os.Exit(1)
	}
	utils.PrintSuccess(".deploy directory initialized as git repository")

	// 删除原始的 .git 目录（可选）
	originalGitPath := filepath.Join(projectPath, ".git")
	if utils.FileExists(originalGitPath) {
		err = os.RemoveAll(originalGitPath)
		if err != nil {
			utils.PrintWarning(fmt.Sprintf("Failed to remove original .git directory: %v", err))
		} else {
			utils.PrintInfo("Removed original .git directory")
		}
	}

	fmt.Println()
	utils.PrintSuccess(fmt.Sprintf("Wordma project '%s' initialized successfully!", filepath.Base(projectPath)))
	utils.PrintInfo("Next steps:")
	fmt.Printf("  1. cd %s\n", projectPath)
	fmt.Printf("  2. wordma install\n")
	fmt.Printf("  3. wordma dev <theme-name>\n")
}