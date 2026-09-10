# Aether CLI

<p align="center">
  <a href="https://discord.gg/pQc9NnGhpG">
    <img src="https://img.shields.io/badge/discord-Join%20our%20Discord-5865F2?logo=discord&logoColor=white&style=for-the-badge" alt="Discord">
  </a>
</p>

`aether` (alias: `aet`) is the official developer toolkit for creating, building, developing, and packaging **extensions** and **themes** for the [Aether Minecraft Launcher](https://github.com/Aether-Launcher/Aether). It lets developers scaffold new projects, hot-reload them into the local launcher, and package them securely into `.aex` (extension) or `.theme` (appearance pack) format.

---

## Installation

Install the CLI with Go (this installs `aether`, the short alias `aet`, and `aether-cli`):

```bash
go install github.com/Aether-Launcher/aether-cli/...@latest
```

Make sure your Go bin directory is on `PATH` (`$(go env GOPATH)/bin` on macOS/Linux, `%USERPROFILE%\go\bin` on Windows), then open a new terminal and run:

```bash
aether help
# or using the short alias:
aet help
```

---

## Commands

You can run any command using either `aether`, `aet`, or `aether-cli`.

### `init` — Scaffold a new project

```bash
# Scaffold a new extension
aether init <name> <id>

# Scaffold a new theme
aether init <name> <id> --theme
```

| Flag | Alias | Description |
|:---|:---|:---|
| `--theme` | `-t` | Scaffold a theme instead of an extension |

**Extension output:**
```
<name>/
├── manifest.json     # Extension metadata (id, name, version, permissions, api)
├── package.json      # npm typings dep (@aethermc/sdk)
├── main.js           # Backend sandbox entry point (Goja)
└── ui/
    └── index.html    # Frontend rendered in a sandboxed iframe
```

**Theme output:**
```
<name>/
├── package.json      # Theme metadata (id, name, version, author, css, overwrite)
├── theme.css         # CSS overrides (Aether :root design tokens pre-filled)
├── overwrite.json    # Optional asset overrides (sidebar-logo, background, etc.)
└── README.md
```

---

### `dev` — Live extension hot-reloading

```bash
# Watch and hot-sync the current extension directory to Aether
aether dev
# or
aet dev
```

Validates `manifest.json`, deploys the extension directly to your local Aether Launcher directory (`%APPDATA%\Aether\extensions\<id>` on Windows, `~/Library/Application Support/Aether/extensions/<id>` on macOS, `~/.local/share/Aether/extensions/<id>` on Linux), and continuously watches for file changes to keep the launcher in sync.

---

### `validate` — Validate an extension or theme

```bash
# Validate extension (reads manifest.json)
aether validate

# Validate theme (reads package.json)
aether validate --theme
```

**Auto-detection:** if `manifest.json` is absent but `package.json` is present, the project is automatically treated as a theme.

| Flag | Alias | Description |
|:---|:---|:---|
| `--theme` | `-t` | Force theme validation mode |

**Extension checks:** `manifest.json` exists, valid JSON, required fields: `id`, `name`, `version`, `main`, `api`.

**Theme checks:** `package.json` exists, valid JSON, required fields: `id`, `name`, `version`, CSS file referenced by `css` field actually exists on disk.

---

### `build` — Package into .aex or .theme

```bash
# Build extension → <id>-<version>.aex
aether build

# Build theme → <id>-<version>.theme
aether build --theme
```

**Auto-detection:** same logic as `validate` — detects project type automatically.

| Flag | Alias | Description |
|:---|:---|:---|
| `--theme` | `-t` | Force theme build mode |

Runs validation first, then packages all project files into a zip-format archive, automatically excluding `.git/`, `node_modules/`, and existing archive files.

---

### `help` — Show help

```bash
aether help
# or
aet help
```

---

## Quick Start — Extension

```bash
aet init my-extension com.example.myextension
cd my-extension
npm install   # installs @aethermc/sdk for TypeScript types
aet dev       # live test in local launcher
aet validate
aet build
# → com.example.myextension-1.0.0.aex
```

## Quick Start — Theme

```bash
aet init my-theme com.example.mytheme --theme
cd my-theme
# Edit theme.css to customise colours, radii, spacing...
aet validate
aet build
# → com.example.mytheme-1.0.0.theme
```

---

## Extension SDK

The scaffolded extension includes a dependency on [`@aethermc/sdk`](https://www.npmjs.com/package/@aethermc/sdk) — the official TypeScript SDK for Aether extensions. It provides full type definitions for the `Aether` global API and helper utilities:

```bash
npm install --save-dev @aethermc/sdk
```

See the [Aether SDK repository](https://github.com/Aether-Launcher/Aether-SDK) for full documentation.

---

## Publishing Your Extension

Once built, submit your `.aex` file alongside an entry in `index.json` to the [Aether Extension Registry](https://github.com/Aether-Launcher/Aether-Extensions) by opening a Pull Request. All extensions must comply with the [Aether Extension API License](https://github.com/Aether-Launcher/Aether-SDK/blob/main/LICENSE) — closed-source extensions are allowed, subject to Aether review.

---

## Development

The CLI is written in Go and uses only the standard library to keep the binary small and dependency-free.

```bash
# Build and install locally
go install ./...

# Run either command
aether help
aet help
```

---

## License

Licensed under the GNU General Public License v3.0 only. See [LICENSE](LICENSE).
