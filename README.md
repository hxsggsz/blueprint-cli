# Blueprint CLI

A fast, concurrent CLI tool written in Go for scaffolding projects from templates. Blueprint copies the directory structure of a chosen template into your current working directory using a concurrent worker pool for maximum performance.

## What is it?

Blueprint is a lightweight, self-hosted project generator. You maintain your own templates (like `react-app`, `go-service`, `rust-cli`, etc.) and Blueprint provides an interactive selection to copy the desired template into the directory you're working in.

## Features

- **Interactive template selection** — Lists all available templates and lets you pick by number
- **Concurrent copying** — Worker pool of goroutines (`NumCPU * 2`) for parallel file copying
- **Parallel directory traversal** — Recursively walks subdirectories in separate goroutines
- **File streaming** — Uses `io.Copy` for direct streaming, no need to load entire files into memory
- **Permission preservation** — Keeps original file modes/permissions
- **Automatic ignoring** — Skips `.DS_Store`, `Thumbs.db`, `.git`, `node_modules`, `dist`, `build` by default
- **Custom ignores** — Add your own ignore patterns via configuration
- **Persistent configuration** — JSON config file at `~/.config/blueprint/config.json`

## Prerequisites

- Go 1.27.0 or higher

## Installation



https://github.com/user-attachments/assets/e35a1c8f-49ad-4a0b-a5ab-dd325b38bb0b





### One-liner (Linux / macOS)

```bash
curl -fsSL https://raw.githubusercontent.com/hxsggsz/blueprint-cli/main/install.sh | bash
```

The installer downloads the latest release binary into `~/.local/bin`, then asks for the local directory path where your templates are stored and writes it as `template_path` in `~/.config/blueprint/config.json`.

You can skip or customize the prompts with environment variables:

| Variable | Description |
|---|---|
| `BLUEPRINT_TEMPLATES_PATH` | Templates directory path (skips the prompt when set) |
| `BLUEPRINT_SKIP_TEMPLATES=1` | Skip template setup entirely |
| `BLUEPRINT_INSTALL_DIR` | Override the binary install directory (default `~/.local/bin`) |
| `BLUEPRINT_CONFIG_DIR` | Override the config directory (default `~/.config/blueprint`) |
| `BLUEPRINT_VERSION` | Install a specific version instead of the latest release |

### From source

```bash
# Build
make build
# or
go build -o blueprint .

# Or run directly
make dev
```

## Usage


https://github.com/user-attachments/assets/a59eba3d-1ce8-4659-81e8-c922c76aeeaa



1. On first run, the config file is automatically created at `~/.config/blueprint/config.json`. Edit it and set `template_path` to an **absolute path** (this path is resolved independently of your current working directory, so it must be absolute):

```json
{
  "template_path": "/path/to/your/templates",
  "ignore_file_paths": ["extra_folder", "some_file.txt"]
}
```

> **Note:** `template_path` must be an absolute path. The config file lives at `~/.config/blueprint/config.json`, so any relative path would be resolved against that location rather than where you run `blueprint`.

2. The `template_path` directory should contain subdirectories, each representing a template:

```
/path/to/your/templates/
├── react-app/
│   ├── package.json
│   ├── src/
│   └── ...
├── go-service/
│   ├── main.go
│   ├── go.mod
│   └── ...
└── rust-cli/
    └── ...
```

3. Run Blueprint from the directory where the project should be scaffolded:

```bash
cd ~/my-new-project
blueprint
```

4. Select the desired template by its number.

## Configuration

| Field | Type | Description |
|---|---|---|
| `template_path` | `string` (required) | Absolute path to the directory containing templates. Must be absolute, since it is resolved relative to the config file (`~/.config/blueprint/config.json`), not your current working directory. |
| `ignore_file_paths` | `[]string` (optional) | Extra file/directory names to skip during copying |

**Default ignores** (always applied):
`.DS_Store`, `Thumbs.db`, `.git`, `node_modules`, `dist`, `build`

## Project Structure

```
blueprint-cli/
├── main.go                    # Entry point: config, template selection, worker pool
├── config/
│   └── config.go              # Configuration management (JSON-based)
└── pkg/
    └── template-copier.go     # Copy engine (traversal, streaming, workers)
```

## How it works

1. **Config** — Reads and validates the JSON configuration file
2. **ListTemplates** — Reads the template directory and returns subdirectory names
3. **selectTemplate** — Interactive menu for the user to pick a template
4. **Start** — Kicks off recursive directory traversal in goroutines
5. **Worker Pool** — `NumCPU * 2` workers consume jobs from a buffered channel (capacity 100)
6. **CopyJob.CreateFile** — Opens source file, creates destination directory if needed, and streams via `io.Copy`

## Tech Stack

- **Language:** Go 1.27.0
- **External dependencies:** None (standard library only)
- **Concurrency:** goroutines + channels + `sync.WaitGroup`
- **I/O:** Streaming via `io.Copy`
- **Configuration:** JSON

## License

MIT
