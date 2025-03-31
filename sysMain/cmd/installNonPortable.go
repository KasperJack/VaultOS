package cmd

import "fmt"

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
