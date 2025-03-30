package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	//"os/exec"
	"gopkg.in/yaml.v3"
)

// Configuration paths
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

// SoftwareConfig represents the structure of the YAML configuration file
type SoftwareConfig map[string]SoftwareDetails

// SoftwareDetails holds the details of a software package
type SoftwareDetails struct {
	Portable   bool     `yaml:"portable"`
	Category   string   `yaml:"category"`
	Executable string   `yaml:"executable"`
	Junctions  []string `yaml:"junctions,omitempty"`
}

// Initialize paths with the correct drive letter
func initPaths() {
	// Get the current executable path to determine the drive letter
	ex, err := os.Executable()
	if err != nil {
		fmt.Println("Error getting executable path:", err)
		os.Exit(1)
	}

	// Extract drive letter from the executable path
	if runtime.GOOS == "windows" {
		DriveLetter = strings.Split(ex, ":")[0] + ":"
	} else {
		// For non-Windows systems, use a default or empty string
		DriveLetter = ""
	}

	// Set the paths with the detected drive letter
	PackageDir = filepath.Join(DriveLetter, "\\system\\package")
	SoftwareYAML = filepath.Join(DriveLetter, "\\system\\config\\software.yaml")
	JunctionsJSON = filepath.Join(DriveLetter, "\\system\\config\\junctions.json")
	AppsDir = filepath.Join(DriveLetter, "\\system\\software\\apps")
	GamesDir = filepath.Join(DriveLetter, "\\system\\software\\games")
	AppsShortcutsDir = filepath.Join(DriveLetter, "\\apps")
	GamesShortcutsDir = filepath.Join(DriveLetter, "\\games")
}

func main() {
	// Initialize paths
	initPaths()

	// Define command-line flags
	installCmd := flag.NewFlagSet("install", flag.ExitOnError)

	// Parse command-line arguments
	if len(os.Args) < 2 {
		fmt.Println("Expected 'install' subcommand")
		os.Exit(1)
	}

	// Handle different subcommands
	switch os.Args[1] {
	case "install":
		installCmd.Parse(os.Args[2:])
		if installCmd.NArg() < 1 {
			fmt.Println("Expected software name")
			os.Exit(1)
		}
		softwareName := installCmd.Arg(0)
		installSoftware(softwareName)
	default:
		fmt.Println("Expected 'install' subcommand")
		os.Exit(1)
	}
}

func installPortable(softwareName, category, executable, currentPath string) {
	var destDir string
	// Determine the destination directory based on the category
	if strings.EqualFold(category, "game") {
		destDir = filepath.Join(GamesDir, filepath.Base(currentPath))
	} else {
		destDir = filepath.Join(AppsDir, filepath.Base(currentPath))
	}

	// Move the directory from currentPath to destDir
	if err := os.Rename(currentPath, destDir); err != nil {
		fmt.Printf("Error moving directory from '%s' to '%s': %v\n", currentPath, destDir, err)
		os.Exit(1)
	}

	// New executable path
	TrueexecutablePath := filepath.Join(destDir, executable)

	relativeDest := strings.TrimPrefix(destDir, DriveLetter)
	executablePath := filepath.Join("drive:"+relativeDest, executable)

	// Load existing YAML file (if it exists)
	config := make(SoftwareConfig)
	yamlData, err := os.ReadFile(SoftwareYAML)
	if err == nil {
		_ = yaml.Unmarshal(yamlData, &config) // Ignore errors, assume empty if invalid
	}

	// Add the new software entry
	config[softwareName] = SoftwareDetails{
		Portable:   true,
		Category:   category,
		Executable: executablePath,
	}

	// Save the updated configuration back to YAML
	newYamlData, err := yaml.Marshal(&config)
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
		createShortcut(TrueexecutablePath, GamesShortcutsDir)

	} else {
		createShortcut(TrueexecutablePath, AppsShortcutsDir)
	}

	// Print the new location of the directory
	fmt.Printf("New location: %s\n", destDir)
}

func createShortcut(source, destination string) error {
	fmt.Println(source)
	fmt.Println(destination)

	return nil
}

// installNonPortable processes a non-portable software installation.
// It prints out the software name, category, executable, the current directory, and each junction.
func installNonPortable(softwareName, category, executable, currentPath string, junctions []string) {
	fmt.Println("Installing Non-Portable Software:")
	fmt.Printf("  Software Name: %s\n", softwareName)
	fmt.Printf("  Category: %s\n", category)
	fmt.Printf("  Executable: %s\n", executable)
	fmt.Printf("  YAML Directory: %s\n", currentPath)
	if len(junctions) > 0 {
		fmt.Println("  Junctions:")
		for i, junction := range junctions {
			fmt.Printf("    %d: %s\n", i+1, junction)
		}
	} else {
		fmt.Println("  No junctions provided.")
	}
}

// ... (previous code remains unchanged)

