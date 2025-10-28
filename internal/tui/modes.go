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
	gameStateRepo *repository.GameStateRepo
	clubRepo      domain.ClubRepository // address this
	fixtureRepo   domain.FixtureRepository
	matchRepo     *repository.MatchRepo
	mode          Mode
	clubs         []*domain.ClubWithPlayers
	fixtures      []*domain.Fixture
	currentMatch  *domain.Match
	menu          *MenuModel
	onboarding    *OnboardingModel
	managerHub    *ManagerHubModel
	prematch      *PreMatchModel
	match         *MatchModel
	width         int
	height        int
}

func NewModel(queries *db.Queries) *AppModel {
	// Create repositories
	gameStateRepo := repository.NewGameStateRepository(queries)
	clubRepo := repository.NewClubRepository(queries)
	fixtureRepo := repository.NewFixtureRepository(queries, clubRepo)
	matchRepo := repository.NewMatchRepository(queries)

	// Get all clubs with players from repository
	clubs, err := clubRepo.GetAll()
	if err != nil {
		panic(err)
	}

	// Get all fixtures
	fixtures, err := fixtureRepo.GetAll()
	if err != nil {
		panic(err)
	}

	return &AppModel{
		gameStateRepo: gameStateRepo,
		clubRepo:      clubRepo,
		fixtureRepo:   fixtureRepo,
		matchRepo:     matchRepo,
		mode:          MenuMode,
		clubs:         clubs,
		fixtures:      fixtures,
		currentMatch:  nil,
		menu: NewMenuModel([]list.Item{
			item("New game"),
			item("Settings"),
		}),
		onboarding: NewOnboardingModel(clubs),
		managerHub: NewManagerHubModel(nil, nil, nil),
		prematch:   NewPreMatchModel(nil),
		match:      NewMatchModel(nil),
		width:      0,
		height:     0,
	}
}

type goToOnboardingMsg struct{}

type goToManagerHubMsg struct {
	ClubID int64
}

type startPreMatchMsg struct{}

type startMatchMsg struct{}

type matchFinishedMsg struct {
	match *domain.Match
}
