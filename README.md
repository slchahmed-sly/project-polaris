# Polaris

**The North Star for your projects.**

A fast, interactive, and beautiful CLI project navigator built in Go.

Polaris eliminates the friction of navigating through complex directory trees to open your code. Register your projects once, and let Polaris instantly guide you to them from anywhere in your terminal.

## Features

- **Interactive UI**: A beautiful, responsive terminal interface powered by Charm's Bubble Tea.
- **Fuzzy Search**: Instantly filter through dozens of projects just by typing.
- **Smart Display**: Automatically truncates system home paths (`~/`) for a clean, readable layout, while preserving absolute paths under the hood.
- **Instant Deletion**: Clean up your registry on the fly directly from the UI.
- **Detached Execution**: Spawns your IDE (`agy`) as a detached background process, keeping your terminal free and fast.
- **Zero Dependencies**: Compiles to a single, lightning-fast Go binary.

## Installation

Ensure you have Go installed on your system.

To build and install Polaris globally, run:

```bash
go install github.com/yourusername/polaris/cmd/polaris@latest
```

> **Note:** Ensure your `~/go/bin` directory is in your system's `$PATH` so you can run the command from anywhere.

## Usage

### 1. Register a Project

Navigate to any project directory you want Polaris to remember, and run:

```bash
polaris add .
```

You can also pass a specific path:

```bash
polaris add /path/to/project
```

### 2. Launch Polaris

From anywhere in your terminal, simply type:

```bash
polaris
```

### 3. Navigation & Controls

Once the Polaris menu is open:

| Key | Action |
|------|--------|
| ↑ / k | Move cursor up |
| ↓ / j | Move cursor down |
| Type | Start typing to instantly search/filter projects |
| Enter | Open the selected project in your IDE |
| x / Del | Remove the selected project from your registry |
| Esc / q | Quit without opening |

## Configuration

Polaris safely stores your registered paths in a standard JSON registry file locally on your machine.

Registry location:

**macOS / Linux**

```text
~/.config/polaris/projects.json
```

> **Note:** Polaris currently defaults to using `agy` as the command to open your IDE. Make sure this alias or command is available in your environment.

## Built With

- Go — The core language
- Bubble Tea — Terminal UI framework
- Bubbles — UI components (List)
- Lip Gloss — Terminal styling
