package main

import (
    "fmt"
    "os"
    "os/exec"
)

func main() {
    // Check if our environment variables are set
    portDrive := os.Getenv("A_DRIVE")
    if portDrive == "" {
        fmt.Println("Environment not initialized. Initializing...")
        
        // Run the batch file to set up environment
        cmd := exec.Command("cmd", "/c", ".\\init.bat")
        err := cmd.Run()
        if err != nil {
            fmt.Printf("Failed to initialize environment: %v\n", err)
            os.Exit(1)
        }
        
        // Note: The current process won't see the new env vars
        // You might need to restart the program or load them directly
        fmt.Println("Environment initialized.")
		fmt.Printf("Using VAR system at %s\n", portDrive)
        os.Exit(0)
    }
    
    // Environment is set, continue with normal operation
    fmt.Printf("Using portable system at %s\n", portDrive)
    
    // Your program logic here...
}