package config




import (
	"fmt"
	"os"
	"strings"
	"os/exec"
	
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

	if os.Getenv("PACKAGE_DIR") == "" {
		fmt.Println("Environment not initialized. Running init.bat...")

		
		ex, err1 := os.Executable()
		if err1 != nil {
			fmt.Println("Error getting executable path:", err1)
			os.Exit(1)
		}
	
	
		DriveLetter = strings.Split(ex, ":")[0] + ":"
		batPath := DriveLetter + "\\system\\config\\init.bat"




		if _, err := os.Stat(batPath); os.IsNotExist(err) {
			fmt.Printf("Init file not found: %s\n", batPath)
			os.Exit(1)
		}



		// Run the batch file to set up the environment
		cmd := exec.Command("cmd", "/C", batPath) // Adjust path if needed
		cmd.Stdout = os.Stdout // Print output
		cmd.Stderr = os.Stderr
		err := cmd.Run()
		if err != nil {
			fmt.Printf("Failed to initialize environment: %v\n", err)
			os.Exit(1)
		}
	}



	// Set the paths
	PackageDir = os.Getenv("PACKAGE_DIR")
	SoftwareYAML = os.Getenv("SOFTWARE_YAML")
	JunctionsJSON = os.Getenv("JUNCTIONS_JSON")
	AppsDir = os.Getenv("APPS_DIR")
	GamesDir = os.Getenv("GAMES_DIR")
	AppsShortcutsDir = os.Getenv("APPS_SHORTCUTS_DIR")
	GamesShortcutsDir = os.Getenv("APPS_SHORTCUTS_DIR")


	if PackageDir == "" || SoftwareYAML == "" || JunctionsJSON == "" || AppsDir == "" || GamesDir == "" || AppsShortcutsDir == "" || GamesShortcutsDir == "" {
		fmt.Println("environment variables are not set.")
		os.Exit(1)
	}


}
