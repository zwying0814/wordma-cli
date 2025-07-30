package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"wordma-cli/utils"
)

var buildCmd = &cobra.Command{
	Use:   "build <theme-name>",
	Short: "Build a theme for production",
	Long:  "Build the specified theme for production deployment",
	Args:  cobra.ExactArgs(1),
	Run:   runBuild,
}

func runBuild(cmd *cobra.Command, args []string) {
	themeName := args[0]

	// 检查 pnpm 是否安装
	if !utils.CheckCommand("pnpm") {
		utils.PrintError("pnpm is required for building themes")
		fmt.Printf("  %s\n", utils.GetInstallInstructions("pnpm"))
		os.Exit(1)
	}

	// 获取项目根目录
	projectRoot, err := utils.GetProjectRoot()
	if err != nil {
		utils.PrintError(fmt.Sprintf("Failed to find project root: %v", err))
		os.Exit(1)
	}

	// 构建主题目录路径
	themePath := filepath.Join(projectRoot, "themes", themeName)
	
	// 检查主题目录是否存在
	if !utils.FileExists(themePath) {
		utils.PrintError(fmt.Sprintf("Theme directory '%s' does not exist", themePath))
		utils.PrintInfo("Available themes:")
		
		themesDir := filepath.Join(projectRoot, "themes")
		if utils.FileExists(themesDir) {
			entries, err := os.ReadDir(themesDir)
			if err == nil {
				for _, entry := range entries {
					if entry.IsDir() {
						fmt.Printf("  - %s\n", entry.Name())
					}
				}
			}
		}
		os.Exit(1)
	}

	// 检查 package.json 是否存在
	packageJsonPath := filepath.Join(themePath, "package.json")
	if !utils.FileExists(packageJsonPath) {
		utils.PrintError(fmt.Sprintf("No package.json found in theme directory '%s'", themePath))
		os.Exit(1)
	}

	utils.PrintInfo(fmt.Sprintf("Building theme '%s' for production...", themeName))
	
	// 执行 pnpm run build
	err = utils.RunCommandInDir(themePath, "pnpm", "run", "build")
	if err != nil {
		utils.PrintError(fmt.Sprintf("Failed to build theme: %v", err))
		os.Exit(1)
	}

	utils.PrintSuccess(fmt.Sprintf("Theme '%s' built successfully!", themeName))
}