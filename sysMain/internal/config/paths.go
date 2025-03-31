package config




import (
	"fmt"
	"os"
	"strings"
	"os/exec"
)


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





func init() {



	ex, err1 := os.Executable()
	if err1 != nil {
		fmt.Println("Error getting executable path:", err1)
		os.Exit(1)
	}


	DriveLetter = strings.Split(ex, ":")[0] + ":"
	

	if os.Getenv("Drive_Letter") != DriveLetter {
		fmt.Println("Environment not initialized. Running init.bat...")



		batfile:= DriveLetter + "\\system\\config\\init.bat"
		if _, err := os.Stat(batfile); os.IsNotExist(err) {
			fmt.Printf("init script  not found: %s\n", batfile)
			os.Exit(1)
		}

		cmd := exec.Command("cmd", "/C", batfile, "setx")
		cmd.Stdout = os.Stdout 
		cmd.Stderr = os.Stderr
		err := cmd.Run()
		if err != nil {
			fmt.Printf("Failed to initialize environment: %v\n", err)
			os.Exit(1)
		}



		fmt.Printf("run the tool again")
		os.Exit(0)



}











	PackageDir = os.Getenv("PACKAGE_DIR")
	SoftwareYAML = os.Getenv("SOFTWARE_YAML")
	JunctionsJSON =  os.Getenv("JUNCTIONS_JSON")
	AppsDir =    os.Getenv("APPS_DIR")
	GamesDir =    os.Getenv("GAMES_DIR")
	AppsShortcutsDir = os.Getenv("APPS_SHORTCUTS_DIR")
	GamesShortcutsDir =  os.Getenv("APPS_SHORTCUTS_DIR")


	if PackageDir == "" || SoftwareYAML == "" || JunctionsJSON == "" || AppsDir == "" || GamesDir == "" || AppsShortcutsDir == "" || GamesShortcutsDir == "" {
		fmt.Println("environment variables are not set.")
		os.Exit(1)
	}


}
