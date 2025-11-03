# vhs-mcp

A [Model Context Protocol (MCP)](https://modelcontextprotocol.io/) server that integrates [VHS](https://github.com/charmbracelet/vhs) terminal recording capabilities with Claude and other MCP-compatible clients.

## Overview

vhs-mcp enables AI assistants to programmatically create terminal session recordings and capture TUI (Text User Interface) application screenshots. By bridging VHS with the Model Context Protocol, it allows for automated terminal documentation and testing workflows.

## Features

- **Terminal Recording Automation**: Execute VHS tape files or tape text to generate terminal recordings
- **TUI Screenshot Capture**: Automatically capture screenshots of terminal applications when specific patterns appear
- **Multiple Output Formats**: Supports PNG, GIF, MP4, WebM, JSON, ASCII, and TXT outputs
- **Flexible Configuration**: Customizable timeouts, terminal dimensions, typing speeds, and more
- **MCP Integration**: Works seamlessly with Claude Desktop and other MCP-compatible clients

## Installation

### Using Nix (Recommended)

If you have Nix with flakes enabled:

```bash
# Run directly
nix run github:Warashi/vhs-mcp

# Build and install
nix build github:Warashi/vhs-mcp
```

### From Source

Requirements:
- Go 1.25.2 or later
- [VHS](https://github.com/charmbracelet/vhs) installed and available in PATH

```bash
# Clone the repository
git clone https://github.com/Warashi/vhs-mcp.git
cd vhs-mcp

# Build
go build -o vhs-mcp main.go

# Run
./vhs-mcp
```

## Tools

### vhs_run_tape

Runs a VHS tape and returns produced file URIs.

**Parameters:**

- `tape_text` (string, optional): VHS tape content as text
- `tape_path` (string, optional): Path to a .tape file
- `timeout_sec` (integer, optional): Timeout in seconds (default: 120)

**Note**: Exactly one of `tape_text` or `tape_path` must be provided.

**Example:**

```json
{
  "tape_text": "Output demo.gif\n\nSet FontSize 32\nSet Width 1200\nSet Height 600\n\nType \"echo 'Hello, World!'\" Enter\nSleep 500ms\n",
  "timeout_sec": 60
}
```

**Returns**: List of file URIs for generated artifacts (PNG, GIF, MP4, WebM, JSON, ASCII, TXT).

### tui_snap_after

Runs a command via VHS, waits for a regex pattern to appear on screen, then takes a PNG screenshot.

**Parameters:**

- `command` (string, required): Command to run (e.g., "htop")
- `wait_regex` (string, required): Regular expression to wait for on screen
- `keys` (array of strings, optional): Additional keystrokes to send (e.g., ["Down", "Enter"])
- `screenshot_name` (string, optional): Output PNG filename (default: "screenshot.png")
- `width` (integer, optional): Terminal width in pixels (default: 1000)
- `height` (integer, optional): Terminal height in pixels (default: 600)
- `typing_ms` (integer, optional): Typing speed in milliseconds (default: 0)
- `timeout_sec` (integer, optional): Timeout in seconds (default: 120)

**Supported Keys**: Up, Down, Left, Right, Enter, Tab, Space, Backspace, PageUp, PageDown, Escape (other strings will be typed as text)

**Example:**

```json
{
  "command": "htop",
  "wait_regex": "Load average",
  "keys": ["Down", "Down"],
  "screenshot_name": "htop-screenshot.png",
  "width": 1200,
  "height": 800,
  "timeout_sec": 30
}
```

**Returns**: File URI for the generated PNG screenshot.

## Development

### Using Nix

Enter the development shell with all required tools:

```bash
nix develop
```

This provides:
- Go toolchain
- Pre-commit hooks (actionlint, treefmt)
- Code formatters (gofmt, gofumpt, goimports, golines, nixfmt)
- Additional tools (keep-sorted, typos, pinact)

### Manual Setup

Requirements:
- Go 1.25.2+
- VHS

```bash
# Install dependencies
go mod download

# Run locally
go run main.go

# Build
go build -o vhs-mcp main.go

# Format code
go fmt ./...
```

## Project Structure

```
.
├── main.go              # Single-file MCP server implementation
├── go.mod               # Go module definition
├── go.sum               # Go dependencies
├── flake.nix            # Nix flake configuration
├── flake.lock           # Nix flake lock file
└── nixpkgs/             # Custom Nix packages
    ├── default.nix      # Package exports
    ├── vhs-mcp/         # Wrapped MCP server
    ├── vhs-mcp-unwrapped/ # Unwrapped Go binary
    ├── vhs/             # VHS package
    └── vhs-unwrapped/   # VHS unwrapped variant
```

## Dependencies

- [github.com/modelcontextprotocol/go-sdk](https://github.com/modelcontextprotocol/go-sdk) - Official MCP Go SDK
- [VHS](https://github.com/charmbracelet/vhs) - Terminal recording tool by Charm

## How It Works

1. The server implements the Model Context Protocol specification using the official Go SDK
2. Two tools are registered: `vhs_run_tape` and `tui_snap_after`
3. Each tool execution:
   - Creates a temporary directory for VHS output
   - Generates a VHS tape file based on the input parameters
   - Executes VHS with timeout protection
   - Collects generated artifacts and returns file URIs

## License

MIT

## Author

[Warashi](https://github.com/Warashi)

## See Also

- [Model Context Protocol](https://modelcontextprotocol.io/)
- [VHS - Write terminal GIFs as code](https://github.com/charmbracelet/vhs)
- [MCP Go SDK](https://github.com/modelcontextprotocol/go-sdk)
