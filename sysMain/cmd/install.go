package cmd



import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"gopkg.in/yaml.v3"
	"sysmain/internal/config"
)


// Function to handle software installation
func Install(args []string) {


	if len(args) == 0 {
		fmt.Println("Please provide a package name to install.")
		return
	}

	name := args[0]
	fmt.Printf("Installing package: %s\n", name)




	// First, try to use the original case version for directory and file access
	softwareDir := filepath.Join(config.PackageDir, name)
	_, err := os.Stat(softwareDir)

	// If the directory doesn't exist, try case variants
	if err != nil && os.IsNotExist(err) {
		// Check if directory exists with first letter capitalized
		capitalized := strings.ToUpper(name[:1]) + name[1:]
		softwareDirCap := filepath.Join(config.PackageDir, capitalized)
		_, err = os.Stat(softwareDirCap)

		if err == nil {
			// Found with capitalized first letter
			name = capitalized
			softwareDir = softwareDirCap
		} else {
			// Check for all uppercase version
			uppercase := strings.ToUpper(name)
			softwareDirUpper := filepath.Join(config.PackageDir, uppercase)
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
	var config config.SoftwareConfig
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

}
