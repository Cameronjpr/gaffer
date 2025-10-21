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
		if m.match.CurrentPhase >= 90 {
			m.match.IsComplete = true
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
		m.match.CurrentPhase++
		m.match.PhaseHistory = append(m.match.PhaseHistory, result)

		if m.match.IsHalfTime() {
			m.isPaused = true
			m.match.CurrentHalf++
			return m, nil
		}

		// Check if a goal was scored in this phase
		goalScored := result.HomeGoals > 0 || result.AwayGoals > 0

		// If a goal was scored, pause longer to let user see it
		if goalScored {
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

// buildTimelineColumn builds a single timeline column with events
func buildTimelineColumn(events []game.PlayerEvent, width int, align lipgloss.Position) string {
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

// buildTimeline creates a centered timeline with home and away events
func buildTimeline(homeEvents, awayEvents []game.PlayerEvent, colWidth int) string {
	// Calculate timeline column width (half of ticker width minus gap)
	timelineWidth := (colWidth / 2) - 2

	homeTimelineStyled := buildTimelineColumn(homeEvents, timelineWidth, lipgloss.Right)
	awayTimelineStyled := buildTimelineColumn(awayEvents, timelineWidth, lipgloss.Left)

	// Add gap between timelines and center the entire timeline
	gap := "  "
	timelineContent := lipgloss.JoinHorizontal(lipgloss.Top, homeTimelineStyled, gap, awayTimelineStyled)
	return lipgloss.NewStyle().
		Width(colWidth).
		Align(lipgloss.Center).
		Render(timelineContent)
}

func (m MatchModel) View() string {
	// Header

	// Footer
	footer := ""
	if !m.match.IsHalfTime() && len(m.match.Commentary) > 0 {
		latestCommentaryMsg :=
			m.match.Commentary[len(m.match.Commentary)-1]
		style :=
			lipgloss.NewStyle().Align(lipgloss.Center).Width(m.width)

		if latestCommentaryMsg.Flash {
			style =
				style.Bold(true).Background(lipgloss.Color("#ff0000"))
		} else {
			style = style.Background(lipgloss.Color("#000000"))
		}

		footer = style.Render(latestCommentaryMsg.Message)
	}

	// Calculate column width
	colWidth := m.width / 3
	timeStr := fmt.Sprintf("(%v:00)", m.match.CurrentPhase)
	if m.match.CurrentPhase == 45 {
		timeStr = "HT"
	} else if m.match.CurrentPhase == 90 {
		timeStr = "FT"
	}
	time := lipgloss.NewStyle().
		Padding(0, 1).
		Render(timeStr)

	// Build score widget and timeline using helper functions
	scoreWidget := buildScoreWidget(
		m.match.Home.Club.Name,
		m.match.Away.Club.Name,
		m.match.Home.Score,
		m.match.Away.Score,
		colWidth,
	)
	timeline := buildTimeline(m.match.Home.PlayerEvents, m.match.Away.PlayerEvents, colWidth)
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

	// Calculate matchInfo height (total height - header - footer)
	// Header and footer each take 1 line
	matchInfoHeight := m.height - 2
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

	return lipgloss.JoinVertical(lipgloss.Top, matchInfo, footer)
}
