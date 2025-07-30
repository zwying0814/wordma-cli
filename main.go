package main

import (
	"fmt"
	"os"

	"wordma-cli/cmd"
)

// Build information (set by ldflags during build)
var (
	Version   = "dev"
	BuildTime = "unknown"
	GitCommit = "unknown"
)

func main() {
	// Set version information for cobra
	cmd.SetVersionInfo(Version, BuildTime, GitCommit)
	
	if err := cmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}