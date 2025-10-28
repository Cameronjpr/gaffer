package tui

import (
	"fmt"
	"time"

	"github.com/cameronjpr/gaffer/internal/domain"
	tea "github.com/charmbracelet/bubbletea"
)

type tickMsg time.Time

func (m *AppModel) Init() tea.Cmd {
	return tick()
}

func (m *AppModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
		// WindowSizeMsg will be passed to child models in the switch below

	case tea.KeyMsg:
		// Handle quit
		if msg.Type == tea.KeyCtrlC || msg.Type == tea.KeyEsc {
			return m, tea.Quit
		}

	case goToOnboardingMsg:
		m.mode = OnboardingMode
		// Send WindowSizeMsg to newly activated model
		m.onboarding.width = m.width
		m.onboarding.height = m.height

		return m, tick()

	case goToManagerHubMsg:
		club, err := m.clubRepo.GetByID(msg.ClubID)
		if err != nil {
			return m, tea.Quit
		}

		// Get only unplayed fixtures for the hub display
		fixtures, err := m.fixtureRepo.GetUnplayedByClubID(club.Club.ID)
		if err != nil {
			return m, tea.Quit
		}

		// Calculate league table from database
		leagueTable, err := m.matchRepo.CalculateLeagueTable(m.clubs, m.fixtures)
		if err != nil {
			return m, tea.Quit
		}

		m.managerHub = NewManagerHubModel(club.Club, fixtures, leagueTable)
		m.mode = ManagerHubMode
		m.managerHub.width = m.width
		m.managerHub.height = m.height
		return m, tick()

	case startPreMatchMsg:
		// Get the next fixture for the selected club
		unplayedFixtures, err := m.fixtureRepo.GetUnplayedByClubID(m.managerHub.ChosenClub.ID)
		if err != nil {
			return m, tea.Quit
		}

		if len(unplayedFixtures) == 0 {
			return m, tea.Quit
		}

		nextMatch := domain.NewMatchFromFixture(unplayedFixtures[0])
		m.currentMatch = nextMatch

		// Update the prematch and match models with the new match
		m.prematch = NewPreMatchModel(m.currentMatch)
		m.match = NewMatchModel(m.currentMatch)

		m.mode = PreMatchMode
		// Send WindowSizeMsg to newly activated model
		m.prematch.width = m.width
		m.prematch.height = m.height
		return m, tick()

	case startMatchMsg:

		// Create the match record in the database
		err := m.matchRepo.Create(m.currentMatch)
		if err != nil {
			// Log error but continue - in production you'd handle this better
			// For now, we still run the match even if DB create fails
		}

		m.mode = MatchMode
		// Send WindowSizeMsg to newly activated model
		m.match.width = m.width
		m.match.height = m.height
		// Initialize the match model (starts controller and begins listening)
		return m, m.match.Init()

	case matchFinishedMsg:
		match := msg.match
		match.ForFixture.Result = msg.match

		// Save the match result to the database
		err := m.matchRepo.SaveResult(match)
		if err != nil {
			fmt.Println("Error saving match result:", err)
		}

		// Refresh the unplayed fixtures list
		unplayedFixtures, err := m.fixtureRepo.GetUnplayedByClubID(m.managerHub.ChosenClub.ID)
		if err != nil {
			return m, tea.Quit
		}

		if len(unplayedFixtures) == 0 {
			return m, tea.Quit
		}

		nextMatch := domain.NewMatchFromFixture(unplayedFixtures[0])
		m.currentMatch = nextMatch

		// Update the child models to point to the new match
		m.prematch = NewPreMatchModel(m.currentMatch)
		m.match = NewMatchModel(m.currentMatch)

		// Recalculate league table with latest results
		leagueTable, err := m.matchRepo.CalculateLeagueTable(m.clubs, m.fixtures)
		if err != nil {
			// Log error but continue
		} else {
			m.managerHub.LeagueTable = leagueTable
		}

		// Update the hub's fixture list to remove completed fixtures
		m.managerHub.Fixtures = unplayedFixtures

		// Go back to the hub
		m.mode = ManagerHubMode
		return m, tick()
	}

	var cmd tea.Cmd
	switch m.mode {
	case MenuMode:
		var newMenu tea.Model
		newMenu, cmd = m.menu.Update(msg)
		m.menu = newMenu.(*MenuModel)

	case OnboardingMode:
		var newOnboarding tea.Model
		newOnboarding, cmd = m.onboarding.Update(msg)
		m.onboarding = newOnboarding.(*OnboardingModel)

	case ManagerHubMode:
		var newManagerHub tea.Model
		newManagerHub, cmd = m.managerHub.Update(msg)
		m.managerHub = newManagerHub.(*ManagerHubModel)

	case PreMatchMode:
		var newPrematch tea.Model
		newPrematch, cmd = m.prematch.Update(msg)
		m.prematch = newPrematch.(*PreMatchModel)

	case MatchMode:
		var newMatch tea.Model
		newMatch, cmd = m.match.Update(msg)
		m.match = newMatch.(*MatchModel)

		// Check if match is complete
		if m.match.match.IsFullTime() {
			// m.mode = ResultsMode
		}

	}

	return m, cmd
}

func (m *AppModel) View() string {
	switch m.mode {
	case MenuMode:
		return m.menu.View()
	case OnboardingMode:
		return m.onboarding.View()
	case ManagerHubMode:
		return m.managerHub.View()
	case PreMatchMode:
		return m.prematch.View()
	case MatchMode:
		return m.match.View()
	}
	return "No mode"
}
