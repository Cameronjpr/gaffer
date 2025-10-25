package tui

import (
	"github.com/cameronjpr/gaffer/internal/db"
	"github.com/cameronjpr/gaffer/internal/domain"
	"github.com/charmbracelet/bubbles/list"
)

type Mode int

const (
	MenuMode Mode = iota
	OnboardingMode
	ManagerHubMode
	PreMatchMode
	MatchMode
)

type AppModel struct {
	queries      *db.Queries
	mode         Mode
	season       *domain.Season
	currentMatch *domain.Match
	menu         *MenuModel
	onboarding   *OnboardingModel
	managerHub   *ManagerHubModel
	prematch     *PreMatchModel
	match        *MatchModel
	width        int
	height       int
}

func NewModel(queries *db.Queries) *AppModel {
	// Get all clubs with players from database
	clubs, err := domain.GetAllClubsWithPlayers(queries)
	if err != nil {
		panic(err)
	}

	season := domain.NewSeason(clubs)
	season.GenerateGameweeks()

	// Get the first fixture and create a match from it
	nextFixture, err := season.GetNextFixture()
	if err != nil {
		// In a real app, you'd handle this gracefully
		panic(err)
	}
	currentMatch := domain.NewMatchFromFixture(nextFixture)

	return &AppModel{
		queries:      queries,
		mode:         MenuMode,
		season:       season,
		currentMatch: currentMatch,
		menu: NewMenuModel([]list.Item{
			item("New game"),
			item("Settings"),
		}),
		onboarding: NewOnboardingModel(season),
		managerHub: NewManagerHubModel(season, nil),
		prematch:   NewPreMatchModel(currentMatch),
		match:      NewMatchModel(currentMatch),
		width:      0,
		height:     0,
	}
}

type goToOnboardingMsg struct{}

type goToManagerHubMsg struct {
	ClubName string
}

type startPreMatchMsg struct{}

type startMatchMsg struct{}

type matchFinishedMsg struct {
	match *domain.Match
}
