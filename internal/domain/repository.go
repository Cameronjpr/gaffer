package domain

// ClubRepository handles persistence of clubs
type ClubRepository interface {
	GetAll() ([]*ClubWithPlayers, error)
	GetByName(name string) (*Club, error)
	GetByID(id int64) (*ClubWithPlayers, error)
}

// FixtureRepository handles persistence of fixtures
type FixtureRepository interface {
	GetAll() ([]*Fixture, error)
	GetByID(id int64) (*Fixture, error)
	GetByClubID(clubID int64) ([]*Fixture, error)
	GetByGameweek(gameweek int) ([]*Fixture, error)
}

// SeasonRepository handles persistence of seasons
type SeasonRepository interface {
	Save(season *Season) error
	Load(id int) (*Season, error)
}
