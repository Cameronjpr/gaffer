package tui

import (
	"github.com/cameronjpr/gaffer/internal/components"
	"github.com/cameronjpr/gaffer/internal/domain"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type ManagerHubModel struct {
	Season      *domain.Season
	ChosenClub  *domain.Club
	Fixtures    []*domain.Fixture
	LeagueTable *domain.LeagueTable
	width       int
	height      int
}

func NewManagerHubModel(season *domain.Season, club *domain.Club, fixtures []*domain.Fixture) *ManagerHubModel {
	return &ManagerHubModel{
		ChosenClub: club,
		Season:     season,
		Fixtures:   fixtures,
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
	colWidth := m.width / 3

	if m.ChosenClub == nil {
		panic("ChosenClub is nil, cannot proceed")
	}
	header := lipgloss.NewStyle().
		Align(lipgloss.Center).
		Width(m.width).
		Padding(1, 2).
		Background(lipgloss.Color(m.ChosenClub.Background)).
		Foreground(lipgloss.Color(m.ChosenClub.Foreground)).
		Render("Managing " + m.ChosenClub.Name)
	footer := lipgloss.NewStyle().
		Align(lipgloss.Center).
		Width(m.width).
		Render("Press [Space] to start pre-match")

	leagueTableView := components.Table(m.Season.GetLeagueTable())
	fixturesView := components.Fixtures(m.Fixtures)

	content := lipgloss.JoinHorizontal(
		lipgloss.Top,
		lipgloss.NewStyle().Padding(2).Width(colWidth).Render(fixturesView),
		lipgloss.NewStyle().Padding(2).Width(colWidth).Render("TBC"),
		lipgloss.NewStyle().Padding(2).Width(colWidth).Render(leagueTableView),
	)

	main := lipgloss.NewStyle().
		Width(m.width).
		Height(m.height-lipgloss.Height(header)-lipgloss.Height(footer)).
		Align(lipgloss.Center, lipgloss.Center).
		Render(content)

	return lipgloss.JoinVertical(lipgloss.Top, header, main, footer)

}
