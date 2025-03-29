package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"gopkg.in/yaml.v3"

)

type Software struct {
	Portable   bool     `yaml:"portable"`
	Category   string   `yaml:"category"`
	Executable string   `yaml:"executable"`
	Junctions  []string `yaml:"junctions,omitempty"`
}

type JunctionEntry struct {
	Paths []struct {
		SourcePath string `json:"SourcePath"`
		TargetPath string `json:"TargetPath"`
	} `json:"Paths"`
}

func getDriveLetter() string {
	cwd, err := os.Getwd()
	if err != nil {
		fmt.Println("Error getting current directory:", err)
		os.Exit(1)
	}
	return strings.ToUpper(cwd[:2])
}

func installSoftware(name string) {
	driveLetter := getDriveLetter()
	packageDir := filepath.Join(driveLetter, "system", "package")
	softwareFile := filepath.Join(driveLetter, "system", "config", "software.yaml")
	junctionFile := filepath.Join(driveLetter, "system", "config", "junctions.json")
	appsDir := filepath.Join(driveLetter, "system", "software", "apps")
	gamesDir := filepath.Join(driveLetter, "system", "software", "games")

	softwarePath := filepath.Join(packageDir, name, name+".yaml")
	data, err := os.ReadFile(softwarePath)
	if err != nil {
		fmt.Println("Error reading software YAML:", err)
		return
	}

	var software Software
	if err := yaml.Unmarshal(data, &software); err != nil {
		fmt.Println("Error parsing YAML:", err)
		return
	}

	destDir := appsDir
	if software.Category == "Games" {
		destDir = gamesDir
	}
	installPath := filepath.Join(destDir, name)
	if err := os.Rename(filepath.Join(packageDir, name), installPath); err != nil {
		fmt.Println("Error moving software:", err)
		return
	}

	software.Executable = filepath.Join(installPath, software.Executable)

	softwareData, err := yaml.Marshal(&software)
	if err != nil {
		fmt.Println("Error marshaling YAML:", err)
		return
	}
	if err := os.WriteFile(softwareFile, softwareData, 0644); err != nil {
		fmt.Println("Error writing software YAML:", err)
		return
	}

	if !software.Portable {
		junctions := map[string]JunctionEntry{}
		jsonData, err := os.ReadFile(junctionFile)
		if err == nil {
			json.Unmarshal(jsonData, &junctions)
		}

		var paths []struct {
			SourcePath string `json:"SourcePath"`
			TargetPath string `json:"TargetPath"`
		}
		for _, j := range software.Junctions {
			paths = append(paths, struct {
				SourcePath string `json:"SourcePath"`
				TargetPath string `json:"TargetPath"`
			}{SourcePath: j, TargetPath: filepath.Join(installPath, "data", filepath.Base(j))})
		}
		junctions[name] = JunctionEntry{Paths: paths}

		newJsonData, err := json.MarshalIndent(junctions, "", "  ")
		if err != nil {
			fmt.Println("Error marshaling JSON:", err)
			return
		}
		if err := os.WriteFile(junctionFile, newJsonData, 0644); err != nil {
			fmt.Println("Error writing junction JSON:", err)
			return
		}
	}

	fmt.Println(name, "installed successfully!")
}

func main() {
	if len(os.Args) < 3 || os.Args[1] != "install" {
		fmt.Println("Usage: tool install <software_name>")
		return
	}
	installSoftware(os.Args[2])
}
