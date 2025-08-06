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
	err = createDeployReadme(deployPath, projectName)
	if err != nil {
		utils.PrintError(fmt.Sprintf("Failed to create README.md: %v", err))
		os.Exit(1)
	}

	// 创建 .gitignore 文件
	err = createDeployGitignore(deployPath)
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

// createDeployReadme 创建 .deploy 目录的 README.md 文件
func createDeployReadme(deployPath, projectName string) error {
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

// createDeployGitignore 创建 .deploy 目录的 .gitignore 文件
func createDeployGitignore(deployPath string) error {
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