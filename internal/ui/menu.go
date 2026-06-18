package ui

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/slchahmed-sly/project-polaris/internal/registry"
)

var docStyle = lipgloss.NewStyle().Margin(1, 2)

// item represents a single project in our list
type item struct {
	fullPath    string
	displayPath string
}

// These three methods implement the list.Item interface required by bubbles/list
func (i item) Title() string       { return filepath.Base(i.fullPath) }
func (i item) Description() string { return i.displayPath }
func (i item) FilterValue() string { return i.fullPath }

type model struct {
	list     list.Model
	registry *registry.Registry
	choice   string
}

func RunMenu(reg *registry.Registry) (string, error) {
	items := make([]list.Item, 0, len(reg.Projects))
	homeDir, _ := os.UserHomeDir()

	for _, p := range reg.Projects {
		display := p
		if homeDir != "" && strings.HasPrefix(p, homeDir) {
			display = strings.Replace(p, homeDir, "~", 1)
		}
		items = append(items, item{fullPath: p, displayPath: display})
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

func (m model) Init() tea.Cmd {
	return nil
}

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
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m model) View() string {
	return docStyle.Render(m.list.View())
}
