package cmd

import (
	"fmt"
	"sysmain/internal/config"
)



func List() {


	config.InitPaths()


	fmt.Println("DriveLetter",config.DriveLetter)
	fmt.Println("PackageDir",config.PackageDir)

	fmt.Println("SoftwareYAML",config.SoftwareYAML)

	fmt.Println("JunctionsJSON",config.JunctionsJSON)

	fmt.Println("AppsDir",config.AppsDir)

	fmt.Println("GamesDir",config.GamesDir)

	fmt.Println("AppsShortcutsDir",config.AppsShortcutsDir)

	fmt.Println("GamesShortcutsDir",config.GamesShortcutsDir)


}