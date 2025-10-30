package tui

import (
	"github.com/cameronjpr/gaffer/internal/domain"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// TacticsModel handles the tactics tab UI
type TacticsModel struct {
	width    int
	height   int
	userTeam *domain.MatchParticipant
}

// NewTacticsModel creates a new tactics model
func NewTacticsModel(userTeam *domain.MatchParticipant) *TacticsModel {
	return &TacticsModel{
		userTeam: userTeam,
	}
}

func (m *TacticsModel) Init() tea.Cmd {
	return nil
}

func (m *TacticsModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil
	}

	return m, nil
}

func (m *TacticsModel) View() string {
	// Placeholder content
	placeholderStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("241")).
		Bold(true).
		Align(lipgloss.Center)

	content := placeholderStyle.Render("Tactics options coming soon...")

	// Center the placeholder
	centered := lipgloss.NewStyle().
		Width(m.width).
		Height(m.height).
		Align(lipgloss.Center, lipgloss.Center).
		Render(content)

	return centered
}
