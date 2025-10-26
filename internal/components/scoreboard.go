package components

import (
	"fmt"

	"github.com/cameronjpr/gaffer/internal/domain"
	"github.com/charmbracelet/lipgloss"
)

func Scoreboard(match *domain.Match, width int) string {
	scoreStr := fmt.Sprintf("%s [%v] – [%v] %s", match.Home.Club.Name, match.Home.Score, match.Away.Score, match.Away.Club.Name)
	return lipgloss.NewStyle().
		Width(width).
		Align(lipgloss.Center).
		Bold(true).
		Render(scoreStr)
}
