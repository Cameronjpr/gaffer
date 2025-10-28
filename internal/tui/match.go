package tui

import (
	"github.com/cameronjpr/gaffer/internal/components"
	"github.com/cameronjpr/gaffer/internal/domain"
	"github.com/cameronjpr/gaffer/internal/simulation"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type MatchModel struct {
	match            *domain.Match
	controller       *simulation.MatchController
	latestEvent      *domain.Event            // For commentary and timeline
	userTeam         *domain.MatchParticipant // Direct pointer to user's controlled team
	width            int
	height           int
	showSubModal     bool
	showTacticsModal bool
}

func NewMatchModel(match *domain.Match, userClubID int64) *MatchModel {
	controller := simulation.NewMatchController(match)

	// Determine which team the user controls (handle nil match during initialization)
	var userTeam *domain.MatchParticipant
	if match != nil {
		if match.Home.Club.ID == userClubID {
			userTeam = match.Home
		} else {
			userTeam = match.Away
		}
	}

	// Don't start controller yet - wait for Init()
	return &MatchModel{
		match:       match,
		controller:  controller,
		userTeam:    userTeam,
		latestEvent: nil,
	}
}

func (m *MatchModel) Init() tea.Cmd {
	// Start the controller goroutine now that TUI is ready
	go m.controller.Run()

	// Start listening for match events
	return waitForMatchEvent(m.controller)
}

// GetOpponentTeam returns the team the user is playing against
func (m *MatchModel) GetOpponentTeam() *domain.MatchParticipant {
	if m.userTeam == m.match.Home {
		return m.match.Away
	}
	return m.match.Home
}

// IsUserControlled checks if the given participant is the user's team
func (m *MatchModel) IsUserControlled(participant *domain.MatchParticipant) bool {
	return participant == m.userTeam
}

// waitForMatchEvent creates a Cmd that blocks until the controller sends an event
func waitForMatchEvent(controller *simulation.MatchController) tea.Cmd {
	return func() tea.Msg {
		// This blocks until an event is received
		// BubbleTea runs this in a goroutine, so blocking is safe
		return <-controller.EventChan()
	}
}

func (m *MatchModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			if m.showSubModal {
				m.controller.SendCommand(simulation.SubstitutePlayerCmd{
					Participant: m.userTeam,
					PlayerIn:    m.userTeam.Bench[0],
					PlayerOut:   m.userTeam.CurrentXI[0],
				})
			}
			return m, waitForMatchEvent(m.controller)
		case tea.KeySpace:
			m.controller.SendCommand(simulation.TogglePausedCmd{})
			return m, waitForMatchEvent(m.controller)
		case tea.KeyRight:
			m.controller.SendCommand(simulation.SpeedUpCmd{})
			return m, waitForMatchEvent(m.controller)
		case tea.KeyLeft:
			m.controller.SendCommand(simulation.SlowDownCmd{})
			return m, waitForMatchEvent(m.controller)
		case tea.KeyRunes:
			switch msg.Runes[0] {
			case 's':
				m.showSubModal = !m.showSubModal
				return m, waitForMatchEvent(m.controller)
			case 't':
				m.showTacticsModal = !m.showTacticsModal
				return m, waitForMatchEvent(m.controller)
			}
		}

	case simulation.MatchUpdateMsg:
		// Update match state from controller
		m.match = msg.Match
		m.latestEvent = msg.LatestEvent

		// Continue listening for next event
		return m, waitForMatchEvent(m.controller)

	case simulation.HalftimeMsg:
		// Controller auto-paused at halftime
		m.match = msg.Match
		// Continue listening - user will press space to resume
		return m, waitForMatchEvent(m.controller)

	case simulation.FulltimeMsg:
		// Match finished
		m.match = msg.Match
		return m, func() tea.Msg {
			return matchFinishedMsg{match: m.match}
		}

	case simulation.MatchPausedMsg:
		// User paused the match
		m.match = msg.Match
		// Continue listening - user will press space to resume
		return m, waitForMatchEvent(m.controller)

	case simulation.MatchResumedMsg:
		// User resumed the match
		m.match = msg.Match
		// Continue listening - user will press space to pause
		return m, waitForMatchEvent(m.controller)

	case simulation.SubstitutionMadeMsg:
		m.match = msg.Match
		return m, waitForMatchEvent(m.controller)
	}

	return m, nil
}

func (m *MatchModel) View() string {
	colWidth := m.width / 3

	// Show substitution modal if active
	if m.showSubModal {
		// TODO: Implement substitution modal content
		modalContent := components.SubstitutionsModal(m.userTeam.CurrentXI, m.userTeam.Bench)
		return components.SimpleModal(m.width, m.height, "Make Substitution – press [S] to close", modalContent)
	}

	if m.showTacticsModal {
		modalContent := "Coming soon"
		return components.SimpleModal(m.width, m.height, "Adjust Tactics – press [T] to close", modalContent)
	}

	// Possession indicator header
	possessionIndicator := lipgloss.NewStyle().
		Width(m.width).
		Border(lipgloss.NormalBorder(), false, false, true, false).
		BorderBottomBackground(lipgloss.Color(m.match.TeamInPossession.Club.Background)).
		Render()

	// Footer with controls
	footer := components.MatchToolbar(m.width, m.match)

	// Calculate heights
	headerHeight := lipgloss.Height(possessionIndicator)
	footerHeight := lipgloss.Height(footer)
	contentHeight := m.height - headerHeight - footerHeight

	// Main match content - three columns
	homeTeamSheet := components.TeamSheet(m.match.Home)
	awayTeamSheet := components.TeamSheet(m.match.Away)
	matchActionView := components.MatchActionView(colWidth, m.controller.GetSpeed(), m.match)

	matchContent := lipgloss.JoinHorizontal(
		lipgloss.Top,
		homeTeamSheet,
		matchActionView,
		awayTeamSheet,
	)

	// Center content vertically and horizontally in available space
	matchContent = components.Centered(m.width, contentHeight, matchContent)

	// Use ScreenLayout to organize header, content, footer
	sections := []components.ScreenSection{
		{Height: headerHeight, Content: possessionIndicator},
		{Height: contentHeight, Content: matchContent},
		{Height: footerHeight, Content: footer},
	}

	return components.ScreenLayout(m.height, sections)
}
