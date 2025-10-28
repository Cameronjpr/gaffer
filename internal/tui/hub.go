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
		case tea.KeyEnter:
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

	// Header with club branding
	header := lipgloss.NewStyle().
		Align(lipgloss.Center).
		Width(m.width).
		Padding(1, 2).
		Bold(true).
		Background(lipgloss.Color(m.ChosenClub.Background)).
		Foreground(lipgloss.Color(m.ChosenClub.Foreground)).
		Render(m.ChosenClub.Name)

	// Footer with instructions
	footer := lipgloss.NewStyle().
		Align(lipgloss.Center).
		Width(m.width).
		Render("Press [Enter] to start pre-match")

	// Main content area - calculate flexible content height
	headerHeight := lipgloss.Height(header)
	footerHeight := lipgloss.Height(footer)
	contentHeight := m.height - headerHeight - footerHeight

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

	// Center content in available space
	centeredContent := components.Centered(m.width, contentHeight, content)

	// Use ScreenLayout to organize header, content, footer
	sections := []components.ScreenSection{
		{Height: headerHeight, Content: header},
		{Height: contentHeight, Content: centeredContent},
		{Height: footerHeight, Content: footer},
	}

	return components.ScreenLayout(m.height, sections)
}
