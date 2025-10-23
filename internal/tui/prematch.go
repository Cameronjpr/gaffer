package tui

import (
	"fmt"

	"github.com/cameronjpr/gaffer/internal/domain"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type PreMatchModel struct {
	match  *domain.Match
	width  int
	height int
}

func NewPreMatchModel(match domain.Match) PreMatchModel {
	return PreMatchModel{
		match: &match,
	}
}

func (m PreMatchModel) Init() tea.Cmd {
	return nil
}

func (m PreMatchModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC:
			return m, tea.Quit
		case tea.KeyEnter:
			return m, func() tea.Msg { return startMatchMsg{} }
		}
	}

	return m, nil
}

func (m PreMatchModel) View() string {
	// Calculate column width
	colWidth := m.width / 3

	// Home team section
	homeContent := lipgloss.JoinVertical(
		lipgloss.Left,
		lipgloss.NewStyle().Bold(true).Render(m.match.Home.Club.Name),
		lipgloss.NewStyle().Italic(true).Render(m.match.Home.Formation),
		"",
		m.match.Home.GetLineup(nil),
	)

	// Match info section
	matchContent := fmt.Sprintf("%s vs %s\n\nPress Enter to start", m.match.Home.Club.Name, m.match.Away.Club.Name)

	// Away team section
	awayContent := lipgloss.JoinVertical(
		lipgloss.Left,
		lipgloss.NewStyle().Bold(true).Render(m.match.Away.Club.Name),
		lipgloss.NewStyle().Italic(true).Render(m.match.Away.Formation),
		"",
		m.match.Away.GetLineup(nil),
	)

	// Create 3-column layout
	layout := lipgloss.JoinHorizontal(
		lipgloss.Top,
		lipgloss.Place(colWidth, m.height, lipgloss.Center, lipgloss.Center, homeContent),
		lipgloss.Place(colWidth, m.height, lipgloss.Center, lipgloss.Center, matchContent),
		lipgloss.Place(colWidth, m.height, lipgloss.Center, lipgloss.Center, awayContent),
	)

	return layout
}
