package repository

import (
	"context"
	"fmt"
	"sort"

	"github.com/cameronjpr/gaffer/internal/db"
	"github.com/cameronjpr/gaffer/internal/domain"
)

type MatchRepo struct {
	queries *db.Queries
}

func NewMatchRepository(queries *db.Queries) *MatchRepo {
	return &MatchRepo{queries: queries}
}

// Create creates a new match record in the database
func (r *MatchRepo) Create(match *domain.Match) error {
	ctx := context.Background()

	_, err := r.queries.CreateMatch(ctx, db.CreateMatchParams{
		FixtureID:              int64(match.ForFixture.ID),
		CurrentMinute:          int64(match.CurrentMinute),
		CurrentHalf:            int64(match.CurrentHalf),
		HomeScore:              0,
		AwayScore:              0,
		ActiveZone:             int64(match.ActiveZone),
		HomeAttackingDirection: int64(match.HomeAttackingDirection),
	})
	if err != nil {
		return fmt.Errorf("failed to create match: %w", err)
	}

	return nil
}

// SaveResult persists the final match result to the database
func (r *MatchRepo) SaveResult(match *domain.Match) error {
	ctx := context.Background()

	homeScore, awayScore := match.GetScore()

	// Complete the match with final scores
	err := r.queries.CompleteMatch(ctx, db.CompleteMatchParams{
		HomeScore:  int64(homeScore),
		AwayScore:  int64(awayScore),
		FixtureID:  int64(match.ForFixture.ID),
	})
	if err != nil {
		return fmt.Errorf("failed to save match result: %w", err)
	}

	return nil
}

// GetByFixtureID retrieves a match by its fixture ID
func (r *MatchRepo) GetByFixtureID(fixtureID int64) (*db.Match, error) {
	ctx := context.Background()

	match, err := r.queries.GetMatchByFixtureID(ctx, fixtureID)
	if err != nil {
		return nil, fmt.Errorf("failed to get match for fixture %d: %w", fixtureID, err)
	}

	return &match, nil
}

// GetCompleted retrieves all completed matches
func (r *MatchRepo) GetCompleted() ([]db.Match, error) {
	ctx := context.Background()

	matches, err := r.queries.GetCompletedMatches(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get completed matches: %w", err)
	}

	return matches, nil
}

// IsFixturePlayed checks if a fixture has a completed match
func (r *MatchRepo) IsFixturePlayed(fixtureID int64) (bool, error) {
	match, err := r.GetByFixtureID(fixtureID)
	if err != nil {
		// If not found, it hasn't been played
		return false, nil
	}

	return match.IsCompleted == 1, nil
}

// GetAllCompletedResults returns a map of fixture ID to (homeScore, awayScore)
func (r *MatchRepo) GetAllCompletedResults() (map[int64][2]int, error) {
	ctx := context.Background()

	matches, err := r.queries.GetCompletedMatches(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get completed matches: %w", err)
	}

	results := make(map[int64][2]int)
	for _, match := range matches {
		results[match.FixtureID] = [2]int{int(match.HomeScore), int(match.AwayScore)}
	}

	return results, nil
}

// CalculateLeagueTable computes league standings from completed matches
func (r *MatchRepo) CalculateLeagueTable(clubs []*domain.ClubWithPlayers, fixtures []*domain.Fixture) (*domain.LeagueTable, error) {
	ctx := context.Background()

	// Get all completed matches
	completedMatches, err := r.queries.GetCompletedMatches(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get completed matches: %w", err)
	}

	// Build a map of fixture ID to match result
	matchResults := make(map[int64]db.Match)
	for _, match := range completedMatches {
		matchResults[match.FixtureID] = match
	}

	// Calculate standings for each club
	positions := make([]domain.LeaguePosition, 0, len(clubs))
	for _, clubWithPlayers := range clubs {
		club := clubWithPlayers.Club
		position := domain.LeaguePosition{
			Club: club,
		}

		// Find all fixtures for this club
		for _, fixture := range fixtures {
			isHome := fixture.HomeTeam.Club.ID == club.ID
			isAway := fixture.AwayTeam.Club.ID == club.ID

			if !isHome && !isAway {
				continue
			}

			// Check if this fixture has been played
			match, played := matchResults[int64(fixture.ID)]
			if !played {
				continue
			}

			position.Played++
			homeScore := int(match.HomeScore)
			awayScore := int(match.AwayScore)

			if isHome {
				position.GoalsFor += homeScore
				position.GoalsAgainst += awayScore

				if homeScore > awayScore {
					position.Won++
					position.Points += 3
				} else if homeScore == awayScore {
					position.Drawn++
					position.Points += 1
				} else {
					position.Lost++
				}
			} else {
				position.GoalsFor += awayScore
				position.GoalsAgainst += homeScore

				if awayScore > homeScore {
					position.Won++
					position.Points += 3
				} else if awayScore == homeScore {
					position.Drawn++
					position.Points += 1
				} else {
					position.Lost++
				}
			}
		}

		position.GoalDifference = position.GoalsFor - position.GoalsAgainst
		positions = append(positions, position)
	}

	// Sort the table
	sort.Sort(domain.ByLeagueStanding(positions))

	return &domain.LeagueTable{
		Positions: positions,
	}, nil
}
