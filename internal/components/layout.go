package components

import (
	"github.com/charmbracelet/lipgloss"
)

func ThreeColumnLayout(width int, c1, c2, c3 string) string {
	colWidth := width / 3

	return lipgloss.JoinHorizontal(
		lipgloss.Top,
		lipgloss.NewStyle().Padding(2).Width(colWidth).Render(c1),
		lipgloss.NewStyle().Padding(2).Width(colWidth).Render(c2),
		lipgloss.NewStyle().Padding(2).Width(colWidth).Render(c3),
	)
}
