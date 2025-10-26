package components

import (
	"fmt"

	"github.com/cameronjpr/gaffer/internal/domain"
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/lipgloss"
)

var baseStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("240"))

func Table(lt domain.LeagueTable) string {
	columns := []table.Column{
		{Title: "Pos.", Width: 4},
		{Title: "Club", Width: 16},
		{Title: "Played", Width: 6},
		{Title: "Points", Width: 6},
	}

	rows := []table.Row{}

	for i, position := range lt.Positions {
		rows = append(rows, table.Row{
			fmt.Sprintf("%d", i+1),
			position.Club.Name,
			fmt.Sprintf("%d", position.Played),
			fmt.Sprintf("%d", position.Points),
		})
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(7),
	)

	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(false)
	t.SetStyles(s)

	return t.View()
}
