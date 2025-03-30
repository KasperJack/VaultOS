package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Global variables for paths
var (
	DriveLetter       string
	PackageDir        string
	SoftwareYAML      string
	JunctionsJSON     string
	AppsDir           string
	GamesDir          string
	AppsShortcutsDir  string
	GamesShortcutsDir string
)

// Initialize paths with the correct drive letter
func InitPaths() {
	ex, err := os.Executable()
	if err != nil {
		fmt.Println("Error getting executable path:", err)
		os.Exit(1)
	}


	DriveLetter = strings.Split(ex, ":")[0] + ":"


	// Set the paths
	PackageDir = filepath.Join(DriveLetter, "\\system\\package")
	SoftwareYAML = filepath.Join(DriveLetter, "\\system\\config\\software.yaml")
	JunctionsJSON = filepath.Join(DriveLetter, "\\system\\config\\junctions.json")
	AppsDir = filepath.Join(DriveLetter, "\\system\\software\\apps")
	GamesDir = filepath.Join(DriveLetter, "\\system\\software\\games")
	AppsShortcutsDir = filepath.Join(DriveLetter, "\\apps")
	GamesShortcutsDir = filepath.Join(DriveLetter, "\\games")
}
