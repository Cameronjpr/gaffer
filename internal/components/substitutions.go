package components

import (
	"github.com/cameronjpr/gaffer/internal/domain"
	"github.com/charmbracelet/lipgloss"
)

func SubstitutionsModal(team []*domain.MatchPlayerParticipant, bench []*domain.MatchPlayerParticipant) string {
	teamStr := "\nTeam:\n"

	for _, participant := range team {
		teamStr += participant.Player.Name + "\n"
	}

	benchStr := "\nBench:\n"
	for _, participant := range bench {
		benchStr += participant.Player.Name + "\n"
	}

	return lipgloss.JoinHorizontal(
		lipgloss.Center,
		lipgloss.NewStyle().Margin(0, 2).Render(teamStr),
		lipgloss.NewStyle().Margin(0, 2).Render(benchStr),
	)
}
