package tui

import (
	"fmt"
	"time"

	"github.com/cameronjpr/gaffer/internal/game"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Model struct {
	match    *game.Match
	width    int
	height   int
	isPaused bool
}

type tickMsg time.Time

func NewModel(match *game.Match) Model {
	return Model{
		match:    match,
		width:    0,
		height:   0,
		isPaused: false,
	}
}

func tick() tea.Cmd {
	return tickWithDuration(time.Millisecond * 100)
}

func tickWithDuration(d time.Duration) tea.Cmd {
	return tea.Tick(d, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func (m Model) Init() tea.Cmd {
	return tick()
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		// Handle quit
		if msg.Type == tea.KeyCtrlC || msg.Type == tea.KeyEsc {
			return m, tea.Quit
		}

		if msg.Type == tea.KeySpace {
			m.isPaused = !m.isPaused
			return m, tick()
		}

	case tickMsg:
		// Check if match is complete
		if m.match.CurrentPhase >= 90 {
			m.match.IsComplete = true
			return m, nil
		}

		// If paused, don't advance
		if m.isPaused {
			return m, nil
		}

		// Auto-advance to next phase
		result := m.match.PlayPhase()

		// Update match state
		m.match.Home.Score += result.HomeGoals
		m.match.Away.Score += result.AwayGoals
		m.match.CurrentPhase++
		m.match.PhaseHistory = append(m.match.PhaseHistory, result)

		if m.match.CurrentPhase == 45 {
			m.isPaused = true
			m.match.CurrentHalf++
			m.match.IsHalfTime = true
			return m, nil
		}

		// Check if a goal was scored in this phase
		goalScored := result.HomeGoals > 0 || result.AwayGoals > 0

		// If a goal was scored, pause longer to let user see it
		if goalScored {
			return m, tickWithDuration(time.Second * 2)
		}

		// Return the next tick command to keep the loop going
		return m, tick()

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	}

	return m, nil
}

func (m Model) View() string {
	// Header
	time := lipgloss.NewStyle().
		Align(lipgloss.Left).
		Width(8).
		Background(lipgloss.Color("#0000ff")).
		Render(fmt.Sprintf("(%v:00)", m.match.CurrentPhase))
	score := lipgloss.NewStyle().
		Align(lipgloss.Right).
		Width(24).
		Background(lipgloss.Color("#0000ff")).
		Render(fmt.Sprintf("Home %d - %d Away", m.match.Home.Score, m.match.Away.Score))

	header := lipgloss.JoinHorizontal(lipgloss.Left, time, score)

	// Footer
	footer := ""
	if len(m.match.Commentary) > 0 {
		latestCommentaryMsg := m.match.Commentary[len(m.match.Commentary)-1]
		if latestCommentaryMsg.Flash {
			footer = lipgloss.NewStyle().
				Align(lipgloss.Center).
				Width(m.width).
				Background(lipgloss.Color("#ff0000")).
				Render(latestCommentaryMsg.Message)
		} else {
			footer = lipgloss.NewStyle().
				Align(lipgloss.Center).
				Width(m.width).
				Render(latestCommentaryMsg.Message)
		}
	}

	// Main content
	matchInProgress := lipgloss.NewStyle().
		Align(lipgloss.Center).
		Width(m.width).
		Render(fmt.Sprintf("Running simulation – Half %v", m.match.CurrentHalf))

	matchHalfTime := lipgloss.NewStyle().
		Align(lipgloss.Center).
		Width(m.width).
		Render(fmt.Sprintf("Half time – %d - %d", m.match.Home.Score, m.match.Away.Score))

	matchOver := lipgloss.NewStyle().
		Align(lipgloss.Center).
		Width(m.width).
		Render(fmt.Sprintf("Match over\n\nFinal score: %d - %d", m.match.Home.Score, m.match.Away.Score))

	matchInfo := ""

	if m.match.IsComplete {
		matchInfo = matchOver
	} else if m.match.IsHalfTime {
		matchInfo = matchHalfTime
	} else {
		matchInfo = matchInProgress
	}

	content := lipgloss.NewStyle().
		Width(m.width).
		// accommodate header and footer
		Height(m.height-1-1).
		Align(lipgloss.Center, lipgloss.Center).
		Render(matchInfo)

	return lipgloss.JoinVertical(lipgloss.Top, header, content, footer)
}
