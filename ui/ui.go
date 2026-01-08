package ui

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	lg "github.com/charmbracelet/lipgloss"
	"github.com/sftsrv/tri/picker"
	"github.com/sftsrv/tri/preview"
	"github.com/sftsrv/tri/theme"
	"github.com/sftsrv/tri/tree"
)

type window struct {
	width  int
	height int
}

type Path string

type Model struct {
	window window

	tree *tree.Tree

	hovered    *tree.Item
	selected   *tree.Item
	pathPicker picker.Model[*tree.Item]

	preview preview.Model
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (w *window) updateWindowSize(width int, height int) {
	w.width = width
	w.height = height
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.window.updateWindowSize(msg.Width, msg.Height)
		m.pathPicker = m.pathPicker.Height(msg.Height - 1).Width(int(float32(msg.Width)*0.25) - 1)
		m.preview, cmd = m.preview.Height(msg.Height - 1).Width(int(float32(msg.Width)*0.75) - 1).Update(msg)
		return m, cmd

	case picker.SelectedMsg[*tree.Item]:
		m.selected = msg.Selected
		return m, tea.Quit

	case picker.HoverMsg[*tree.Item]:
		m.hovered = msg.Hovered

		if msg.Hovered.IsFile() {
			m.preview = m.preview.SetPath(msg.Hovered.GetPath())
		}

		return m, nil

	case tea.KeyMsg:
		str := msg.String()

		switch str {
		case "ctrl+c":
			return m, tea.Quit

		case "left", "h":
			if m.hovered != nil && (str == "left" || !m.pathPicker.IsSearching()) {
				m.hovered.Collapse()
				m.pathPicker, cmd = m.pathPicker.Items(tree.ToItems(m.tree)).Update(msg)
				return m, cmd
			}

		case "right", "l":
			if m.hovered != nil && (str == "right" || !m.pathPicker.IsSearching()) {
				m.hovered.Expand()
				m.pathPicker, cmd = m.pathPicker.Items(tree.ToItems(m.tree)).Update(msg)
				return m, cmd
			}

		case "]":
			if !m.pathPicker.IsSearching() {
				m.tree.ExpandAll()
				m.pathPicker, cmd = m.pathPicker.Items(tree.ToItems(m.tree)).Update(msg)
				return m, cmd
			}

		case "[":
			if !m.pathPicker.IsSearching() {
				m.tree.CollapseAll()
				m.pathPicker, cmd = m.pathPicker.Items(tree.ToItems(m.tree)).Update(msg)
				return m, cmd
			}
		}

	case tea.MouseMsg:
		m.preview, cmd = m.preview.Update(msg)
		return m, cmd

	}

	m.pathPicker, cmd = m.pathPicker.Update(msg)
	return m, cmd
}

func helpView(m Model) string {
	item := func(icon string, title string) string {
		return lg.JoinHorizontal(
			lg.Top,
			theme.Secondary.MarginLeft(2).Render(icon),
			theme.Faded.MarginRight(2).Render(" ", title),
		)
	}

	help := ""
	if m.pathPicker.IsSearching() {
		help += item("esc", "close search")
		help += item("↓↑", "navigate")
		help += item("→", "expand")
		help += item("←", "collapse")
	} else {
		help += item("/", "search")
		help += item("↓↑/jk", "navigate")
		help += item("→/l", "expand")
		help += item("←/h", "collapse")
		help += item("]/[", "expand/collapse all")
	}

	return lg.NewStyle().Render(help)
}

func (m Model) View() string {
	return lg.JoinVertical(
		lg.Left,
		lg.JoinHorizontal(lg.Top, m.pathPicker.View(), m.preview.View()),
		helpView(m),
	)
}

func initialModel(f *tree.Tree, previewCmd string) Model {
	items := tree.ToItems(f)

	return Model{
		tree:       f,
		pathPicker: picker.New[*tree.Item]().Title("Items").Accent(theme.ColorPrimary).Items(items),
		preview:    preview.New(previewCmd),
	}
}

func Run(f *tree.Tree, previewCmd string) {
	m := initialModel(f, previewCmd)
	p := tea.NewProgram(m,
		tea.WithAltScreen(),
		tea.WithMouseCellMotion(),
	)

	result, err := p.Run()
	if err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}

	selected := result.(Model).selected

	if selected != nil {
		fmt.Println(result.(Model).selected.GetPath())
	}
}
