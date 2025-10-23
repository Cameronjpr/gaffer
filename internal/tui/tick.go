package tui

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

const (
	defaultTickDuration   = time.Millisecond * 300
	goalTickDuration      = time.Millisecond * 3000
	addedTimeTickDuration = time.Millisecond * 1000
)

func tick() tea.Cmd {
	return tickWithDuration(defaultTickDuration)
}

func tickWithDuration(d time.Duration) tea.Cmd {
	return tea.Tick(d, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}
