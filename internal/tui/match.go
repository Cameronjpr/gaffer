package tui

import (
	"github.com/cameronjpr/gaffer/internal/components"
	"github.com/cameronjpr/gaffer/internal/domain"
	"github.com/cameronjpr/gaffer/internal/simulation"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// MatchTab represents the different tabs in the match interface
type MatchTab int

const (
	MatchViewTab MatchTab = iota
	SubstitutionsTab
	TacticsTab
)

type MatchModel struct {
	match             *domain.Match
	controller        *simulation.MatchController
	latestEvent       *domain.Event            // For commentary and timeline
	userTeam          *domain.MatchParticipant // Direct pointer to user's controlled team
	width             int
	height            int
	currentTab        MatchTab
	substitutionModel *SubstitutionModel
	tacticsModel      *TacticsModel
}

func NewMatchModel(match *domain.Match, userClubID int64) *MatchModel {
	controller := simulation.NewMatchController(match)

	// Determine which team the user controls (handle nil match during initialization)
	var userTeam *domain.MatchParticipant
	var substitutionModel *SubstitutionModel
	var tacticsModel *TacticsModel

	if match != nil {
		if match.Home.Club.ID == userClubID {
			userTeam = match.Home
		} else {
			userTeam = match.Away
		}

		// Initialize sub-models only if we have a valid userTeam
		if userTeam != nil {
			substitutionModel = NewSubstitutionModel(userTeam.CurrentXI, userTeam.Bench)
			tacticsModel = NewTacticsModel(userTeam)
		}
	}

	// Don't start controller yet - wait for Init()
	return &MatchModel{
		match:             match,
		controller:        controller,
		userTeam:          userTeam,
		latestEvent:       nil,
		currentTab:        MatchViewTab, // Start on match view
		substitutionModel: substitutionModel,
		tacticsModel:      tacticsModel,
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
	// Handle substitution execution message
	switch msg := msg.(type) {
	case executeSubstitutionMsg:
		// Send substitution command to controller
		m.controller.SendCommand(simulation.SubstitutePlayerCmd{
			Participant: m.userTeam,
			PlayerOut:   msg.playerOut,
			PlayerIn:    msg.playerIn,
		})
		// Switch back to match view after substitution
		m.currentTab = MatchViewTab
		return m, waitForMatchEvent(m.controller)
	}

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		// Forward to sub-models if they exist
		if m.substitutionModel != nil {
			m.substitutionModel.Update(msg)
		}
		if m.tacticsModel != nil {
			m.tacticsModel.Update(msg)
		}
		return m, nil

	case tea.KeyMsg:
		// Global tab switching (m/s/t keys)
		switch msg.Type {
		case tea.KeyRunes:
			switch msg.Runes[0] {
			case 'm':
				// Switch to match view
				m.currentTab = MatchViewTab
				return m, nil
			case 's':
				// Switch to substitutions tab
				if m.currentTab != SubstitutionsTab {
					m.currentTab = SubstitutionsTab
					// Pause match when entering subs tab
					m.controller.SendCommand(simulation.PauseMatchCmd{})
				}
				return m, nil
			case 't':
				// Switch to tactics tab
				if m.currentTab != TacticsTab {
					m.currentTab = TacticsTab
					// Pause match when entering tactics tab
					m.controller.SendCommand(simulation.PauseMatchCmd{})
				}
				return m, nil
			}
		}

		// Forward key events to the active tab
		if m.currentTab == SubstitutionsTab && m.substitutionModel != nil {
			var cmd tea.Cmd
			_, cmd = m.substitutionModel.Update(msg)
			return m, cmd
		} else if m.currentTab == TacticsTab && m.tacticsModel != nil {
			var cmd tea.Cmd
			_, cmd = m.tacticsModel.Update(msg)
			return m, cmd
		}

		// Match view specific controls
		if m.currentTab == MatchViewTab {
			switch msg.Type {
			case tea.KeySpace:
				m.controller.SendCommand(simulation.TogglePausedCmd{})
				return m, waitForMatchEvent(m.controller)
			case tea.KeyRight:
				m.controller.SendCommand(simulation.SpeedUpCmd{})
				return m, waitForMatchEvent(m.controller)
			case tea.KeyLeft:
				m.controller.SendCommand(simulation.SlowDownCmd{})
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
		// Switch back to match view after substitution executes
		m.currentTab = MatchViewTab
		return m, waitForMatchEvent(m.controller)
	}

	return m, nil
}

func (m *MatchModel) View() string {
	// Shared header with game state
	header := components.MatchHeader(m.width, m.match, m.userTeam)

	// Content and footer based on current tab
	var content string
	var hotkeys []components.HotkeyBinding

	switch m.currentTab {
	case MatchViewTab:
		content = m.renderMatchView()
		hotkeys = []components.HotkeyBinding{
			{Key: "[S]", Description: "Subs"},
			{Key: "[T]", Description: "Tactics"},
			{Key: "Space", Description: "Pause"},
			{Key: "←→", Description: "Speed"},
			{Key: "Esc", Description: "Menu"},
		}

	case SubstitutionsTab:
		if m.substitutionModel != nil {
			content = m.substitutionModel.View()
		} else {
			content = "Substitution view not available"
		}
		hotkeys = []components.HotkeyBinding{
			{Key: "[M]", Description: "Match"},
			{Key: "[T]", Description: "Tactics"},
			{Key: "←→/Tab", Description: "Switch"},
			{Key: "Space", Description: "Select"},
			{Key: "Enter", Description: "Confirm"},
		}

	case TacticsTab:
		if m.tacticsModel != nil {
			content = m.tacticsModel.View()
		} else {
			content = "Tactics view not available"
		}
		hotkeys = []components.HotkeyBinding{
			{Key: "[M]", Description: "Match"},
			{Key: "[S]", Description: "Subs"},
		}
	}

	footer := components.HotkeyGuide(m.width, hotkeys)

	// Calculate heights
	headerHeight := lipgloss.Height(header)
	footerHeight := lipgloss.Height(footer)
	contentHeight := m.height - headerHeight - footerHeight

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

// renderMatchView renders the main match view tab
func (m *MatchModel) renderMatchView() string {
	if m.match == nil {
		return "Match not initialized"
	}

	colWidth := m.width / 3

	// Main match content - three columns
	homeTeamSheet := components.TeamSheet(m.match.Home)
	awayTeamSheet := components.TeamSheet(m.match.Away)
	matchActionView := components.MatchActionView(colWidth, m.controller.GetSpeed(), m.match)

	// Add padding/spacing between columns
	spacer := "  "

	matchContent := lipgloss.JoinHorizontal(
		lipgloss.Top,
		homeTeamSheet,
		spacer,
		matchActionView,
		spacer,
		awayTeamSheet,
	)

	return matchContent
}
