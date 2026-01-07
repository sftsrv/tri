package preview

import (
	"os/exec"
	"strconv"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	lg "github.com/charmbracelet/lipgloss"
	"github.com/sftsrv/tri/theme"
)

type Model struct {
	path     string
	ready    bool
	width    int
	height   int
	viewport viewport.Model
}

func New() Model {
	return Model{}
}

func (m Model) SetPath(path string) Model {
	m.path = path

	content := preview(path, m.width)
	m.viewport.SetContent(content)
	return m
}

func (m Model) ClearPath() Model {
	m.path = ""
	m.viewport.SetContent("")
	return m
}

func (m Model) Width(width int) Model {
	m.width = width
	return m
}

func (m Model) Height(height int) Model {
	m.height = height
	return m
}

func preview(path string, width int) string {
	_, err := exec.LookPath("bat")

	if path == "" {
		return "No file selected"
	}

	var cmd *exec.Cmd
	if err == nil {
		cmd = exec.Command("bat", "--color=always", "--number", "--terminal-width", strconv.Itoa(width), path)
	} else {
		cmd = exec.Command("cat", path)
	}

	output, err := cmd.CombinedOutput()
	if err != nil {
		return "ERROR: " + string(output)
	}

	return string(output)
}

func (m Model) View() string {
	return lg.JoinVertical(
		lg.Center,
		lg.NewStyle().Width(m.width).PaddingLeft(1).PaddingRight(1).Background(theme.ColorSecondary).Render(m.path),
		m.viewport.View(),
	)
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg.(type) {

	case tea.WindowSizeMsg:
		if !m.ready {
			// Since this program is using the full size of the viewport we
			// need to wait until we've received the window dimensions before
			// we can initialize the viewport. The initial dimensions come in
			// quickly, though asynchronously, which is why we wait for them
			// here.
			m.viewport = viewport.New(m.width, m.height-1)
			m.viewport.YPosition = 0
			m.viewport.SetContent("")
			m.ready = true
		} else {
			m.viewport.Width = m.width
			m.viewport.Height = m.height
		}
	}

	// Handle keyboard and mouse events in the viewport
	m.viewport, cmd = m.viewport.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}
