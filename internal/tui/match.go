package tui

import (
	"fmt"
	"time"

	"github.com/cameronjpr/gaffer/internal/domain"
	"github.com/cameronjpr/gaffer/internal/simulation"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type MatchModel struct {
	match    *domain.Match
	engine   *simulation.Engine
	isPaused bool
	width    int
	height   int
}

func NewMatchModel(match *domain.Match) MatchModel {
	return MatchModel{
		match:    match,
		engine:   simulation.NewEngine(match),
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

			if m.match.IsHalfTime() {
				m.match.StartSecondHalf()
			}
			return m, tick()
		}
	case tickMsg:
		// Check if match is complete
		if m.match.IsFullTime() {
			return m, func() tea.Msg {
				return matchFinishedMsg{match: m.match}
			}
		}

		// If paused, don't advance
		if m.isPaused {
			return m, nil
		}

		// Auto-advance to next phase
		result := m.engine.PlayPhase()

		// Update match state
		m.match.Home.Score += result.HomeGoals
		m.match.Away.Score += result.AwayGoals
		m.match.CurrentMinute++
		m.match.PhaseHistory = append(m.match.PhaseHistory, result)

		if m.match.IsHalfTime() {
			m.isPaused = true
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
func buildTimelineColumnFromEvents(events []domain.Event, width int, align lipgloss.Position) string {
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

// buildZoneIndicator creates a 5x4 grid showing current active zone
// Layout: 5 rows (lanes: LW, LH, C, RH, RW) × 4 columns (West to East)
func buildZoneIndicator(zone domain.PitchZone, teamInPossession *domain.MatchParticipant) string {
	// Map zones to grid positions (horizontal pitch: West to East, left to right)
	// Row represents lane (0=Left Wing, 4=Right Wing)
	// Column represents depth (0=West, 3=East)
	zoneMap := map[domain.PitchZone][2]int{
		// Column 0: West end
		domain.WestLeftWing:  {0, 0}, // row=Left Wing lane, col=West
		domain.WestLeftHalf:  {1, 0}, // row=Left Half lane, col=West
		domain.WestCentre:    {2, 0}, // row=Centre lane, col=West
		domain.WestRightHalf: {3, 0}, // row=Right Half lane, col=West
		domain.WestRightWing: {4, 0}, // row=Right Wing lane, col=West

		// Column 1: West-Mid
		domain.WestMidLeftWing:  {0, 1},
		domain.WestMidLeftHalf:  {1, 1},
		domain.WestMidCentre:    {2, 1},
		domain.WestMidRightHalf: {3, 1},
		domain.WestMidRightWing: {4, 1},

		// Column 2: East-Mid
		domain.EastMidLeftWing:  {0, 2},
		domain.EastMidLeftHalf:  {1, 2},
		domain.EastMidCentre:    {2, 2},
		domain.EastMidRightHalf: {3, 2},
		domain.EastMidRightWing: {4, 2},

		// Column 3: East end
		domain.EastLeftWing:  {0, 3},
		domain.EastLeftHalf:  {1, 3},
		domain.EastCentre:    {2, 3},
		domain.EastRightHalf: {3, 3},
		domain.EastRightWing: {4, 3},
	}

	// Build 5x4 grid (5 rows for lanes, 4 columns for depth)
	grid := [5][4]string{}
	for z, pos := range zoneMap {
		row, col := pos[0], pos[1]
		if z == zone {
			grid[row][col] = "   ●   " // Active zone
		} else {
			grid[row][col] = "   ·   " // Inactive zone
		}
	}

	// Render grid (5 rows × 4 columns, West to East left-to-right)
	result := ""
	for row := range 5 {
		for col := range 4 {
			result += grid[row][col]
			if col < 3 {
				result += " "
			}
		}
		if row < 4 {
			result += "\n\n"
		}
	}

	return result
}

// buildTimelineFromEvents creates a centered timeline with home and away events
func buildTimelineFromEvents(homeEvents, awayEvents []domain.Event, colWidth int) string {
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
	if !m.match.IsHalfTime() && !m.match.IsFullTime() && len(m.match.Events) > 0 {
		latestEvent := m.match.Events[len(m.match.Events)-1]
		commentary := GenerateCommentary(latestEvent, m.match)

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
	} else if m.match.CurrentHalf == domain.FirstHalf && m.match.IsInAddedTime() {
		timeStr += fmt.Sprintf("+%v'", m.match.GetAddedTime(domain.FirstHalf))
	} else if m.match.CurrentHalf == domain.SecondHalf && m.match.IsInAddedTime() {
		timeStr += fmt.Sprintf("+%v'", m.match.GetAddedTime(domain.SecondHalf))
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

	timeline := buildTimelineFromEvents(homeEvents, awayEvents, colWidth)
	zoneIndicator := buildZoneIndicator(m.match.ActiveZone, m.match.TeamInPossession)
	gap := "  "

	possessionIndicator := lipgloss.NewStyle().
		Width(m.width).
		Border(lipgloss.NormalBorder(), false, false, true, false).
		BorderBottomBackground(lipgloss.Color(m.match.TeamInPossession.Club.Background)).
		Render()

	tickerContent := lipgloss.JoinVertical(
		lipgloss.Center,
		scoreWidget,
		lipgloss.NewStyle().Italic(true).Render(time),
		gap,
		lipgloss.NewStyle().Italic(true).Render(timeline),
		gap,
		lipgloss.NewStyle().Faint(true).Render(zoneIndicator),
	)

	// Home team section
	homeContent := lipgloss.JoinVertical(
		lipgloss.Left,
		lipgloss.NewStyle().Bold(true).Render(m.match.Home.Club.Name),
		lipgloss.NewStyle().Italic(true).Render(m.match.Home.Formation),
		"",
		m.match.Home.GetLineup(m.match),
	)

	// Away team section
	awayContent := lipgloss.JoinVertical(
		lipgloss.Left,
		lipgloss.NewStyle().Bold(true).Render(m.match.Away.Club.Name),
		lipgloss.NewStyle().Italic(true).Render(m.match.Away.Formation),
		"",
		m.match.Away.GetLineup(m.match),
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
		headerPlaced := lipgloss.Place(colWidth, headerHeight, lipgloss.Center, lipgloss.Center, header)
		return lipgloss.JoinVertical(lipgloss.Top, headerPlaced, matchInfo, footerPlaced)
	}

	return matchInfo
}
