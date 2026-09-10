package cmd

import (
	"fmt"
)

func PrintHelp() {
	fmt.Println("Aether CLI - Developer toolkit for Aether extensions and themes")
	fmt.Println("\nUsage:")
	fmt.Println("  aether [command]    (aliases: aet, aether-cli)")
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
