package tui

import (
	"github.com/cameronjpr/gaffer/internal/components"
	"github.com/cameronjpr/gaffer/internal/domain"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type ManagerHubModel struct {
	ChosenClub  *domain.Club
	Fixtures    []*domain.Fixture
	LeagueTable *domain.LeagueTable
	width       int
	height      int
}

func NewManagerHubModel(club *domain.Club, fixtures []*domain.Fixture, leagueTable *domain.LeagueTable) *ManagerHubModel {
	return &ManagerHubModel{
		ChosenClub:  club,
		Fixtures:    fixtures,
		LeagueTable: leagueTable,
	}
}
func (m *ManagerHubModel) Init() tea.Cmd {
	return nil
}

func (m *ManagerHubModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeySpace:
			return m, func() tea.Msg {
				return startPreMatchMsg{}
			}
		}
	}

	return m, nil
}

func (m *ManagerHubModel) View() string {
	if m.ChosenClub == nil {
		panic("ChosenClub is nil, cannot proceed")
	}
	header := lipgloss.NewStyle().
		Align(lipgloss.Center).
		Width(m.width).
		Padding(1, 2).
		Bold(true).
		Background(lipgloss.Color(m.ChosenClub.Background)).
		Foreground(lipgloss.Color(m.ChosenClub.Foreground)).
		Render(m.ChosenClub.Name)
	footer := lipgloss.NewStyle().
		Align(lipgloss.Center).
		Width(m.width).
		Render("Press [Space] to start pre-match")

	leagueTableView := ""
	if m.LeagueTable != nil {
		leagueTableView = components.Table(*m.LeagueTable)
	}
	fixturesView := components.Fixtures(m.Fixtures)

	content := components.ThreeColumnLayout(
		m.width,
		fixturesView,
		"",
		leagueTableView,
	)

	main := lipgloss.NewStyle().
		Width(m.width).
		Height(m.height-lipgloss.Height(header)-lipgloss.Height(footer)).
		Align(lipgloss.Center, lipgloss.Center).
		Render(content)

	return lipgloss.JoinVertical(lipgloss.Top, header, main, footer)

}
