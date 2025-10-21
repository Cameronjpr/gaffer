package tui

import (
	"time"

	"github.com/cameronjpr/gaffer/internal/game"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type Mode int

const (
	MenuMode Mode = iota
	PreMatchMode
	MatchMode
)

type AppModel struct {
	mode     Mode
	menu     MenuModel
	prematch PreMatchModel
	match    MatchModel
	width    int
	height   int
}

func NewModel() AppModel {
	// For now, hardcode Leeds vs Arsenal until manager mode is implemented
	homeClub := game.GetClubByName("Manchester City")
	awayClub := game.GetClubByName("Arsenal")
	match := game.NewMatch(homeClub, awayClub)

	return AppModel{
		mode: MenuMode,
		menu: NewMenuModel([]list.Item{
			item("New game"),
			item("Settings"),
		}),
		prematch: NewPreMatchModel(match),
		match:    NewMatchModel(match),
		width:    0,
		height:   0,
	}
}

type startPreMatchMsg struct{}
type startMatchMsg struct{}

func tick() tea.Cmd {
	return tickWithDuration(time.Millisecond * 200)
}

func tickWithDuration(d time.Duration) tea.Cmd {
	return tea.Tick(d, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}
