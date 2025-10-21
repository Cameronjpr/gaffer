package tui

import (
	"fmt"
	"io"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type MenuModel struct {
	width  int
	height int
	list   list.Model
}

func NewMenuModel(items []list.Item) MenuModel {
	delegate := itemDelegate{}
	l := list.New(items, delegate, 0, 0)
	l.SetShowTitle(false)
	l.SetShowStatusBar(false)
	l.SetShowHelp(false)
	l.SetFilteringEnabled(false)

	return MenuModel{
		list: l,
	}
}

func (m MenuModel) Init() tea.Cmd {
	return nil
}

func (m MenuModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.list.SetSize(msg.Width, msg.Height/2)
		return m, nil

	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			// Check which item is selected
			selected := m.list.SelectedItem()
			if i, ok := selected.(item); ok {
				choice := string(i)
				switch choice {
				case "New game":
					return m, func() tea.Msg { return startPreMatchMsg{} }
				case "Settings":
					// TODO: return settings mode message
					return m, nil
				}
			}
			return m, nil
		}
	}

	// Delegate to list for navigation (up/down/etc)
	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m MenuModel) View() string {
	title := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("170")).
		Render("G A F F E R")

	help := lipgloss.NewStyle().
		Foreground(lipgloss.Color("241")).
		Render("\n↑/↓: navigate • enter: select • esc: quit")

	content := lipgloss.JoinVertical(
		lipgloss.Center,
		"\n"+title,
		"\n",
		m.list.View(),
		help,
	)

	return lipgloss.Place(
		m.width,
		m.height,
		lipgloss.Center,
		lipgloss.Center,
		content,
	)
}

// item represents a menu item
type item string

func (i item) FilterValue() string { return "" }

// itemDelegate handles rendering of menu items
type itemDelegate struct{}

func (d itemDelegate) Height() int                             { return 1 }
func (d itemDelegate) Spacing() int                            { return 0 }
func (d itemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d itemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(item)
	if !ok {
		return
	}

	str := string(i)

	// Style for unselected items
	style := lipgloss.NewStyle().
		PaddingLeft(2).
		Foreground(lipgloss.Color("252"))

	// Style for selected item
	if index == m.Index() {
		style = lipgloss.NewStyle().
			PaddingLeft(2).
			Foreground(lipgloss.Color("170")).
			Bold(true)
		str = "> " + str
	}

	fmt.Fprint(w, style.Render(str))
}
