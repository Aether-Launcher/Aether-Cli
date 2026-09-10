package cmd

import (
	"fmt"
	"os"

	"github.com/wayback09/aether-cli/pkg/scaffold"
)

func RunInit(args []string) error {
	isTheme := false
	var filteredArgs []string
	for _, arg := range args {
		if arg == "--theme" || arg == "-t" {
			isTheme = true
		} else {
			filteredArgs = append(filteredArgs, arg)
		}
	}

	bin := BinaryName()
	if len(filteredArgs) < 2 {
		if isTheme {
			return fmt.Errorf("usage: %s init --theme <theme-name> <theme-id>", bin)
		}
		return fmt.Errorf("usage: %s init <extension-name> <extension-id> [--theme]", bin)
	}

	name := filteredArgs[0]
	id := filteredArgs[1]

	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("error getting current directory: %w", err)
	}

	targetDir := cwd + "/" + name
	if err := os.Mkdir(targetDir, 0755); err != nil {
		if os.IsExist(err) {
			return fmt.Errorf("directory %s already exists", name)
		}
		return fmt.Errorf("failed to create directory: %w", err)
	}

	if isTheme {
		fmt.Printf("Initializing theme '%s' (%s)...\n", name, id)
		if err := scaffold.CreateTheme(targetDir, name, id); err != nil {
			return fmt.Errorf("failed to initialize theme: %w", err)
		}
		fmt.Println("Initialization complete.")
		fmt.Printf("Next steps:\n  cd %s && %s build --theme\n", name, bin)
	} else {
		fmt.Printf("Initializing extension '%s' (%s)...\n", name, id)
		if err := scaffold.CreateExtension(targetDir, name, id); err != nil {
			return fmt.Errorf("failed to initialize extension: %w", err)
		}
		fmt.Println("Initialization complete.")
		fmt.Printf("Next steps:\n  cd %s && %s build\n", name, bin)
	}
	return nil
}
