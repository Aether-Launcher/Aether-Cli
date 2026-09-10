package cmd

import (
	"fmt"
	"io"
	"os"
	"os/signal"
	"path/filepath"
	"runtime"
	"syscall"
	"time"

	"github.com/wayback09/aether-cli/pkg/manifest"
)

func RunDev(args []string) error {
	dir := "."
	if len(args) > 0 && args[0] != "" && args[0] != "--help" && args[0] != "-h" {
		dir = args[0]
	}

	cwd, err := filepath.Abs(dir)
	if err != nil {
		return fmt.Errorf("failed to resolve directory: %w", err)
	}

	// 1. Validate Extension
	fmt.Println("🔍 Validating extension...")
	m, err := manifest.Validate(cwd)
	if err != nil {
		return fmt.Errorf("❌ Extension validation failed: %w", err)
	}

	// 2. Resolve Aether Local Extension Data Directory
	destDir, err := getAetherExtensionDir(m.ID)
	if err != nil {
		return fmt.Errorf("❌ Could not locate Aether data directory: %w", err)
	}

	fmt.Println("==================================================")
	fmt.Printf("🚀 Aether Dev Mode Active\n")
	fmt.Printf("📦 Extension: %s (id: %s, v%s)\n", m.Name, m.ID, m.Version)
	fmt.Printf("📂 Source:    %s\n", cwd)
	fmt.Printf("🎯 Target:    %s\n", destDir)
	fmt.Println("==================================================")

	// Initial Sync
	if err := syncDirectory(cwd, destDir); err != nil {
		return fmt.Errorf("❌ Failed initial sync to Aether: %w", err)
	}
	fmt.Println("✅ Extension deployed into Aether Launcher!")
	fmt.Println("👀 Watching for file changes... (Press Ctrl+C to stop)")

	// 3. File Watcher Loop
	modTimes := make(map[string]time.Time)
	scanFiles(cwd, modTimes)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-sigChan:
			fmt.Println("\n👋 Stopping dev mode watcher. Happy coding!")
			return nil
		case <-ticker.C:
			changedFiles := checkModTimes(cwd, modTimes)
			if len(changedFiles) > 0 {
				nowStr := time.Now().Format("15:04:05")
				fmt.Printf("[%s] 🔄 Changed: %v\n", nowStr, changedFiles[0])
				if err := syncDirectory(cwd, destDir); err != nil {
					fmt.Printf("[%s] ⚠️ Sync error: %v\n", nowStr, err)
				} else {
					fmt.Printf("[%s] ✨ Extension updated in Aether!\n", nowStr)
				}
			}
		}
	}
}

func getAetherExtensionDir(extensionID string) (string, error) {
	// Check for local dev .aether folder first
	if stat, err := os.Stat(".aether"); err == nil && stat.IsDir() {
		abs, _ := filepath.Abs(".aether")
		return filepath.Join(abs, "extensions", extensionID), nil
	}

	userConfig, err := os.UserConfigDir()
	if err != nil {
		userHome, homeErr := os.UserHomeDir()
		if homeErr != nil {
			return "", err
		}
		if runtime.GOOS == "windows" {
			userConfig = filepath.Join(userHome, "AppData", "Roaming")
		} else {
			userConfig = filepath.Join(userHome, ".config")
		}
	}

	aetherDataDir := filepath.Join(userConfig, "Aether")
	target := filepath.Join(aetherDataDir, "extensions", extensionID)
	return target, nil
}

func scanFiles(src string, modTimes map[string]time.Time) {
	_ = filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return nil
		}
		rel, err := filepath.Rel(src, path)
		if err != nil {
			return nil
		}
		// Ignore node_modules, .git, etc.
		if rel == ".git" || filepath.HasPrefix(rel, ".git"+string(filepath.Separator)) || filepath.HasPrefix(rel, "node_modules") {
			return nil
		}
		modTimes[rel] = info.ModTime()
		return nil
	})
}

func checkModTimes(src string, modTimes map[string]time.Time) []string {
	var changed []string
	currentFiles := make(map[string]bool)

	_ = filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return nil
		}
		rel, err := filepath.Rel(src, path)
		if err != nil {
			return nil
		}
		if rel == ".git" || filepath.HasPrefix(rel, ".git"+string(filepath.Separator)) || filepath.HasPrefix(rel, "node_modules") {
			return nil
		}
		currentFiles[rel] = true
		lastTime, exists := modTimes[rel]
		if !exists || info.ModTime().After(lastTime) {
			modTimes[rel] = info.ModTime()
			changed = append(changed, rel)
		}
		return nil
	})

	for f := range modTimes {
		if !currentFiles[f] {
			delete(modTimes, f)
			changed = append(changed, f+" (deleted)")
		}
	}

	return changed
}

func syncDirectory(src, dest string) error {
	if err := os.MkdirAll(dest, 0755); err != nil {
		return err
	}

	return filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		rel, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}
		if rel == "." {
			return nil
		}
		if rel == ".git" || filepath.HasPrefix(rel, ".git"+string(filepath.Separator)) || filepath.HasPrefix(rel, "node_modules") {
			if info.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}

		targetPath := filepath.Join(dest, rel)

		if info.IsDir() {
			return os.MkdirAll(targetPath, 0755)
		}

		return copyFile(path, targetPath)
	})
}

func copyFile(src, dest string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	if err := os.MkdirAll(filepath.Dir(dest), 0755); err != nil {
		return err
	}

	out, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	return err
}
