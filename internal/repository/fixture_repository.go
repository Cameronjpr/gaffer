package repository

import (
	"context"
	"fmt"

	"github.com/cameronjpr/gaffer/internal/db"
	"github.com/cameronjpr/gaffer/internal/domain"
)

type FixtureRepo struct {
	queries  *db.Queries
	clubRepo *ClubRepo
}

func NewFixtureRepository(queries *db.Queries, clubRepo *ClubRepo) *FixtureRepo {
	return &FixtureRepo{
		queries:  queries,
		clubRepo: clubRepo,
	}
}

// GetAll fetches all fixtures from the database
func (r *FixtureRepo) GetAll() ([]*domain.Fixture, error) {
	ctx := context.Background()

	dbFixtures, err := r.queries.GetAllFixtures(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get fixtures: %w", err)
	}

	fixtures := make([]*domain.Fixture, len(dbFixtures))
	for i, dbFixture := range dbFixtures {
		fixture, err := r.dbFixtureToDomain(dbFixture)
		if err != nil {
			return nil, fmt.Errorf("failed to convert fixture %d: %w", dbFixture.ID, err)
		}
		fixtures[i] = fixture
	}

	return fixtures, nil
}

// GetByID fetches a single fixture by ID
func (r *FixtureRepo) GetByID(id int64) (*domain.Fixture, error) {
	ctx := context.Background()

	dbFixture, err := r.queries.GetFixtureByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get fixture %d: %w", id, err)
	}

	return r.dbFixtureToDomain(dbFixture)
}

// GetByClubID fetches all fixtures for a specific club
func (r *FixtureRepo) GetByClubID(clubID int64) ([]*domain.Fixture, error) {
	ctx := context.Background()

	dbFixtures, err := r.queries.GetFixturesByClubID(ctx, db.GetFixturesByClubIDParams{
		HomeTeamID: clubID,
		AwayTeamID: clubID,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get fixtures for club %d: %w", clubID, err)
	}

	fixtures := make([]*domain.Fixture, len(dbFixtures))
	for i, dbFixture := range dbFixtures {
		fixture, err := r.dbFixtureToDomain(dbFixture)
		if err != nil {
			return nil, fmt.Errorf("failed to convert fixture %d: %w", dbFixture.ID, err)
		}
		fixtures[i] = fixture
	}

	return fixtures, nil
}

// GetByGameweek fetches all fixtures for a specific gameweek
func (r *FixtureRepo) GetByGameweek(gameweek int) ([]*domain.Fixture, error) {
	ctx := context.Background()

	// Get all fixtures and filter by gameweek
	// Note: This could be optimized with a dedicated SQL query
	dbFixtures, err := r.queries.GetAllFixtures(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get fixtures: %w", err)
	}

	fixtures := make([]*domain.Fixture, 0)
	for _, dbFixture := range dbFixtures {
		if int(dbFixture.Gameweek) == gameweek {
			fixture, err := r.dbFixtureToDomain(dbFixture)
			if err != nil {
				return nil, fmt.Errorf("failed to convert fixture %d: %w", dbFixture.ID, err)
			}
			fixtures = append(fixtures, fixture)
		}
	}

	return fixtures, nil
}

// dbFixtureToDomain converts a database fixture to a domain Fixture
func (r *FixtureRepo) dbFixtureToDomain(dbFixture db.Fixture) (*domain.Fixture, error) {
	// Load the home team with players
	homeTeam, err := r.clubRepo.GetByID(dbFixture.HomeTeamID)
	if err != nil {
		return nil, fmt.Errorf("failed to get home team %d: %w", dbFixture.HomeTeamID, err)
	}

	// Load the away team with players
	awayTeam, err := r.clubRepo.GetByID(dbFixture.AwayTeamID)
	if err != nil {
		return nil, fmt.Errorf("failed to get away team %d: %w", dbFixture.AwayTeamID, err)
	}

	return &domain.Fixture{
		ID:       int(dbFixture.ID),
		Gameweek: int(dbFixture.Gameweek),
		HomeTeam: homeTeam,
		AwayTeam: awayTeam,
		Result:   nil, // Match results would be loaded separately if needed
	}, nil
}

// Ensure FixtureRepo implements domain.FixtureRepository
var _ domain.FixtureRepository = (*FixtureRepo)(nil)
