<img width="2176" height="996" alt="polaris-log" src="https://github.com/user-attachments/assets/4bb82f4e-5b75-45b4-b115-aa18d24c79dc" />

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

### Changing the IDE Command
Polaris defaults to using `agy .` as the command to open your IDE. If `agy` is not available on your system, or if you prefer to use a different IDE, you can easily configure the command using `set-cmd`:

```bash
# Configure Polaris to use VS Code
polaris set-cmd code .

# Configure Polaris to use Cursor
polaris set-cmd cursor .
```
This configuration is natively executed using Go's `exec.Command`, which means it works seamlessly across both **macOS/Linux** (resolving standard binaries in `$PATH`) and **Windows** (automatically resolving `.cmd`, `.bat`, or `.exe` files in `%PATH%`).

### Registry Location
Polaris safely stores your registered paths and configuration in a standard JSON registry file locally on your machine. The location differs depending on your OS:

**macOS**
```text
~/Library/Application Support/polaris/projects.json
```

**Linux**
```text
~/.config/polaris/projects.json
```

**Windows**
```text
%AppData%\polaris\projects.json
```

## Contributing

Contributions are welcome and greatly appreciated! Please read our [Contributing Guidelines](CONTRIBUTING.md) for details on how to get started, and our [Code of Conduct](CODE_OF_CONDUCT.md) to ensure a welcoming environment for everyone.

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.

## Built With

- Go — The core language
- Bubble Tea — Terminal UI framework
- Bubbles — UI components (List)
- Lip Gloss — Terminal styling
