package components

import (
	"fmt"

	"github.com/cameronjpr/gaffer/internal/domain"
	"github.com/charmbracelet/lipgloss"
)

func TeamSheet(participant *domain.MatchParticipant) string {
	return lipgloss.JoinVertical(
		lipgloss.Left,
		lipgloss.NewStyle().
			Bold(true).
			Padding(1).
			Background(lipgloss.Color(participant.Club.Background)).
			Foreground(lipgloss.Color(participant.Club.Foreground)).
			Render(participant.Club.Name),
		"\n",

		lipgloss.NewStyle().Italic(true).Render(participant.Formation),
		lipgloss.NewStyle().Bold(true).Render(fmt.Sprintf("%.2f avg.", participant.GetAverageQuality())),
		"",
		participant.GetLineup(nil),
	)
}
