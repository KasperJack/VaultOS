package config




import (
	"fmt"
	"os"
	"strings"
	
)


var (
	DriveLetter       string
	Resolver          string

	PackageDir        string
	SoftwareYAML      string
	JunctionsJSON     string
	AppsDir           string
	GamesDir          string
	AppsShortcutsDir  string
	GamesShortcutsDir string
)



func init() {




	ex, err1 := os.Executable()
	if err1 != nil {
		fmt.Println("Error getting executable path:", err1)
		os.Exit(1)
	}


	DriveLetter = strings.Split(ex, ":")[0] + ":"
	Resolver = DriveLetter + "\\system\\config\\resolve.ps1"




	if _, err := os.Stat(Resolver); os.IsNotExist(err) {
		fmt.Printf("Resolver  not found: %s\n", Resolver)
		os.Exit(1)
	}


}




func InitPaths() {

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
