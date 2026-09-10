package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func BinaryName() string {
	if len(os.Args) > 0 {
		base := filepath.Base(os.Args[0])
		ext := filepath.Ext(base)
		name := strings.TrimSuffix(base, ext)
		if name == "aet" || name == "aether" || name == "aether-cli" {
			return name
		}
	}
	return "aether"
}

func PrintHelp() {
	bin := BinaryName()
	fmt.Println("Aether CLI - Developer toolkit for Aether extensions and themes")
	fmt.Println("\nUsage:")
	fmt.Printf("  %s [command]    (aliases: aet, aether, aether-cli)\n", bin)
	fmt.Println("\nAvailable Commands:")
	fmt.Println("  init        Scaffolds a new Aether extension or theme (--theme)")
	fmt.Println("  validate    Validates the manifest.json of an extension or package.json of a theme")
	fmt.Println("  build       Builds and packages the extension/theme (.aex / .theme)")
	fmt.Println("  dev         Deploys extension to local Aether Launcher & watches for file changes")
}

func Execute(command string, args []string) error {
	switch command {
	case "init":
		return RunInit(args)
	case "validate":
		return RunValidate(args)
	case "build":
		return RunBuild(args)
	case "dev":
		return RunDev(args)
	case "help", "--help", "-h":
		PrintHelp()
		return nil
	default:
		return fmt.Errorf("unknown command %q", command)
	}
}
