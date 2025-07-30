package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"wordma-cli/utils"
)

var devCmd = &cobra.Command{
	Use:   "dev <theme-name>",
	Short: "Start development server for a theme",
	Long:  "Start the development server for the specified theme in the themes directory",
	Args:  cobra.ExactArgs(1),
	Run:   runDev,
}

func runDev(cmd *cobra.Command, args []string) {
	themeName := args[0]

	// 检查 pnpm 是否安装
	if !utils.CheckCommand("pnpm") {
		utils.PrintError("pnpm is required for running development server")
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

	utils.PrintInfo(fmt.Sprintf("Starting development server for theme '%s'...", themeName))
	
	// 执行 pnpm run dev
	err = utils.RunCommandInDir(themePath, "pnpm", "run", "dev")
	if err != nil {
		utils.PrintError(fmt.Sprintf("Failed to start development server: %v", err))
		os.Exit(1)
	}
}