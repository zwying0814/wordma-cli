package cmd

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "wordma",
	Short: "Wordma CLI - A scaffolding tool for wordma static blog projects",
	Long: `Wordma CLI is a command-line tool for managing wordma static blog projects.
It provides commands for project initialization, dependency management, development, and building.`,
	Version: "1.0.0",
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
}