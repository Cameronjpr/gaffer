package tui

import (
	"fmt"

	"github.com/cameronjpr/gaffer/internal/domain"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type ManagerHubModel struct {
	Season *domain.Season
	width  int
	height int
}

func NewManagerHubModel(season *domain.Season) ManagerHubModel {
	return ManagerHubModel{
		Season: season,
	}
}

func (m ManagerHubModel) Init() tea.Cmd {
	return nil
}

func (m ManagerHubModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
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

func (m ManagerHubModel) View() string {
	header := lipgloss.NewStyle().
		Align(lipgloss.Center).
		Width(m.width).
		Border(lipgloss.NormalBorder(), false, false, true, false).
		Render("Manager Hub")
	footer := lipgloss.NewStyle().
		Align(lipgloss.Center).
		Width(m.width).
		Render("Press [Space] to start pre-match")

	leagueTableStr := "League Table:\n"
	for i, position := range m.Season.LeagueTable.Positions {
		leagueTableStr += fmt.Sprintf("%d.  %s\n", i+1, position.Club.Name)
	}

	fixturesStr := "Fixtures:\n"
	for _, gameweek := range m.Season.Gameweeks[:4] { // First 5 fixtures
		for _, fixture := range gameweek.Fixtures {
			fixturesStr += fmt.Sprintf("%s vs %s\n", fixture.HomeTeam.Name, fixture.AwayTeam.Name)
		}
	}

	content := lipgloss.JoinHorizontal(
		lipgloss.Top,
		lipgloss.NewStyle().Padding(2).Render(fixturesStr),
		lipgloss.NewStyle().Padding(2).Render(leagueTableStr),
	)

	main := lipgloss.NewStyle().
		Width(m.width).
		Height(m.height-lipgloss.Height(header)-lipgloss.Height(footer)).
		Align(lipgloss.Center, lipgloss.Center).
		Render(content)

	return lipgloss.JoinVertical(lipgloss.Top, header, main, footer)

}
