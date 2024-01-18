package main

import (
	"fmt"
	"os"
)

func main() {
	if os.Geteuid() == 0 {
		fmt.Println("Error: Running as root is not allowed.")
		os.Exit(1)
	}
}
