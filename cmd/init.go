package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"wordma-cli/utils"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize a new wordma project in current directory",
	Long:  "Initialize a new wordma static blog project in the current directory by cloning the template repository",
	Args:  cobra.NoArgs,
	Run:   runInit,
}

func runInit(cmd *cobra.Command, args []string) {
	// 获取当前工作目录
	currentDir, err := os.Getwd()
	if err != nil {
		utils.PrintError(fmt.Sprintf("Failed to get current directory: %v", err))
		os.Exit(1)
	}
	
	// 检查必要的依赖
	if !utils.CheckCommand("git") {
		utils.PrintError("Git is required for project initialization")
		fmt.Printf("  %s\n", utils.GetInstallInstructions("git"))
		os.Exit(1)
	}

	// 获取当前目录名作为项目名
	projectName := filepath.Base(currentDir)
	
	// 检查当前目录是否为空（忽略可能存在的wordma可执行文件）
	isEmpty, err := utils.IsDirEmpty(currentDir, "wordma", "wordma.exe")
	if err != nil {
		utils.PrintError(fmt.Sprintf("Failed to check directory: %v", err))
		os.Exit(1)
	}
	
	if !isEmpty {
		utils.PrintError("Current directory is not empty")
		utils.PrintInfo("Please run 'wordma init' in an empty directory")
		os.Exit(1)
	}

	utils.PrintInfo(fmt.Sprintf("Initializing wordma project '%s' in current directory...", projectName))

	// 第一步：克隆远程仓库到临时目录，然后移动内容
	utils.PrintInfo("Cloning wordma template repository...")
	repoURL := "https://github.com/zwying0814/wordma.git"
	tempDir := filepath.Join(os.TempDir(), "wordma-temp")
	
	// 清理可能存在的临时目录
	if utils.FileExists(tempDir) {
		os.RemoveAll(tempDir)
	}
	
	err = utils.RunCommand("git", "clone", "-b", "main", repoURL, tempDir)
	if err != nil {
		utils.PrintError(fmt.Sprintf("Failed to clone repository: %v", err))
		os.Exit(1)
	}
	utils.PrintSuccess("Repository cloned successfully")

	// 第二步：移动文件到当前目录
	utils.PrintInfo("Moving files to current directory...")
	err = moveDirectoryContents(tempDir, currentDir)
	if err != nil {
		utils.PrintError(fmt.Sprintf("Failed to move files: %v", err))
		os.Exit(1)
	}
	
	// 清理临时目录
	os.RemoveAll(tempDir)
	utils.PrintSuccess("Files moved successfully")

	// 第三步：初始化 .deploy 目录为 git 项目
	deployPath := filepath.Join(currentDir, ".deploy")
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

	utils.PrintSuccess(".deploy directory initialized as git repository")

	// 删除原始的 .git 目录（可选）
	originalGitPath := filepath.Join(currentDir, ".git")
	if utils.FileExists(originalGitPath) {
		err = os.RemoveAll(originalGitPath)
		if err != nil {
			utils.PrintWarning(fmt.Sprintf("Failed to remove original .git directory: %v", err))
		} else {
			utils.PrintInfo("Removed original .git directory")
		}
	}

	fmt.Println()
	utils.PrintSuccess(fmt.Sprintf("Wordma project '%s' initialized successfully!", projectName))
	utils.PrintInfo("Next steps:")
	fmt.Printf("  1. wordma install\n")
	fmt.Printf("  2. wordma dev <theme-name>\n")
}

// moveDirectoryContents 移动目录内容（支持跨驱动器）
func moveDirectoryContents(src, dst string) error {
	entries, err := os.ReadDir(src)
	if err != nil {
		return err
	}
	
	for _, entry := range entries {
		srcPath := filepath.Join(src, entry.Name())
		dstPath := filepath.Join(dst, entry.Name())
		
		if entry.IsDir() {
			// 递归处理目录
			err = copyDirectory(srcPath, dstPath)
		} else {
			// 复制文件
			err = copyFile(srcPath, dstPath)
		}
		
		if err != nil {
			return err
		}
	}
	
	return nil
}



// copyFile 复制文件
func copyFile(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()
	
	destFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destFile.Close()
	
	_, err = destFile.ReadFrom(sourceFile)
	if err != nil {
		return err
	}
	
	// 复制文件权限
	sourceInfo, err := os.Stat(src)
	if err != nil {
		return err
	}
	
	return os.Chmod(dst, sourceInfo.Mode())
}

// copyDirectory 递归复制目录
func copyDirectory(src, dst string) error {
	// 创建目标目录
	err := os.MkdirAll(dst, 0755)
	if err != nil {
		return err
	}
	
	entries, err := os.ReadDir(src)
	if err != nil {
		return err
	}
	
	for _, entry := range entries {
		srcPath := filepath.Join(src, entry.Name())
		dstPath := filepath.Join(dst, entry.Name())
		
		if entry.IsDir() {
			err = copyDirectory(srcPath, dstPath)
		} else {
			err = copyFile(srcPath, dstPath)
		}
		
		if err != nil {
			return err
		}
	}
	
	return nil
}