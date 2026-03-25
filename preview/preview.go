package preview

import (
	"os/exec"
	"strconv"
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	lg "github.com/charmbracelet/lipgloss"
	"github.com/sftsrv/tri/theme"
)

type Model struct {
	cmd      string
	path     string
	ready    bool
	width    int
	height   int
	active   *exec.Cmd
	viewport viewport.Model
}

func New(cmd string) Model {
	return Model{
		cmd: cmd,
	}
}

func (m Model) SetPath(path string) (Model, tea.Cmd) {
	m.path = path

	if m.active != nil && m.active.Process != nil {
		m.active.Process.Kill()
	}

	m.viewport.SetContent("Loading Path: " + path)
	active, cmd := preview(m.cmd, path, m.width)

	m.active = active
	return m, cmd
}

func (m Model) SetContent(preview PreviewResultMsg) Model {
	m.viewport.SetContent(preview.content)
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

func resolveCommand(command string, width int) (string, []string) {
	command = strings.TrimSpace(command)
	if len(command) > 0 {
		parts := strings.Split(command, " ")
		bin := parts[0]
		args := parts[1:]

		return bin, args
	}

	_, err := exec.LookPath("bat")
	if err == nil {
		return "bat", []string{"--color=always", "--number", "--terminal-width", strconv.Itoa(width)}
	} else {
		return "cat", []string{}
	}
}

type PreviewResultMsg struct {
	content string
}

func preview(command string, path string, width int) (*exec.Cmd, tea.Cmd) {
	if path == "" {
		return nil, func() tea.Msg {
			return PreviewResultMsg{"No file selected"}
		}
	}

	bin, args := resolveCommand(command, width)
	cmd := exec.Command(bin, append(args, path)...)

	return cmd, func() tea.Msg {
		output, err := cmd.CombinedOutput()
		if err != nil {
			return PreviewResultMsg{"ERROR reading " + path + ":" + string(err.Error())}
		}

		return PreviewResultMsg{string(output)}
	}
}

func (m Model) View() string {
	return lg.JoinVertical(
		lg.Center,
		lg.NewStyle().
			Width(m.width).
			PaddingLeft(1).
			PaddingRight(1).
			Background(theme.ColorSecondary).
			Render(m.path),
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
