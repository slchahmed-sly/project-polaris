package registry

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
)

// Registry holds the state of our saved project paths.
type Registry struct {
	// Projects is exported (capitalized) so the JSON package can see it,
	// but the slice itself is managed internally.
	Projects []string `json:"projects"`

	// Command is the user-configured command to open a project.
	Command []string `json:"command,omitempty"`

	// filePath is unexported (lowercase), encapsulating the storage detail.
	filePath string
}

// New initializes a Registry, resolving the correct cross-platform config path.
func New() (*Registry, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return nil, err
	}

	// Create a dedicated directory: e.g., ~/.config/projnav
	appDir := filepath.Join(configDir, "polaris")

	// os.MkdirAll is like `mkdir -p`. 0755 provides standard read/execute permissions.
	if err := os.MkdirAll(appDir, 0755); err != nil {
		return nil, err
	}

	return &Registry{
		Projects: make([]string, 0),
		filePath: filepath.Join(appDir, "projects.json"),
	}, nil
}

// Load reads the JSON file from disk into memory.
func (r *Registry) Load() error {
	data, err := os.ReadFile(r.filePath)
	if err != nil {
		// If the file doesn't exist yet, that's fine. We just return nil (no error).
		if errors.Is(err, os.ErrNotExist) {
			return nil
		}
		return err
	}

	// Deserialize the JSON directly into our struct.
	return json.Unmarshal(data, r)
}

// Save writes the current state back to disk.
func (r *Registry) Save() error {
	data, err := json.MarshalIndent(r, "", "  ")
	if err != nil {
		return err
	}

	// 0644 gives read/write to the owner, read-only to everyone else.
	return os.WriteFile(r.filePath, data, 0644)
}

// Add safely appends a new path and persists it, preventing duplicates.
func (r *Registry) Add(path string) error {
	// Always store absolute paths so execution works regardless of where the CLI is run.
	absPath, err := filepath.Abs(path)
	if err != nil {
		return err
	}

	// Prevent duplicates
	for _, p := range r.Projects {
		if p == absPath {
			return nil
		}
	}

	r.Projects = append(r.Projects, absPath)
	return r.Save()
}

// Remove deletes a path from the registry and saves the updated list.
func (r *Registry) Remove(path string) error {
	for i, p := range r.Projects {
		if p == path {
			// Remove the item from the slice while maintaining order
			r.Projects = append(r.Projects[:i], r.Projects[i+1:]...)
			return r.Save()
		}
	}
	return nil // Path wasn't in the list anyway
}

// SetCommand updates the command used to open projects and saves it.
func (r *Registry) SetCommand(cmd []string) error {
	r.Command = cmd
	return r.Save()
}
