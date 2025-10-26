package repository

import (
	"context"
	"fmt"

	"github.com/cameronjpr/gaffer/internal/db"
	"github.com/cameronjpr/gaffer/internal/domain"
)

type ClubRepo struct {
	queries *db.Queries
}

func NewClubRepository(queries *db.Queries) *ClubRepo {
	return &ClubRepo{queries: queries}
}

// GetAll fetches all clubs with their players from the database
func (r *ClubRepo) GetAll() ([]*domain.ClubWithPlayers, error) {
	ctx := context.Background()

	dbClubs, err := r.queries.GetAllClubs(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get clubs: %w", err)
	}

	clubsWithPlayers := make([]*domain.ClubWithPlayers, len(dbClubs))
	for i, dbClub := range dbClubs {
		dbPlayers, err := r.queries.GetPlayersByClubID(ctx, dbClub.ID)
		if err != nil {
			return nil, fmt.Errorf("failed to get players for club %s: %w", dbClub.Name, err)
		}

		players := make([]domain.Player, len(dbPlayers))
		for j, p := range dbPlayers {
			players[j] = domain.Player{
				Name:    p.Name,
				Quality: int(p.Quality),
			}
		}

		clubsWithPlayers[i] = &domain.ClubWithPlayers{
			Club:    dbClubToDomain(dbClub),
			Players: players,
		}
	}

	return clubsWithPlayers, nil
}

// GetByName returns a club by name from the database
func (r *ClubRepo) GetByName(name string) (*domain.Club, error) {
	ctx := context.Background()

	dbClub, err := r.queries.GetClubByName(ctx, name)
	if err != nil {
		return nil, fmt.Errorf("failed to get club %s: %w", name, err)
	}

	return dbClubToDomain(dbClub), nil
}

// GetByID returns a club with its players from the database
func (r *ClubRepo) GetByID(clubID int64) (*domain.ClubWithPlayers, error) {
	ctx := context.Background()

	dbClub, err := r.queries.GetClubByID(ctx, clubID)
	if err != nil {
		return nil, fmt.Errorf("failed to get club: %w", err)
	}

	dbPlayers, err := r.queries.GetPlayersByClubID(ctx, clubID)
	if err != nil {
		return nil, fmt.Errorf("failed to get players: %w", err)
	}

	players := make([]domain.Player, len(dbPlayers))
	for i, p := range dbPlayers {
		players[i] = domain.Player{
			Name:    p.Name,
			Quality: int(p.Quality),
		}
	}

	return &domain.ClubWithPlayers{
		Club:    dbClubToDomain(dbClub),
		Players: players,
	}, nil
}

// dbClubToDomain converts a database club to a domain Club
func dbClubToDomain(dbClub db.Club) *domain.Club {
	return &domain.Club{
		ID:         dbClub.ID,
		Name:       dbClub.Name,
		Strength:   int(dbClub.Strength),
		Background: dbClub.BackgroundColor,
		Foreground: dbClub.ForegroundColor,
	}
}

// Ensure ClubRepo implements domain.ClubRepository
var _ domain.ClubRepository = (*ClubRepo)(nil)
