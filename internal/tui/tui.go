package tui

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

type tickMsg time.Time

func (m AppModel) Init() tea.Cmd {
	return tick()
}

func (m AppModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
		// WindowSizeMsg will be passed to child models in the switch below

	case tea.KeyMsg:
		// Handle quit
		if msg.Type == tea.KeyCtrlC || msg.Type == tea.KeyEsc {
			return m, tea.Quit
		}

	case startPreMatchMsg:
		m.mode = PreMatchMode
		// Send WindowSizeMsg to newly activated model
		m.prematch.width = m.width
		m.prematch.height = m.height
		return m, tick()

	case startMatchMsg:
		m.mode = MatchMode
		// Send WindowSizeMsg to newly activated model
		m.match.width = m.width
		m.match.height = m.height
		return m, tick()
	}

	var cmd tea.Cmd
	switch m.mode {
	case MenuMode:
		var newMenu tea.Model
		newMenu, cmd = m.menu.Update(msg)
		m.menu = newMenu.(MenuModel)

	case PreMatchMode:
		var newPrematch tea.Model
		newPrematch, cmd = m.prematch.Update(msg)
		m.prematch = newPrematch.(PreMatchModel)

	case MatchMode:
		var newMatch tea.Model
		newMatch, cmd = m.match.Update(msg)
		m.match = newMatch.(MatchModel)

		// Check if match is complete
		if m.match.match.IsFullTime() {
			// m.mode = ResultsMode
		}

	}

	return m, cmd
}

func (m AppModel) View() string {
	switch m.mode {
	case MenuMode:
		return m.menu.View()
	case PreMatchMode:
		return m.prematch.View()
	case MatchMode:
		return m.match.View()
	}
	return "No mode"
}
