package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func main() {
	if os.Geteuid() == 0 {
		fmt.Println("Error: This script should not be run with sudo privileges.")
		os.Exit(1)
	}

	sourceBinary := "/home/gamer/.local/bin/vinegar"
	destinationLink := "/usr/local/bin/vinegar"

	if _, err := os.Stat(destinationLink); err == nil {
		fmt.Printf("Error: %s already exists.\n", destinationLink)
		fmt.Print("Do you want to remove it? (yes/no): ")

		reader := bufio.NewReader(os.Stdin)
		response, _ := reader.ReadString('\n')
		response = strings.TrimSpace(response)

		if response != "yes" {
			fmt.Println("Installation aborted.")
			os.Exit(1)
		}

		// Use sudo for removing the existing link
		cmd := exec.Command("sudo", "rm", destinationLink)
		err := cmd.Run()

		if err != nil {
			fmt.Printf("Error: Failed to remove existing link %s\n", destinationLink)
			os.Exit(1)
		}

		fmt.Printf("Removed existing link: %s\n", destinationLink)
	}

	// Create the symbolic link
	cmd := exec.Command("sudo", "ln", "-s", sourceBinary, destinationLink)
	err := cmd.Run()

	if err != nil {
		fmt.Printf("Error: Failed to create symbolic link %s\n", destinationLink)
		os.Exit(1)
	}

	fmt.Printf("Symbolic link created: %s -> %s\n", sourceBinary, destinationLink)
}
