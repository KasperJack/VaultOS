package cmd


import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"gopkg.in/yaml.v3"
	"sysmain/internal/config"
	"sysmain/internal/utils"
)


func installPortable(softwareName, category, executable, currentPath string) {
	var destDir string
	// Determine the destination directory based on the category
	if strings.EqualFold(category, "game") {
		destDir = filepath.Join(config.GamesDir, filepath.Base(currentPath))
	} else {
		destDir = filepath.Join(config.AppsDir, filepath.Base(currentPath))
	}

	// Move the directory from currentPath to destDir
	if err := os.Rename(currentPath, destDir); err != nil {
		fmt.Printf("Error moving directory from '%s' to '%s': %v\n", currentPath, destDir, err)
		os.Exit(1)
	}

	// New executable path
	TrueexecutablePath := filepath.Join(destDir, executable)

	relativeDest := strings.TrimPrefix(destDir, config.DriveLetter)
	executablePath := filepath.Join("drive:"+relativeDest, executable)

	SoftwareYAML := config.SoftwareYAML
	// Load existing YAML file (if it exists)
	configMap := make(config.SoftwareConfig)
	yamlData, err := os.ReadFile(SoftwareYAML)
	if err == nil {
		_ = yaml.Unmarshal(yamlData, &configMap) // Ignore errors, assume empty if invalid
	}



	// Add the new software entry
	configMap[softwareName] = config.SoftwareDetails{
		Portable:   true,
		Category:   category,
		Executable: executablePath,
	}

	// Save the updated configuration back to YAML
	newYamlData, err := yaml.Marshal(&configMap)
	if err != nil {
		fmt.Printf("Error encoding YAML: %v\n", err)
		os.Exit(1)
	}

	err = os.WriteFile(SoftwareYAML, newYamlData, 0644)
	if err != nil {
		fmt.Printf("Error writing YAML file: %v\n", err)
		os.Exit(1)
	}

	if strings.EqualFold(category, "game") {
		utils.CreateShortcut(TrueexecutablePath, config.GamesShortcutsDir)

	} else {
		utils.CreateShortcut(TrueexecutablePath, config.AppsShortcutsDir)
	}

	// Print the new location of the directory
	fmt.Printf("New location: %s\n", destDir)
}