package tui

import (
	"fmt"

	"github.com/cameronjpr/gaffer/internal/domain"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type ManagerHubModel struct {
	Season     *domain.Season
	ChosenClub *domain.Club
	width      int
	height     int
}

func NewManagerHubModel(season *domain.Season, club *domain.Club) ManagerHubModel {
	return ManagerHubModel{
		ChosenClub: club,
		Season:     season,
	}
}

func (m ManagerHubModel) Init() tea.Cmd {
	return nil
}

func (m ManagerHubModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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

func (m ManagerHubModel) View() string {
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

	leagueTableStr := "Pos.  |  Name  |  Played  |  Points\n"
	for i, position := range m.Season.GetLeagueTable().Positions {

		leagueTableStr += fmt.Sprintf("%d. |  %s  |  %d  |  %d pts\n", i+1, position.Club.Name, position.Played, position.Points)
	}

	fixturesStr := "Fixtures:\n"
	clubFixtures := m.Season.GetFixturesForClub(m.ChosenClub)
	// Show up to 5 fixtures
	numToShow := len(clubFixtures)
	if numToShow > 5 {
		numToShow = 5
	}
	if numToShow == 0 {
		fixturesStr += "No fixtures scheduled\n"
	} else {
		for _, fixture := range clubFixtures[:numToShow] {
			status := ""
			if fixture.Result != nil {
				status = fmt.Sprintf(" %d-%d", fixture.Result.Home.Score, fixture.Result.Away.Score)
			}
			fixturesStr += fmt.Sprintf("%s vs %s%s\n", fixture.HomeTeam.Name, fixture.AwayTeam.Name, status)
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
