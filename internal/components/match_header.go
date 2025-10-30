package components

import (
	"fmt"

	"github.com/cameronjpr/gaffer/internal/domain"
	"github.com/charmbracelet/lipgloss"
)

// MatchHeader renders the persistent header for match screens
// Displays: Club name | Score | Match time
func MatchHeader(width int, match *domain.Match, userTeam *domain.MatchParticipant) string {
	if match == nil {
		return ""
	}

	// Determine home and away for score display
	homeClub := match.Home.Club.Name
	awayClub := match.Away.Club.Name
	homeScore := match.Home.Score
	awayScore := match.Away.Score

	// Left section: Your club name
	leftSection := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("170")).
		Render(fmt.Sprintf(" %s", userTeam.Club.Name))

	// Center section: Score
	scoreText := fmt.Sprintf("%s %d - %d %s", homeClub, homeScore, awayScore, awayClub)
	centerSection := lipgloss.NewStyle().
		Bold(true).
		Align(lipgloss.Center).
		Render(scoreText)

	// Right section: Match time
	minutes := match.CurrentMinute
	half := "1st Half"
	if match.CurrentHalf == domain.SecondHalf {
		half = "2nd Half"
	}
	timeText := fmt.Sprintf("%s %d' ", half, minutes)
	rightSection := lipgloss.NewStyle().
		Foreground(lipgloss.Color("241")).
		Render(timeText)

	// Calculate widths for three sections
	sideWidth := (width - lipgloss.Width(centerSection)) / 2

	leftStyled := lipgloss.NewStyle().
		Width(sideWidth).
		Align(lipgloss.Left).
		Render(leftSection)

	rightStyled := lipgloss.NewStyle().
		Width(sideWidth).
		Align(lipgloss.Right).
		Render(rightSection)

	// Join all sections
	headerContent := lipgloss.JoinHorizontal(
		lipgloss.Top,
		leftStyled,
		centerSection,
		rightStyled,
	)

	// Wrap in a bordered style
	header := lipgloss.NewStyle().
		Width(width).
		Border(lipgloss.NormalBorder(), false, false, true, false).
		BorderForeground(lipgloss.Color("240")).
		Render(headerContent)

	return header
}
