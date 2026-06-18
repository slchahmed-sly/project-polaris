package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/slchahmed-sly/project-polaris/internal/registry"
	"github.com/slchahmed-sly/project-polaris/internal/ui"
)

func main() {
	reg, err := registry.New()
	if err != nil {
		log.Fatalf("Failed to initialize registry: %v", err)
	}

	if err := reg.Load(); err != nil {
		log.Fatalf("Failed to load registry: %v", err)
	}

	// If no arguments are passed, launch the interactive UI
	if len(os.Args) < 2 {
		handleUI(reg)
		return
	}

	subcommand := os.Args[1]

	switch subcommand {
	case "add":
		handleAdd(reg)
	case "list":
		handleList(reg)
	case "set-cmd":
		handleSetCmd(reg)
	default:
		fmt.Printf("Unknown command: %s\n", subcommand)
		printUsage()
		os.Exit(1)
	}
}

func handleUI(reg *registry.Registry) {
	if len(reg.Projects) == 0 {
		fmt.Println("No projects registered. Use 'polaris add <path>' first.")
		return
	}

	// 1. Launch the UI and wait for a selection
	selectedPath, err := ui.RunMenu(reg)
	if err != nil {
		log.Fatalf("UI Error: %v", err)
	}

	// 2. If the user quit without selecting (e.g., pressed Esc)
	if selectedPath == "" {
		fmt.Println("Cancelled.")
		return
	}

	// 3. Spawn the IDE process
	openIDE(reg, selectedPath)
}

func openIDE(reg *registry.Registry, targetPath string) {
	cmdArgs := reg.Command
	if len(cmdArgs) == 0 {
		// Fallback default
		cmdArgs = []string{"agy", "."}
	}

	// Construct the command using the configured executable and arguments
	cmd := exec.Command(cmdArgs[0], cmdArgs[1:]...)

	// Set the working directory to the project the user selected
	cmd.Dir = targetPath

	// CRITICAL: Detach the I/O streams.
	// If we don't do this, the CLI will wait for the IDE to close,
	// or the IDE might die when the CLI exits.
	cmd.Stdin = nil
	cmd.Stdout = nil
	cmd.Stderr = nil

	// Start() runs the process in the background.
	// (Unlike Run(), which blocks until the process finishes).
	if err := cmd.Start(); err != nil {
		log.Fatalf("Failed to open IDE: %v", err)
	}

	fmt.Printf("Successfully opened %s\n", targetPath)
}

// ... keep your existing handleAdd, handleList, and printUsage functions here ...
func handleAdd(reg *registry.Registry) {
	if len(os.Args) < 3 {
		fmt.Println("Usage: polaris add <path>")
		os.Exit(1)
	}

	targetPath := os.Args[2]
	if err := reg.Add(targetPath); err != nil {
		log.Fatalf("Failed to add path: %v", err)
	}

	fmt.Printf("Successfully registered: %s\n", targetPath)
}

func handleSetCmd(reg *registry.Registry) {
	if len(os.Args) < 3 {
		fmt.Println("Usage: polaris set-cmd <command> [args...]")
		os.Exit(1)
	}

	commandArgs := os.Args[2:]
	if err := reg.SetCommand(commandArgs); err != nil {
		log.Fatalf("Failed to set command: %v", err)
	}

	fmt.Printf("Successfully set default command to: %v\n", commandArgs)
}

func handleList(reg *registry.Registry) {
	if len(reg.Projects) == 0 {
		fmt.Println("No projects registered yet. Use 'polaris add <path>' to add one.")
		return
	}

	fmt.Println("Registered Projects:")
	for i, p := range reg.Projects {
		fmt.Printf("  %d: %s\n", i+1, p)
	}
}

func printUsage() {
	fmt.Println("Project Navigator (polaris)")
	fmt.Println("Usage:")
	fmt.Println("  polaris add <path>   - Register a new project directory")
	fmt.Println("  polaris list         - List all registered projects")
	fmt.Println("  polaris set-cmd      - Set the command used to open projects (e.g. 'polaris set-cmd code .')")
}
