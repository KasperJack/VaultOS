package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)




// DriveExec
// AppLauncher
// Run


// Phantom → A module that redirects commands dynamically.

func main() {

	if len(os.Args) < 2 {
		fmt.Println("Usage:")
		fmt.Println(" run list         - list available software")
		fmt.Println(" run <software>   - launch the software")
		os.Exit(1)
	}

	// drive letter.
	wd, err := os.Getwd()
	if err != nil || len(wd) < 1 {
		fmt.Println("Error: Unable to determine current working directory.")
		os.Exit(1)
	}
	driveLetter := string(wd[0])
	
	// Build the directory path where batch files are stored.
	batDir := fmt.Sprintf(`%s:\System\scripts\bat`, driveLetter)
	
	command := os.Args[1]

	// "list" command: list all available batch files (without the .bat extension)
	if strings.ToLower(command) == "list" {
		files, err := os.ReadDir(batDir)
		if err != nil {
			fmt.Printf("Error: Unable to read directory %s\n", batDir)
			os.Exit(1)
		}
		fmt.Println("Available software:")
		for _, file := range files {
			if !file.IsDir() && strings.HasSuffix(file.Name(), ".bat") {
				name := file.Name()
				if len(name) > 4 {
					name = name[:len(name)-4]
				}
				fmt.Println(name)
			}
		}
		os.Exit(0)
	}


	requested := os.Args[1]
	var matchedFile string

	files, err := os.ReadDir(batDir)
	if err != nil {
		fmt.Printf("Error: Unable to read directory %s\n", batDir)
		os.Exit(1)
	}

	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".bat") {
			baseName := file.Name()[:len(file.Name())-4] // remove .bat
			if strings.EqualFold(baseName, requested) {
				matchedFile = file.Name()
				break
			}
		}
	}

	if matchedFile == "" {
		fmt.Printf("Error: Software '%s' not found.\n", requested)
		os.Exit(1)
	}

	// Construct the full path using the matched file name.
	batFile := filepath.Join(batDir, matchedFile)

	// Execute the batch file.
	cmd := exec.Command(batFile)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Printf("Error launching '%s': %v\n", matchedFile, err)
		os.Exit(1)
	}
}
