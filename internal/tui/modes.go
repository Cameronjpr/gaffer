package tui

import (
	"github.com/cameronjpr/gaffer/internal/domain"
	"github.com/charmbracelet/bubbles/list"
)

type Mode int

const (
	MenuMode Mode = iota
	ManagerHubMode
	PreMatchMode
	MatchMode
)

type AppModel struct {
	mode         Mode
	season       *domain.Season
	currentMatch *domain.Match
	menu         MenuModel
	managerHub   ManagerHubModel
	prematch     PreMatchModel
	match        MatchModel
	width        int
	height       int
}

func NewModel() AppModel {
	// For now, hardcode Manchester City vs Arsenal until manager mode is implemented
	mci := domain.GetClubByName("Manchester City")
	ars := domain.GetClubByName("Arsenal")
	season := domain.NewSeason([]*domain.Club{mci, ars})
	season.GenerateGameweeks()

	// Get the first fixture and create a match from it
	nextFixture, err := season.GetNextFixture()
	if err != nil {
		// In a real app, you'd handle this gracefully
		panic(err)
	}
	currentMatch := domain.NewMatchFromFixture(nextFixture)

	return AppModel{
		mode:         MenuMode,
		season:       season,
		currentMatch: currentMatch,
		menu: NewMenuModel([]list.Item{
			item("New game"),
			item("Settings"),
		}),
		managerHub: NewManagerHubModel(season),
		prematch:   NewPreMatchModel(currentMatch),
		match:      NewMatchModel(currentMatch),
		width:      0,
		height:     0,
	}
}

type goToManagerHubMsg struct{}

type startPreMatchMsg struct{}

type startMatchMsg struct{}

type matchFinishedMsg struct {
	match *domain.Match
}
