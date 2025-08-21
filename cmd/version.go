package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"wordma-cli/utils"
)

// GitHubRelease represents a GitHub release
type GitHubRelease struct {
	TagName string `json:"tag_name"`
	Name    string `json:"name"`
	Body    string `json:"body"`
	HTMLURL string `json:"html_url"`
}

// VersionInfo holds version information
type VersionInfo struct {
	Current string
	Latest  string
	NeedsUpdate bool
	ReleaseURL  string
	Changelog   string
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show version information",
	Long:  "Display current version and check for updates",
	Run:   runVersion,
}

func init() {
	// Add version command as subcommand
	rootCmd.AddCommand(versionCmd)
}

// runVersion handles the version command
func runVersion(cmd *cobra.Command, args []string) {
	versionInfo, err := getVersionInfo()
	if err != nil {
		utils.PrintError(fmt.Sprintf("Failed to check version: %v", err))
		printCurrentVersion()
		return
	}

	printVersionInfo(versionInfo)
}

// getVersionInfo retrieves current and latest version information
func getVersionInfo() (*VersionInfo, error) {
	currentVer := getCurrentVersion()
	latestVer, releaseURL, changelog, err := getLatestVersion()
	if err != nil {
		return nil, err
	}

	needsUpdate := compareVersions(currentVer, latestVer) < 0

	return &VersionInfo{
		Current:     currentVer,
		Latest:      latestVer,
		NeedsUpdate: needsUpdate,
		ReleaseURL:  releaseURL,
		Changelog:   changelog,
	}, nil
}

// getCurrentVersion returns the current version
func getCurrentVersion() string {
	if version == "dev" {
		return "dev"
	}
	// Remove 'v' prefix if present to match the format from GitHub API
	return strings.TrimPrefix(version, "v")
}

// getLatestVersion fetches the latest version from GitHub
func getLatestVersion() (string, string, string, error) {
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	resp, err := client.Get("https://api.github.com/repos/zwying0814/wordma-cli/releases/latest")
	if err != nil {
		return "", "", "", fmt.Errorf("failed to fetch latest release: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", "", "", fmt.Errorf("GitHub API returned status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", "", "", fmt.Errorf("failed to read response body: %w", err)
	}

	var release GitHubRelease
	if err := json.Unmarshal(body, &release); err != nil {
		return "", "", "", fmt.Errorf("failed to parse release data: %w", err)
	}

	// Remove 'v' prefix if present
	latestVer := strings.TrimPrefix(release.TagName, "v")
	
	return latestVer, release.HTMLURL, release.Body, nil
}

// compareVersions compares two semantic versions
// Returns: -1 if v1 < v2, 0 if v1 == v2, 1 if v1 > v2
func compareVersions(v1, v2 string) int {
	if v1 == "dev" {
		return -1 // dev is always considered older
	}
	
	if v2 == "dev" {
		return 1
	}

	// Parse version strings
	parts1 := parseVersion(v1)
	parts2 := parseVersion(v2)

	// Compare each part
	for i := 0; i < 3; i++ {
		if parts1[i] < parts2[i] {
			return -1
		}
		if parts1[i] > parts2[i] {
			return 1
		}
	}

	return 0
}

// parseVersion parses a version string into major, minor, patch numbers
func parseVersion(version string) [3]int {
	parts := strings.Split(version, ".")
	result := [3]int{0, 0, 0}

	for i, part := range parts {
		if i >= 3 {
			break
		}
		if num, err := strconv.Atoi(part); err == nil {
			result[i] = num
		}
	}

	return result
}

// printVersionInfo displays version information
func printVersionInfo(info *VersionInfo) {
	fmt.Printf("\n%s Wordma CLI Version Information\n", utils.ColorText("📦", "blue"))
	fmt.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n\n")

	// Current version
	fmt.Printf("%s Current Version: %s\n", 
		utils.ColorText("🔹", "blue"), 
		utils.ColorText(info.Current, "green"))

	// Latest version
	fmt.Printf("%s Latest Version:  %s\n", 
		utils.ColorText("🔹", "blue"), 
		utils.ColorText(info.Latest, "green"))

	// System info
	fmt.Printf("%s Platform:        %s/%s\n", 
		utils.ColorText("🔹", "blue"), 
		runtime.GOOS, runtime.GOARCH)

	fmt.Printf("\n")

	// Update status
	if info.NeedsUpdate {
		fmt.Printf("%s %s\n", 
			utils.ColorText("🚀", "yellow"), 
			utils.ColorText("A new version is available!", "yellow"))
		fmt.Printf("\n%s To update, run: %s\n", 
			utils.ColorText("💡", "blue"), 
			utils.ColorText("wordma update", "cyan"))
		fmt.Printf("%s Release page: %s\n", 
			utils.ColorText("🔗", "blue"), 
			utils.ColorText(info.ReleaseURL, "cyan"))
	} else {
		fmt.Printf("%s %s\n", 
			utils.ColorText("✅", "green"), 
			utils.ColorText("You are using the latest version!", "green"))
	}

	fmt.Printf("\n")
}

// printCurrentVersion prints only current version info (fallback)
func printCurrentVersion() {
	fmt.Printf("\n%s Wordma CLI Version: %s\n", 
		utils.ColorText("📦", "blue"), 
		utils.ColorText(getCurrentVersion(), "green"))
	fmt.Printf("%s Platform: %s/%s\n\n", 
		utils.ColorText("🔹", "blue"), 
		runtime.GOOS, runtime.GOARCH)
}