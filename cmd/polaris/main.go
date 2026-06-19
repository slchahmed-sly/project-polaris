// Package main is the entry point for the Polaris CLI application.
// Polaris is a fast, interactive CLI project navigator that helps you
// register project paths and quickly open them in your preferred IDE.
package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/charmbracelet/lipgloss"
	"github.com/slchahmed-sly/project-polaris/internal/registry"
	"github.com/slchahmed-sly/project-polaris/internal/ui"
)

// main is the entry point of the Polaris application. It initializes the
// registry, loads saved projects, and routes the command line arguments
// to the appropriate handler.
func main() {
	reg, err := registry.New()
	if err != nil {
		log.Fatalf("Failed to initialize registry: %v", err)
	}

	if err := reg.Load(); err != nil {
		log.Fatalf("Failed to load registry: %v", err)
	}

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
	case "help", "-h", "--help":
		handleHelp()
	default:
		fmt.Printf("Unknown command: %s\n", subcommand)
		printUsage()
		os.Exit(1)
	}
}

// handleUI launches the interactive terminal user interface using Bubble Tea.
// It allows the user to fuzzy search and select a project to open.
func handleUI(reg *registry.Registry) {
	if len(reg.Projects) == 0 {
		fmt.Println("No projects registered. Use 'polaris add <path>' first.")
		return
	}

	selectedPath, err := ui.RunMenu(reg)
	if err != nil {
		log.Fatalf("UI Error: %v", err)
	}

	if selectedPath == "" {
		fmt.Println("Cancelled.")
		return
	}

	openIDE(reg, selectedPath)
}

// openIDE spawns the user's configured IDE as a detached background process
// targeting the selected project path. It also bumps the project to the top
// of the recently used list in the registry.
func openIDE(reg *registry.Registry, targetPath string) {
	cmdArgs := reg.Command
	if len(cmdArgs) == 0 {
		cmdArgs = []string{"agy", "."}
	}

	cmd := exec.Command(cmdArgs[0], cmdArgs[1:]...)

	cmd.Dir = targetPath

	// CRITICAL: Detach the I/O streams.
	// If we don't do this, the CLI will wait for the IDE to close,
	// or the IDE might die when the CLI exits.
	cmd.Stdin = nil
	cmd.Stdout = nil
	cmd.Stderr = nil

	if err := cmd.Start(); err != nil {
		log.Fatalf("Failed to open IDE: %v", err)
	}

	reg.Bump(targetPath)

	fmt.Printf("Successfully opened %s\n", targetPath)
}

// handleAdd registers a new project path into the Polaris registry.
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

// handleSetCmd configures the default command used to open projects
// (e.g., "code .", "cursor .", "agy .").
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

// handleList prints out all currently registered project paths to the terminal.
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

// printUsage prints a brief usage summary of the available Polaris commands.
func printUsage() {
	fmt.Println("Project Navigator (polaris)")
	fmt.Println("Usage:")
	fmt.Println("  polaris add <path>   - Register a new project directory")
	fmt.Println("  polaris list         - List all registered projects")
	fmt.Println("  polaris set-cmd      - Set the command used to open projects (e.g. 'polaris set-cmd code .')")
	fmt.Println("  polaris help         - Show detailed help information")
}

// handleHelp prints detailed help information, including how to use the
// CLI commands and how to navigate the interactive UI.
func handleHelp() {
	titleStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("39"))
	subtitleStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("241"))
	headerStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("43"))
	cmdStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("78"))

	fmt.Printf("%s - %s\n\n", titleStyle.Render("Polaris"), subtitleStyle.Render("The North Star for your projects."))

	fmt.Println("Polaris is a fast, interactive CLI project navigator. It helps you save paths to your")
	fmt.Println("frequently used projects and open them in your favorite IDE instantly from anywhere.")

	fmt.Println(headerStyle.Render("How it works:"))
	fmt.Println("1. Navigate to a project you want to save.")
	fmt.Printf("2. Run '%s' to register it.\n", cmdStyle.Render("polaris add ."))
	fmt.Printf("3. Run '%s' from anywhere to launch the interactive UI.\n", cmdStyle.Render("polaris"))
	fmt.Println("4. Select the project, and Polaris will open it in your IDE.")

	fmt.Println(headerStyle.Render("Commands:"))
	fmt.Printf("  %-25s Register a new project directory (e.g., 'polaris add .')\n", cmdStyle.Render("add <path>"))
	fmt.Printf("  %-25s List all registered projects\n", cmdStyle.Render("list"))
	fmt.Printf("  %-25s Set the command used to open projects (e.g., 'polaris set-cmd code .')\n", cmdStyle.Render("set-cmd"))
	fmt.Printf("  %-25s Show this detailed help message\n\n", cmdStyle.Render("help, -h, --help"))

	fmt.Println(headerStyle.Render("In the interactive UI:"))
	fmt.Println("  Type to fuzzy search your projects.")
	fmt.Println("  Press 'Enter' to open the selected project.")
	fmt.Println("  Press 'x' or 'Delete' to remove a project from your registry.")
	fmt.Println("  Press 'Esc' or 'q' to quit without opening.")
}
