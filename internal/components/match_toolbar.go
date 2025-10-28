package components

import (
	"github.com/cameronjpr/gaffer/internal/domain"
	"github.com/charmbracelet/lipgloss"
)

func MatchToolbar(width int, m *domain.Match) string {
	itemWidth := min(width/3, 20)

	return Centered(
		width,
		4,
		lipgloss.JoinHorizontal(lipgloss.Center,
			lipgloss.NewStyle().Padding(2).Width(itemWidth).Render("[S]ubstition"),
			lipgloss.NewStyle().Padding(2).Width(itemWidth).Render("[T]actics"),
		),
	)
}
