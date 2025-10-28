package components

import (
	"fmt"

	"github.com/cameronjpr/gaffer/internal/domain"
	"github.com/charmbracelet/lipgloss"
)

func MatchActionView(width int, speed string, match *domain.Match) string {
	gap := " "

	return lipgloss.JoinVertical(
		lipgloss.Center,
		fmt.Sprintf("Speed: %s", speed),
		Scoreboard(
			match,
			width,
		),
		Clock(match),
		gap,
		EventsTimeline(match, width),
		gap,
		Pitch(match),
		gap,
		CommentaryBar(match, width),
	)
}
