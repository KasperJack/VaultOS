package utils

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
	"sysmain/internal/config"
)



func CreateShortcut(source, destination string) error {
	fmt.Println(source)
	fmt.Println(destination)

	return nil
}




// GetPathFromScript calls the PowerShell script and returns the path for the given variable
func GetPathFromScript(variable string) string {
	cmd := exec.Command("powershell", "-NoProfile", "-ExecutionPolicy", "Bypass", "-File", config.Resolver, variable)

	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Run()

	return strings.TrimSpace(out.String())
}