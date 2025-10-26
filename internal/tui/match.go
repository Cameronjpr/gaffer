package tui

import (
	"fmt"

	"github.com/cameronjpr/gaffer/internal/components"
	"github.com/cameronjpr/gaffer/internal/domain"
	"github.com/cameronjpr/gaffer/internal/simulation"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type MatchModel struct {
	match       *domain.Match
	controller  *simulation.MatchController
	latestEvent *domain.Event // For commentary and timeline
	width       int
	height      int
}

func NewMatchModel(match *domain.Match) *MatchModel {
	controller := simulation.NewMatchController(match)
	// Don't start controller yet - wait for Init()

	return &MatchModel{
		match:       match,
		controller:  controller,
		latestEvent: nil,
	}
}

func (m *MatchModel) Init() tea.Cmd {
	// Start the controller goroutine now that TUI is ready
	go m.controller.Run()

	// Start listening for match events
	return waitForMatchEvent(m.controller)
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

	case simulation.MatchUpdateMsg:
		// Update match state from controller
		m.match = msg.Match
		m.latestEvent = msg.LatestEvent

		// Continue listening for next event
		return m, waitForMatchEvent(m.controller)

	case simulation.HalftimeMsg:
		// Controller auto-paused at halftime
		m.match = msg.Match
		m.match.StartSecondHalf() // Prepare for second half
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
	}
	return m, nil
}

func (m *MatchModel) View() string {
	// Footer - generate commentary on-demand from latest event
	footer := ""
	if !m.match.IsHalfTime() && !m.match.IsFullTime() && m.latestEvent != nil {
		commentary := GenerateCommentary(*m.latestEvent, m.match)

		style := lipgloss.NewStyle().Align(lipgloss.Center).Width(m.width)

		// Use the commentary's For field to determine styling
		if commentary.EventType == domain.GoalEvent && commentary.For != nil {
			style = style.Bold(true).
				Background(lipgloss.Color(commentary.For.Club.Background)).
				Foreground(lipgloss.Color(commentary.For.Club.Foreground))
		} else {
			style = style.Background(lipgloss.Color("#000000"))
		}

		footer = style.Render(commentary.Message)
	}

	// Calculate column width
	colWidth := m.width / 3
	timeStr := fmt.Sprintf("(%v:00)", m.match.CurrentMinute)
	if m.match.IsHalfTime() {
		timeStr = "HT"
	} else if m.match.IsFullTime() {
		timeStr = "FT"
	} else if m.match.IsFirstHalf() && m.match.IsInAddedTime() {
		timeStr += fmt.Sprintf("+%v'", m.match.GetAddedTime(domain.FirstHalf))
	} else if m.match.IsSecondHalf() && m.match.IsInAddedTime() {
		timeStr += fmt.Sprintf("+%v'", m.match.GetAddedTime(domain.SecondHalf))
	}
	time := lipgloss.NewStyle().
		Padding(0, 1).
		Render(timeStr)

	scoreWidget := components.Scoreboard(
		m.match,
		colWidth,
	)

	// Filter events for each team's timeline (only show goal events)
	var homeEvents []domain.Event
	var awayEvents []domain.Event
	for _, event := range m.match.Events {
		if event.Type == domain.GoalEvent {
			if event.For == m.match.Home {
				homeEvents = append(homeEvents, event)
			} else if event.For == m.match.Away {
				awayEvents = append(awayEvents, event)
			}
		}
	}

	timelineView := components.EventsTimeline(homeEvents, awayEvents, colWidth)
	pitchView := components.Pitch(m.match.ActiveZone, m.match)
	gap := "  "

	possessionIndicator := lipgloss.NewStyle().
		Width(m.width).
		Border(lipgloss.NormalBorder(), false, false, true, false).
		BorderBottomBackground(lipgloss.Color(m.match.TeamInPossession.Club.Background)).
		Render()

	tickerContent := lipgloss.JoinVertical(
		lipgloss.Center,
		fmt.Sprintf("Speed: %s", m.controller.GetSpeed()),
		scoreWidget,
		lipgloss.NewStyle().Italic(true).Render(time),
		gap,
		timelineView,
		gap,
		pitchView,
	)

	header := lipgloss.JoinVertical(
		lipgloss.Center,
		possessionIndicator,
	)

	// Calculate heights - footer takes 1 line, rest is for match info
	headerHeight := 1
	footerHeight := 1
	matchInfoHeight := m.height - footerHeight - headerHeight
	if matchInfoHeight < 0 {
		matchInfoHeight = 0
	}

	homeTeamSheet := components.TeamSheet(m.match.Home)
	awayTeamSheet := components.TeamSheet(m.match.Away)
	matchInfo := lipgloss.JoinHorizontal(
		lipgloss.Top,
		lipgloss.Place(colWidth, matchInfoHeight, lipgloss.Center, lipgloss.Center, homeTeamSheet),
		lipgloss.Place(colWidth, matchInfoHeight, lipgloss.Center, lipgloss.Center, tickerContent),
		lipgloss.Place(colWidth, matchInfoHeight, lipgloss.Center, lipgloss.Center, awayTeamSheet),
	)

	// Place footer at the bottom of the screen
	if footer != "" {
		footerPlaced := lipgloss.Place(m.width, footerHeight, lipgloss.Center, lipgloss.Bottom, footer)
		headerPlaced := lipgloss.Place(colWidth, headerHeight, lipgloss.Center, lipgloss.Center, header)
		return lipgloss.JoinVertical(lipgloss.Top, headerPlaced, matchInfo, footerPlaced)
	}

	return matchInfo
}
