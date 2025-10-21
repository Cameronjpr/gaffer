package tui

import (
	"fmt"
	"time"

	"github.com/cameronjpr/gaffer/internal/game"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type MatchModel struct {
	match    *game.Match
	isPaused bool
	width    int
	height   int
}

func NewMatchModel(match game.Match) MatchModel {
	return MatchModel{
		match:    &match,
		isPaused: false,
	}
}

func (m MatchModel) Init() tea.Cmd {
	return nil
}

func (m MatchModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {

		case tea.KeySpace:
			m.isPaused = !m.isPaused
			if m.isPaused {
				return m, nil
			}

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
	}
	return m, nil
}

func (m MatchModel) View() string {
	// Header
	timeStr := fmt.Sprintf("(%v:00)", m.match.CurrentPhase)
	scoreStr := fmt.Sprintf("%v [%v] – [%v] %v", m.match.Home.Club.Name, m.match.Home.Score, m.match.Away.Score, m.match.Away.Club.Name)

	time := lipgloss.NewStyle().
		Background(lipgloss.Color("#0000ff")).
		Padding(0, 1).
		Render(timeStr)

	score := lipgloss.NewStyle().
		Background(lipgloss.Color("#0000ff")).
		Padding(0, 1).
		Width(m.width - lipgloss.Width(time)).
		Align(lipgloss.Right).
		Render(scoreStr)

	header := lipgloss.JoinHorizontal(lipgloss.Left, time, score)

	// Footer
	footer := ""
	if m.match.IsHalfTime {
		footer = lipgloss.NewStyle().
			Align(lipgloss.Center).
			Width(m.width).
			Render(fmt.Sprintf("Half time – %d - %d", m.match.Home.Score, m.match.Away.Score))
	} else if len(m.match.Commentary) > 0 {
		latestCommentaryMsg :=
			m.match.Commentary[len(m.match.Commentary)-1]
		style :=
			lipgloss.NewStyle().Align(lipgloss.Center).Width(m.width)

		if latestCommentaryMsg.Flash {
			style =
				style.Bold(true).Background(lipgloss.Color("#ff0000"))
		} else {
			style = style.Background(lipgloss.Color("#000000"))
		}

		footer = style.Render(latestCommentaryMsg.Message)
	}

	// Main content
	var matchInfoText string
	if m.match.IsComplete {
		matchInfoText = fmt.Sprintf("Match over\n\nFinal score: %d - %d", m.match.Home.Score, m.match.Away.Score)
	} else if m.match.IsHalfTime {
		matchInfoText = fmt.Sprintf("Half time – %d - %d", m.match.Home.Score, m.match.Away.Score)
	} else {
		matchInfoText = fmt.Sprintf("Running simulation – Half %v", m.match.CurrentHalf)
	}

	// Calculate content height (total height - header - footer)
	// Header and footer each take 1 line
	contentHeight := m.height - 2
	if contentHeight < 0 {
		contentHeight = 0
	}

	content := lipgloss.Place(m.width, contentHeight, lipgloss.Center, lipgloss.Center, matchInfoText)

	return lipgloss.JoinVertical(lipgloss.Top, header, content, footer)
}