// Function to handle software installation
func installSoftware(name string) {
	fmt.Printf("Installing software: %s\n", name)

	// First, try to use the original case version for directory and file access
	softwareDir := filepath.Join(PackageDir, name)
	_, err := os.Stat(softwareDir)

	// If the directory doesn't exist, try case variants
	if err != nil && os.IsNotExist(err) {
		// Check if directory exists with first letter capitalized
		capitalized := strings.ToUpper(name[:1]) + name[1:]
		softwareDirCap := filepath.Join(PackageDir, capitalized)
		_, err = os.Stat(softwareDirCap)

		if err == nil {
			// Found with capitalized first letter
			name = capitalized
			softwareDir = softwareDirCap
		} else {
			// Check for all uppercase version
			uppercase := strings.ToUpper(name)
			softwareDirUpper := filepath.Join(PackageDir, uppercase)
			_, err = os.Stat(softwareDirUpper)

			if err == nil {
				// Found with all uppercase
				name = uppercase
				softwareDir = softwareDirUpper
			}
		}
	}

	// Check if the software package directory exists after potential case adjustments
	dirInfo, err := os.Stat(softwareDir)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Printf("Error: Software package directory '%s' does not exist\n", softwareDir)
		} else {
			fmt.Printf("Error checking software directory: %v\n", err)
		}
		os.Exit(1)
	}

	// Verify it's a directory
	if !dirInfo.IsDir() {
		fmt.Printf("Error: '%s' is not a directory\n", softwareDir)
		os.Exit(1)
	}

	fmt.Printf("Found software package directory: %s\n", softwareDir)

	// Check for the yaml file matching the software name
	yamlFile := filepath.Join(softwareDir, name+".yaml")
	_, err = os.Stat(yamlFile)
	if err != nil {
		if os.IsNotExist(err) {
			// Try different case variants for the yaml file
			variants := []string{
				strings.ToLower(name) + ".yaml",
				strings.ToUpper(name) + ".yaml",
				strings.ToUpper(name[:1]) + name[1:] + ".yaml",
			}

			found := false
			for _, variant := range variants {
				testPath := filepath.Join(softwareDir, variant)
				_, err = os.Stat(testPath)
				if err == nil {
					yamlFile = testPath
					found = true
					break
				}
			}

			if !found {
				fmt.Printf("Error: YAML file for '%s' does not exist\n", name)
				os.Exit(1)
			}
		} else {
			fmt.Printf("Error checking YAML file: %v\n", err)
			os.Exit(1)
		}
	}

	fmt.Printf("Found software configuration file: %s\n", yamlFile)

	// Read and parse the YAML file
	yamlData, err := os.ReadFile(yamlFile)
	if err != nil {
		fmt.Printf("Error reading YAML file: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("YAML file content read successfully\n")

	// Parse the YAML content
	var config SoftwareConfig
	err = yaml.Unmarshal(yamlData, &config)
	if err != nil {
		fmt.Printf("Error parsing YAML file: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("YAML file parsed successfully\n")

	// Try to find the software name in the YAML, checking different case variants
	// First, try direct match
	details, exists := config[name]

	// If not found, try case-insensitive search
	if !exists {
		actualKey := ""
		for key := range config {
			if strings.EqualFold(key, name) {
				actualKey = key
				details = config[key]
				exists = true
				break
			}
		}

		if exists {
			fmt.Printf("Found software '%s' in YAML as '%s'\n", name, actualKey)
		}
	} else {
		fmt.Printf("Software '%s' found in YAML configuration\n", name)
	}

	// Exit if software not found after case-insensitive search
	if !exists {
		fmt.Printf("Error: Software name '%s' not found in YAML configuration\n", name)
		os.Exit(1)
	}

	// Check for required fields
	if details.Executable == "" {
		fmt.Printf("Error: No executable specified for '%s'\n", name)
		os.Exit(1)
	}

	fmt.Printf("Executable: %s\n", details.Executable)
	fmt.Printf("Portable: %t\n", details.Portable)
	fmt.Printf("Category: %s\n", details.Category)

	// Validate junctions based on the portability
	if details.Portable {
		// Portable software must NOT have junctions
		if len(details.Junctions) > 0 {
			fmt.Printf("Error: Portable software '%s' should not have junctions defined\n", name)
			os.Exit(1)
		}
	} else {
		// Non-portable software MUST have at least one junction
		if len(details.Junctions) == 0 {
			fmt.Printf("Error: Non-portable software '%s' must have at least one junction defined\n", name)
			os.Exit(1)
		}
	}

	// Determine which installation function to use based on the portable flag
	if details.Portable {
		installPortable(name, details.Category, details.Executable, softwareDir)
	} else {
		installNonPortable(name, details.Category, details.Executable, softwareDir, details.Junctions)
	}

	// Placeholder for further installation logic
	fmt.Println("YAML validation complete - Ready to install software")
}

// ... (rest of the code remains unchanged)
