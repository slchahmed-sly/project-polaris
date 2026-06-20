// Package ui provides the interactive terminal user interface for Polaris
// using the Charm Bubble Tea framework.
package ui

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/slchahmed-sly/project-polaris/internal/registry"
)

var (
	docStyle       = lipgloss.NewStyle().Margin(1, 2)
	gitCleanStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("42"))
	gitDirtyStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("214"))
	gitRemoteStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("39"))
	gitNoneStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
)

// item represents a single project in our list
type item struct {
	fullPath    string
	displayPath string
	gitFetched  bool
	isGitRepo   bool
	isDirty     bool
	hasRemote   bool
	maxWidth    int
}

// These three methods implement the list.Item interface required by bubbles/list
func (i item) Title() string { return filepath.Base(i.fullPath) }
func (i item) Description() string {
	paddedPath := fmt.Sprintf("%-*s", i.maxWidth, i.displayPath)

	if !i.gitFetched {
		return paddedPath + "   " + gitNoneStyle.Render("(checking...)")
	}
	if !i.isGitRepo {
		return paddedPath
	}

	var tags []string

	if i.isDirty {
		tags = append(tags, gitDirtyStyle.Render("✗"))
	} else {
		tags = append(tags, gitCleanStyle.Render("✓"))
	}

	if i.hasRemote {
		tags = append(tags, gitRemoteStyle.Render("Remote"))
	}

	return fmt.Sprintf("%s    %s", paddedPath, strings.Join(tags, "  "))
}
func (i item) FilterValue() string { return i.fullPath }

type gitStatusMsg struct {
	index     int
	isGitRepo bool
	isDirty   bool
	hasRemote bool
}

func checkGitStatus(index int, path string) tea.Cmd {
	return func() tea.Msg {
		msg := gitStatusMsg{index: index}

		// Check if it's a git repo
		cmd := exec.Command("git", "rev-parse", "--is-inside-work-tree")
		cmd.Dir = path
		if err := cmd.Run(); err != nil {
			return msg // not a git repo
		}
		msg.isGitRepo = true

		// Check if dirty
		cmd = exec.Command("git", "status", "--porcelain")
		cmd.Dir = path
		if out, err := cmd.Output(); err == nil {
			if len(strings.TrimSpace(string(out))) > 0 {
				msg.isDirty = true
			}
		}

		// Check if remote exists
		cmd = exec.Command("git", "remote", "-v")
		cmd.Dir = path
		if out, err := cmd.Output(); err == nil {
			if len(strings.TrimSpace(string(out))) > 0 {
				msg.hasRemote = true
			}
		}

		return msg
	}
}

// model represents the application state for the Bubble Tea UI.
type model struct {
	list     list.Model
	registry *registry.Registry
	choice   string
}

// RunMenu initializes and starts the Bubble Tea program. It returns the
// path of the selected project or an empty string if the user cancels.
func RunMenu(reg *registry.Registry) (string, error) {
	items := make([]list.Item, 0, len(reg.Projects))
	homeDir, _ := os.UserHomeDir()

	var maxPathLen int
	for _, p := range reg.Projects {
		display := p
		if homeDir != "" && strings.HasPrefix(p, homeDir) {
			display = strings.Replace(p, homeDir, "~", 1)
		}
		if len(display) > maxPathLen {
			maxPathLen = len(display)
		}
	}

	for _, p := range reg.Projects {
		display := p
		if homeDir != "" && strings.HasPrefix(p, homeDir) {
			display = strings.Replace(p, homeDir, "~", 1)
		}
		items = append(items, item{
			fullPath:    p,
			displayPath: display,
			maxWidth:    maxPathLen,
		})
	}

	m := model{registry: reg}

	delegate := list.NewDefaultDelegate()
	blue := lipgloss.Color("39")
	delegate.Styles.SelectedTitle = delegate.Styles.SelectedTitle.Foreground(blue).BorderLeftForeground(blue)
	delegate.Styles.SelectedDesc = delegate.Styles.SelectedDesc.Foreground(blue).BorderLeftForeground(blue)

	m.list = list.New(items, delegate, 0, 0)
	m.list.Title = "Polaris - Select a Project"

	m.list.AdditionalShortHelpKeys = func() []key.Binding {
		return []key.Binding{
			key.NewBinding(key.WithKeys("x", "delete"), key.WithHelp("x/del", "remove project")),
		}
	}

	p := tea.NewProgram(m, tea.WithAltScreen())
	finalModel, err := p.Run()
	if err != nil {
		return "", err
	}

	if m, ok := finalModel.(model); ok {
		return m.choice, nil
	}
	return "", nil
}

// Init is called when the Bubble Tea program starts.
func (m model) Init() tea.Cmd {
	var cmds []tea.Cmd
	for i, itm := range m.list.Items() {
		if p, ok := itm.(item); ok {
			cmds = append(cmds, checkGitStatus(i, p.fullPath))
		}
	}
	return tea.Batch(cmds...)
}

// Update handles incoming events (like key presses or window resizes) and
// updates the model accordingly.
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.KeyMsg:
		if m.list.FilterState() == list.Filtering {
			break
		}

		switch msg.String() {
		case "ctrl+c", "q", "esc":
			return m, tea.Quit

		case "enter":
			if selectedItem, ok := m.list.SelectedItem().(item); ok {
				m.choice = selectedItem.fullPath
				return m, tea.Quit
			}

		case "x", "delete":
			if selectedItem, ok := m.list.SelectedItem().(item); ok {
				m.registry.Remove(selectedItem.fullPath)

				index := m.list.Index()
				m.list.RemoveItem(index)
			}
		}

	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)

	case gitStatusMsg:
		items := m.list.Items()
		if msg.index >= 0 && msg.index < len(items) {
			if itm, ok := items[msg.index].(item); ok {
				itm.gitFetched = true
				itm.isGitRepo = msg.isGitRepo
				itm.isDirty = msg.isDirty
				itm.hasRemote = msg.hasRemote
				cmd := m.list.SetItem(msg.index, itm)
				return m, cmd
			}
		}
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

// View renders the UI based on the current model state.
func (m model) View() string {
	return docStyle.Render(m.list.View())
}
