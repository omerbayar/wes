// Package tui implements wes's terminal user interface.
package tui

import (
	"path/filepath"
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"

	"github.com/omerbayar/wes/internal/core"
)

var (
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FAFAFA")).
			Background(lipgloss.Color("#7D56F4")).
			Padding(0, 1)
	pathStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("241"))
	dirStyle      = lipgloss.NewStyle().Foreground(lipgloss.Color("#7D56F4")).Bold(true)
	selectedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#FF79C6")).Bold(true)
	helpStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("241"))
	errStyle      = lipgloss.NewStyle().Foreground(lipgloss.Color("#FF5555"))
)

// headerLines and footerLines are the rows reserved around the list.
const (
	headerLines = 3 // title, path, blank
	footerLines = 2 // blank, help
)

type model struct {
	path    string
	entries []core.Entry
	err     error
	cursor  int
	offset  int // index of the first visible entry
	height  int // terminal height; 0 until the first WindowSizeMsg
}

func newModel(path string) model {
	abs, err := filepath.Abs(path)
	if err != nil {
		abs = path
	}
	entries, err := core.ListDir(abs)
	return model{path: abs, entries: entries, err: err}
}

func (m model) Init() tea.Cmd { return nil }

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.height = msg.Height
	case tea.KeyPressMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(m.entries)-1 {
				m.cursor++
			}
		}
	}
	m.clampOffset()
	return m, nil
}

// visibleRows is how many entries fit on screen.
func (m model) visibleRows() int {
	if m.height == 0 {
		return len(m.entries)
	}
	return max(1, m.height-headerLines-footerLines)
}

// clampOffset scrolls the window so the cursor stays visible.
func (m *model) clampOffset() {
	rows := m.visibleRows()
	if m.cursor < m.offset {
		m.offset = m.cursor
	}
	if m.cursor >= m.offset+rows {
		m.offset = m.cursor - rows + 1
	}
}

func (m model) View() tea.View {
	var b strings.Builder
	b.WriteString(titleStyle.Render("wes") + "\n")
	b.WriteString(pathStyle.Render(m.path) + "\n\n")

	switch {
	case m.err != nil:
		b.WriteString(errStyle.Render("error: "+m.err.Error()) + "\n")
	case len(m.entries) == 0:
		b.WriteString(helpStyle.Render("(empty directory)") + "\n")
	default:
		end := min(len(m.entries), m.offset+m.visibleRows())
		for i := m.offset; i < end; i++ {
			e := m.entries[i]
			name := e.Name
			if e.IsDir {
				name += "/"
			}
			switch {
			case i == m.cursor:
				b.WriteString(selectedStyle.Render("> "+name) + "\n")
			case e.IsDir:
				b.WriteString("  " + dirStyle.Render(name) + "\n")
			default:
				b.WriteString("  " + name + "\n")
			}
		}
	}

	b.WriteString("\n" + helpStyle.Render("↑/↓: move • q: quit"))

	v := tea.NewView(b.String())
	v.AltScreen = true
	return v
}

// Run starts the TUI showing the contents of dir.
func Run(dir string) error {
	_, err := tea.NewProgram(newModel(dir)).Run()
	return err
}
