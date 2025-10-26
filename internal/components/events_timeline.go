package components

import (
	"github.com/cameronjpr/gaffer/internal/domain"
	"github.com/charmbracelet/lipgloss"
)

// buildTimelineFromEvents creates a centered timeline with home and away events
func EventsTimeline(homeEvents, awayEvents []domain.Event, colWidth int) string {
	// Calculate timeline column width (half of ticker width minus gap)
	timelineWidth := (colWidth / 2) - 2

	homeTimelineStyled := EventsTimelineForTeam(homeEvents, timelineWidth, lipgloss.Right)
	awayTimelineStyled := EventsTimelineForTeam(awayEvents, timelineWidth, lipgloss.Left)

	// Add gap between timelines and center the entire timeline
	gap := "  "
	timelineContent := lipgloss.JoinHorizontal(lipgloss.Top, homeTimelineStyled, gap, awayTimelineStyled)
	return lipgloss.NewStyle().
		Width(colWidth).
		Italic(true).
		Align(lipgloss.Center).
		Render(timelineContent)
}

func EventsTimelineForTeam(events []domain.Event, width int, align lipgloss.Position) string {
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
