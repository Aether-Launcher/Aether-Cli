package scaffold

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/Aether-Launcher/aether-cli/pkg/manifest"
)

func CreateExtension(dirPath string, name string, id string) error {
	// Create main directories
	if err := os.MkdirAll(filepath.Join(dirPath, "ui"), 0755); err != nil {
		return err
	}

	// Create manifest.json
	m := manifest.Manifest{
		ID:          id,
		Name:        name,
		Version:     "1.0.0",
		Author:      "Author",
		Description: "A new Aether extension.",
		Main:        "main.js",
		API:         "1.0",
		Permissions: []string{"ui:sidebar"},
	}

	manifestBytes, _ := json.MarshalIndent(m, "", "  ")
	if err := os.WriteFile(filepath.Join(dirPath, "manifest.json"), manifestBytes, 0644); err != nil {
		return err
	}

	// Create main.js
	mainJs := `// Backend Script (Goja Sandbox)
Aether.ui.registerSidebarPage({
	id: "my-page",
	label: "` + name + `",
	url: "ui/index.html"
});
`
	if err := os.WriteFile(filepath.Join(dirPath, "main.js"), []byte(mainJs), 0644); err != nil {
		return err
	}

	// Create ui/index.html
	indexHtml := `<!DOCTYPE html>
<html>
<head>
	<meta charset="utf-8">
	<title>` + name + `</title>
	<style>
		body { font-family: sans-serif; color: white; background: #1a1a1a; padding: 20px; }
	</style>
</head>
<body>
	<h1>Hello from ` + name + `</h1>
	<p>Your extension frontend is working.</p>
</body>
</html>
`
	if err := os.WriteFile(filepath.Join(dirPath, "ui", "index.html"), []byte(indexHtml), 0644); err != nil {
		return err
	}

	// Create package.json for typings
	pkgJson := `{
  "name": "` + id + `",
  "version": "1.0.0",
  "dependencies": {
    "@aethermc/sdk": "latest"
  }
}
`
	if err := os.WriteFile(filepath.Join(dirPath, "package.json"), []byte(pkgJson), 0644); err != nil {
		return err
	}

	fmt.Printf("Successfully scaffolded extension: %s\n", name)
	return nil
}

func CreateTheme(dirPath string, name string, id string) error {
	// Create package.json for theme
	pkgJson := `{
  "id": "` + id + `",
  "name": "` + name + `",
  "version": "1.0.0",
  "author": "Author",
  "description": "A new Aether theme.",
  "css": "theme.css",
  "overwrite": "overwrite.json"
}
`
	if err := os.WriteFile(filepath.Join(dirPath, "package.json"), []byte(pkgJson), 0644); err != nil {
		return err
	}

	// Create theme.css
	themeCss := `/* theme.css */
:root {
  --bg-color: #0a0a12;
  --sidebar-bg: #10101c;
  --panel-bg: #16162a;

  --accent-color: #7c3aed;
  --accent-hover: #6d28d9;

  --border-radius: 10px;
}
`
	if err := os.WriteFile(filepath.Join(dirPath, "theme.css"), []byte(themeCss), 0644); err != nil {
		return err
	}

	// Create overwrite.json
	overwriteJson := `{}`
	if err := os.WriteFile(filepath.Join(dirPath, "overwrite.json"), []byte(overwriteJson), 0644); err != nil {
		return err
	}

	// Create README.md
	readmeMd := `# ` + name + ` Theme

A new Aether launcher appearance pack.

## Installation
Build the theme package using the Aether CLI:
` + "```bash\n" + `aether-cli build --theme
` + "```\n" + `
Then install the produced ` + "`.theme`" + ` file in **Settings → Appearance** inside the Aether Launcher.
`
	if err := os.WriteFile(filepath.Join(dirPath, "README.md"), []byte(readmeMd), 0644); err != nil {
		return err
	}

	fmt.Printf("Successfully scaffolded theme: %s\n", name)
	return nil
}
