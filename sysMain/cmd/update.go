package cmd

import "fmt"

func Update(args []string) {
	if len(args) == 0 {
		fmt.Println("Please provide a package name to update.")
		return
	}

	packageName := args[0]
	fmt.Printf("updating package: %s\n", packageName)

}