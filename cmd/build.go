package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/wayback09/aether-cli/pkg/manifest"
	"github.com/wayback09/aether-cli/pkg/packager"
)

func RunBuild(args []string) error {
	isTheme := false
	for _, arg := range args {
		if arg == "--theme" || arg == "-t" {
			isTheme = true
		}
	}

	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("error getting current directory: %w", err)
	}

	// Auto-detect
	if !isTheme {
		manifestPath := filepath.Join(cwd, "manifest.json")
		if _, err := os.Stat(manifestPath); os.IsNotExist(err) {
			pkgPath := filepath.Join(cwd, "package.json")
			if _, err := os.Stat(pkgPath); err == nil {
				isTheme = true
			}
		}
	}

	if isTheme {
		fmt.Println("Validating theme...")
		t, err := manifest.ValidateTheme(cwd)
		if err != nil {
			return fmt.Errorf("theme validation failed: %w", err)
		}

		fmt.Printf("Packaging theme %s v%s...\n", t.Name, t.Version)
		outPath, err := packager.BuildTheme(cwd, t)
		if err != nil {
			return fmt.Errorf("theme build failed: %w", err)
		}

		fmt.Printf("[OK] Theme build successful: %s\n", outPath)
	} else {
		fmt.Println("Validating extension...")
		m, err := manifest.Validate(cwd)
		if err != nil {
			return fmt.Errorf("extension validation failed: %w", err)
		}

		fmt.Printf("Packaging extension %s v%s...\n", m.Name, m.Version)
		outPath, err := packager.Build(cwd, m)
		if err != nil {
			return fmt.Errorf("extension build failed: %w", err)
		}

		fmt.Printf("[OK] Extension build successful: %s\n", outPath)
	}

	return nil
}
