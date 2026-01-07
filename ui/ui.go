package ui

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/sftsrv/tri/picker"
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

	hovered    tree.Item
	path       string
	pathPicker picker.Model[tree.Item]
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (w *window) updateWindowSize(width int, height int) {
	w.width = width
	w.height = height
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.window.updateWindowSize(msg.Width, msg.Height)
		m.pathPicker = m.pathPicker.Height(msg.Height - 1)
		return m, nil

	case picker.SelectedMsg[tree.Item]:
		// handle selection

	case picker.HoverMsg[tree.Item]:
		m.hovered = msg.Hovered
		return m, nil

	case tea.KeyMsg:
		str := msg.String()

		switch str {
		case "ctrl+c":
			return m, tea.Quit

		case "enter":
			// enter pressed, exit with selection

		case "left":
			m.hovered.Collapse()
			m.pathPicker = m.pathPicker.Items(tree.ToItems(m.tree))
			return m, nil

		case "right":
			m.hovered.Expand()
			m.pathPicker = m.pathPicker.Items(tree.ToItems(m.tree))
			return m, nil
		}

		var cmd tea.Cmd
		m.pathPicker, cmd = m.pathPicker.Update(msg)

		return m, cmd

	}

	return m, nil
}

func (m Model) View() string {
	return m.pathPicker.View()
}

func initialModel(f *tree.Tree) Model {
	items := tree.ToItems(f)

	return Model{
		tree:       f,
		pathPicker: picker.New[tree.Item]().Title("Items").Accent(theme.ColorPrimary).Items(items).Searching(true),
	}
}

func Run(f *tree.Tree) {
	m := initialModel(f)
	p := tea.NewProgram(m)

	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
