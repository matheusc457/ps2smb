package main

import (
	"fmt"
	"os"
)

const version = "0.1.0-dev"

func main() {
	fmt.Printf("ps2smb v%s\n", version)
	fmt.Println("PlayStation 2 SMB Configuration Tool")
	fmt.Println("\nRun 'ps2smb --help' for usage information")
	os.Exit(0)
}

