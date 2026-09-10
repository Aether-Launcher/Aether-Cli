package main

import (
	"fmt"
	"os"

	"github.com/wayback09/aether-cli/cmd"
)

func main() {
	if len(os.Args) < 2 {
		cmd.PrintHelp()
		os.Exit(0)
	}

	command := os.Args[1]
	args := os.Args[2:]

	if err := cmd.Execute(command, args); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}
