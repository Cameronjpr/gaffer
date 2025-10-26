package tui

import (
	"github.com/cameronjpr/gaffer/internal/db"
	"github.com/cameronjpr/gaffer/internal/domain"
	"github.com/cameronjpr/gaffer/internal/repository"
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
	clubRepo     domain.ClubRepository
	fixtureRepo  domain.FixtureRepository
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
	// Create repositories
	clubRepo := repository.NewClubRepository(queries)
	fixtureRepo := repository.NewFixtureRepository(queries, clubRepo)

	// Get all clubs with players from repository
	clubs, err := clubRepo.GetAll()
	if err != nil {
		panic(err)
	}

	season := domain.NewSeason(clubs)
	season.GenerateAllFixtures()

	// Get the first fixture and create a match from it
	nextFixture, err := season.GetNextFixture()
	if err != nil {
		// In a real app, you'd handle this gracefully
		panic(err)
	}
	currentMatch := domain.NewMatchFromFixture(nextFixture)

	return &AppModel{
		clubRepo:     clubRepo,
		fixtureRepo:  fixtureRepo,
		mode:         MenuMode,
		season:       season,
		currentMatch: currentMatch,
		menu: NewMenuModel([]list.Item{
			item("New game"),
			item("Settings"),
		}),
		onboarding: NewOnboardingModel(season),
		managerHub: NewManagerHubModel(season, nil, nil),
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
