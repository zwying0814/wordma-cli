package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"wordma-cli/utils"
)

var installCmd = &cobra.Command{
	Use:     "install",
	Aliases: []string{"i"},
	Short:   "Install project dependencies",
	Long:    "Install all dependencies for the monorepo using pnpm",
	Run:     runInstall,
}

func runInstall(cmd *cobra.Command, args []string) {
	// 检查 pnpm 是否安装
	if !utils.CheckCommand("pnpm") {
		utils.PrintError("pnpm is required for dependency installation")
		fmt.Printf("  %s\n", utils.GetInstallInstructions("pnpm"))
		os.Exit(1)
	}

	// 获取项目根目录
	projectRoot, err := utils.GetProjectRoot()
	if err != nil {
		utils.PrintError(fmt.Sprintf("Failed to find project root: %v", err))
		os.Exit(1)
	}

	// 检查是否在 wordma 项目中
	packageJsonPath := fmt.Sprintf("%s/package.json", projectRoot)
	if !utils.FileExists(packageJsonPath) {
		utils.PrintError("No package.json found. Make sure you're in a wordma project directory.")
		os.Exit(1)
	}

	utils.PrintInfo("Installing dependencies with pnpm...")
	
	err = utils.RunCommandInDir(projectRoot, "pnpm", "install")
	if err != nil {
		utils.PrintError(fmt.Sprintf("Failed to install dependencies: %v", err))
		os.Exit(1)
	}

	utils.PrintSuccess("Dependencies installed successfully!")
}