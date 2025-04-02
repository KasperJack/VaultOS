package cmd

import "fmt"

func Remove(args []string) {
	if len(args) == 0 {
		fmt.Println("Please provide a package name to remove.")
		return
	}

	packageName := args[0]
	fmt.Printf("removing package: %s\n", packageName)

}