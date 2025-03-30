package main

import (
	"fmt"
	"os"
	"sysmain/cmd"
)

func main() {


	if len(os.Args) < 2 {
		fmt.Println("Usage: pkgmgr <command> [options]")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "install":
		cmd.Install(os.Args[2:])
	case "remove":
		cmd.Remove(os.Args[2:])
	case "update":
		cmd.Update(os.Args[2:])
	case "list":
		cmd.List()
	default:
		fmt.Println("Unknown command. Available commands: install, remove, update, list")
		os.Exit(1)
	}
}












