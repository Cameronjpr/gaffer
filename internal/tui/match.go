package tui

import (
	"fmt"
	"time"

	"github.com/cameronjpr/gaffer/internal/game"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type MatchModel struct {
	match    *game.Match
	isPaused bool
	width    int
	height   int
}

func NewMatchModel(match game.Match) MatchModel {
	return MatchModel{
		match:    &match,
		isPaused: false,
	}
}

func (m MatchModel) Init() tea.Cmd {
	return nil
}

func (m MatchModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {

		case tea.KeySpace:
			m.isPaused = !m.isPaused
			if m.isPaused {
				return m, nil
			}

			return m, tick()
		}
	case tickMsg:
		// Check if match is complete
		if m.match.CurrentMinute >= 90 {
			return m, nil
		}

		// If paused, don't advance
		if m.isPaused {
			return m, nil
		}

		// Auto-advance to next phase
		result := m.match.PlayPhase()

		// Update match state
		m.match.Home.Score += result.HomeGoals
		m.match.Away.Score += result.AwayGoals
		m.match.CurrentMinute++
		m.match.PhaseHistory = append(m.match.PhaseHistory, result)

		if m.match.IsHalfTime() {
			m.match.StartSecondHalf()
			return m, nil
		}

		// Check if a goal was scored in this phase
		goalScored := result.HomeGoals > 0 || result.AwayGoals > 0

		// If a goal was scored, pause longer to let user see it
		if goalScored {
			return m, tickWithDuration(time.Second * 3)
		}

		if m.match.IsInAddedTime() {
			return m, tickWithDuration(time.Second * 2)
		}

		// Return the next tick command to keep the loop going
		return m, tick()
	}
	return m, nil
}

// buildScoreWidget creates a centered score widget with team names and scores
func buildScoreWidget(homeTeam, awayTeam string, homeScore, awayScore, width int) string {
	scoreStr := fmt.Sprintf("%s [%v] – [%v] %s", homeTeam, homeScore, awayScore, awayTeam)
	return lipgloss.NewStyle().
		Width(width).
		Align(lipgloss.Center).
		Bold(true).
		Render(scoreStr)
}

// buildTimelineColumnFromEvents builds a single timeline column with events
func buildTimelineColumnFromEvents(events []game.Event, width int, align lipgloss.Position) string {
	timeline := ""
	for _, event := range events {
		timeline += event.String()
		timeline += "\n"
	}
	return lipgloss.NewStyle().
		Width(width).
		Align(align).
		Render(timeline)
}

// buildTimelineFromEvents creates a centered timeline with home and away events
func buildTimelineFromEvents(homeEvents, awayEvents []game.Event, colWidth int) string {
	// Calculate timeline column width (half of ticker width minus gap)
	timelineWidth := (colWidth / 2) - 2

	homeTimelineStyled := buildTimelineColumnFromEvents(homeEvents, timelineWidth, lipgloss.Right)
	awayTimelineStyled := buildTimelineColumnFromEvents(awayEvents, timelineWidth, lipgloss.Left)

	// Add gap between timelines and center the entire timeline
	gap := "  "
	timelineContent := lipgloss.JoinHorizontal(lipgloss.Top, homeTimelineStyled, gap, awayTimelineStyled)
	return lipgloss.NewStyle().
		Width(colWidth).
		Align(lipgloss.Center).
		Render(timelineContent)
}

func (m MatchModel) View() string {
	// Footer - generate commentary on-demand from latest event
	footer := ""
	if !m.match.IsHalfTime() && len(m.match.Events) > 0 {
		latestEvent := m.match.Events[len(m.match.Events)-1]
		commentary := game.GenerateCommentary(latestEvent, m.match)

		style := lipgloss.NewStyle().Align(lipgloss.Center).Width(m.width)

		// Use the commentary's For field to determine styling
		if commentary.EventType == game.GoalEvent && commentary.For != nil {
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
	} else if m.match.CurrentHalf == 1 && m.match.IsInAddedTime() {
		timeStr += fmt.Sprintf("+%v'", m.match.GetAddedTime(game.FirstHalf))
	} else if m.match.CurrentHalf == 2 && m.match.IsInAddedTime() {
		timeStr += fmt.Sprintf("+%v'", m.match.GetAddedTime(game.SecondHalf))
	}
	time := lipgloss.NewStyle().
		Padding(0, 1).
		Render(timeStr)

	// Build score widget using helper function
	scoreWidget := buildScoreWidget(
		m.match.Home.Club.Name,
		m.match.Away.Club.Name,
		m.match.Home.Score,
		m.match.Away.Score,
		colWidth,
	)

	// Filter events for each team's timeline (only show goal events)
	var homeEvents []game.Event
	var awayEvents []game.Event
	for _, event := range m.match.Events {
		if event.Type == game.GoalEvent {
			if event.For == m.match.Home {
				homeEvents = append(homeEvents, event)
			} else if event.For == m.match.Away {
				awayEvents = append(awayEvents, event)
			}
		}
	}

	timeline := buildTimelineFromEvents(homeEvents, awayEvents, colWidth)
	gap := "  "

	tickerContent := lipgloss.JoinVertical(
		lipgloss.Center,
		scoreWidget,
		lipgloss.NewStyle().Italic(true).Render(time),
		gap,
		lipgloss.NewStyle().Italic(true).Render(timeline),
	)

	// Home team section
	homeContent := lipgloss.JoinVertical(
		lipgloss.Left,
		lipgloss.NewStyle().Bold(true).Render(m.match.Home.Club.Name),
		lipgloss.NewStyle().Italic(true).Render(m.match.Home.Formation),
		"",
		m.match.Home.GetLineup(),
	)

	// Away team section
	awayContent := lipgloss.JoinVertical(
		lipgloss.Left,
		lipgloss.NewStyle().Bold(true).Render(m.match.Away.Club.Name),
		lipgloss.NewStyle().Italic(true).Render(m.match.Away.Formation),
		"",
		m.match.Away.GetLineup(),
	)

	// Calculate heights - footer takes 1 line, rest is for match info
	footerHeight := 1
	if footer == "" {
		footerHeight = 0
	}
	matchInfoHeight := m.height - footerHeight
	if matchInfoHeight < 0 {
		matchInfoHeight = 0
	}

	// Create 3-column layout
	matchInfo := lipgloss.JoinHorizontal(
		lipgloss.Top,
		lipgloss.Place(colWidth, matchInfoHeight, lipgloss.Center, lipgloss.Center, homeContent),
		lipgloss.Place(colWidth, matchInfoHeight, lipgloss.Center, lipgloss.Center, tickerContent),
		lipgloss.Place(colWidth, matchInfoHeight, lipgloss.Center, lipgloss.Center, awayContent),
	)

	// Place footer at the bottom of the screen
	if footer != "" {
		footerPlaced := lipgloss.Place(m.width, footerHeight, lipgloss.Center, lipgloss.Bottom, footer)
		return lipgloss.JoinVertical(lipgloss.Top, matchInfo, footerPlaced)
	}

	return matchInfo
}
