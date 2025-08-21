package cmd

import (
	"github.com/spf13/cobra"
)

var (
	version   = "dev"
	buildTime = "unknown"
	gitCommit = "unknown"
)

var rootCmd = &cobra.Command{
	Use:   "wordma",
	Short: "Wordma CLI - A scaffolding tool for wordma static blog projects",
	Long: `Wordma CLI is a command-line tool for managing wordma static blog projects.
It provides commands for project initialization, dependency management, development, and building.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Show help if no subcommand is provided
		cmd.Help()
	},
}

// SetVersionInfo sets the version information
func SetVersionInfo(v, bt, gc string) {
	version = v
	buildTime = bt
	gitCommit = gc
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.AddCommand(doctorCmd)
	rootCmd.AddCommand(initCmd)
	rootCmd.AddCommand(installCmd)
	rootCmd.AddCommand(devCmd)
	rootCmd.AddCommand(buildCmd)
	rootCmd.AddCommand(addCmd)
	rootCmd.AddCommand(updateCmd)
	rootCmd.AddCommand(deployCmd)
}