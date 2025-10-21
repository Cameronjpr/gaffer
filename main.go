package main

import (
	"fmt"
	"os"

	"github.com/cameronjpr/gaffer/internal/tui"
	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	model := tui.NewModel()

	p := tea.NewProgram(model, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}
