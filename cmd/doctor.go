package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"wordma-cli/utils"
)

var doctorCmd = &cobra.Command{
	Use:   "doctor",
	Short: "Check system dependencies",
	Long:  "Check if required tools (nodejs, pnpm, git) are installed on the system",
	Run:   runDoctor,
}

func runDoctor(cmd *cobra.Command, args []string) {
	utils.PrintInfo("Checking system dependencies...")
	fmt.Println()

	allGood := true

	// 检查 Node.js
	if utils.CheckCommand("node") {
		version, err := utils.GetCommandVersion("node", "--version")
		if err != nil {
			utils.PrintWarning("Node.js is installed but version check failed")
		} else {
			utils.PrintSuccess(fmt.Sprintf("Node.js %s", version))
		}
	} else {
		utils.PrintError("Node.js is not installed")
		fmt.Printf("  %s\n", utils.GetInstallInstructions("nodejs"))
		allGood = false
	}

	// 检查 pnpm
	if utils.CheckCommand("pnpm") {
		version, err := utils.GetCommandVersion("pnpm", "--version")
		if err != nil {
			utils.PrintWarning("pnpm is installed but version check failed")
		} else {
			utils.PrintSuccess(fmt.Sprintf("pnpm %s", version))
		}
	} else {
		utils.PrintError("pnpm is not installed")
		fmt.Printf("  %s\n", utils.GetInstallInstructions("pnpm"))
		allGood = false
	}

	// 检查 Git
	if utils.CheckCommand("git") {
		version, err := utils.GetCommandVersion("git", "--version")
		if err != nil {
			utils.PrintWarning("Git is installed but version check failed")
		} else {
			utils.PrintSuccess(version)
		}
	} else {
		utils.PrintError("Git is not installed")
		fmt.Printf("  %s\n", utils.GetInstallInstructions("git"))
		allGood = false
	}

	fmt.Println()
	if allGood {
		utils.PrintSuccess("All dependencies are installed! You're ready to use wordma.")
	} else {
		utils.PrintError("Some dependencies are missing. Please install them before using wordma.")
	}
}