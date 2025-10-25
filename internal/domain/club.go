package domain

import (
	"context"
	"fmt"

	"github.com/cameronjpr/gaffer/internal/db"
)

// Club represents a football club with permanent attributes
type Club struct {
	ID         int64
	Name       string
	Strength   int // out of 20
	Background string
	Foreground string
}

// ClubWithPlayers is a view model for when you need club + players together
type ClubWithPlayers struct {
	Club    *Club
	Players []Player
}

// dbClubToDomain converts a database club to a domain Club
func dbClubToDomain(dbClub db.Club) *Club {
	return &Club{
		ID:         dbClub.ID,
		Name:       dbClub.Name,
		Strength:   int(dbClub.Strength),
		Background: dbClub.BackgroundColor,
		Foreground: dbClub.ForegroundColor,
	}
}

// GetAllClubs fetches all clubs from the database
func GetAllClubs(queries *db.Queries) ([]*Club, error) {
	ctx := context.Background()

	dbClubs, err := queries.GetAllClubs(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get clubs: %w", err)
	}

	clubs := make([]*Club, len(dbClubs))
	for i, dbClub := range dbClubs {
		clubs[i] = dbClubToDomain(dbClub)
	}

	return clubs, nil
}

// GetAllClubsWithPlayers fetches all clubs with their players from the database
func GetAllClubsWithPlayers(queries *db.Queries) ([]*ClubWithPlayers, error) {
	ctx := context.Background()

	dbClubs, err := queries.GetAllClubs(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get clubs: %w", err)
	}

	clubsWithPlayers := make([]*ClubWithPlayers, len(dbClubs))
	for i, dbClub := range dbClubs {
		dbPlayers, err := queries.GetPlayersByClubID(ctx, dbClub.ID)
		if err != nil {
			return nil, fmt.Errorf("failed to get players for club %s: %w", dbClub.Name, err)
		}

		players := make([]Player, len(dbPlayers))
		for j, p := range dbPlayers {
			players[j] = Player{
				Name:    p.Name,
				Quality: int(p.Quality),
			}
		}

		clubsWithPlayers[i] = &ClubWithPlayers{
			Club:    dbClubToDomain(dbClub),
			Players: players,
		}
	}

	return clubsWithPlayers, nil
}

// GetClubByName returns a club by name from the database
func GetClubByName(queries *db.Queries, name string) (*Club, error) {
	ctx := context.Background()

	dbClub, err := queries.GetClubByName(ctx, name)
	if err != nil {
		return nil, fmt.Errorf("failed to get club %s: %w", name, err)
	}

	return dbClubToDomain(dbClub), nil
}

// GetClubWithPlayers returns a club with its players from the database
func GetClubWithPlayers(queries *db.Queries, clubID int64) (*ClubWithPlayers, error) {
	ctx := context.Background()

	dbClub, err := queries.GetClubByID(ctx, clubID)
	if err != nil {
		return nil, fmt.Errorf("failed to get club: %w", err)
	}

	dbPlayers, err := queries.GetPlayersByClubID(ctx, clubID)
	if err != nil {
		return nil, fmt.Errorf("failed to get players: %w", err)
	}

	players := make([]Player, len(dbPlayers))
	for i, p := range dbPlayers {
		players[i] = Player{
			Name:    p.Name,
			Quality: int(p.Quality),
		}
	}

	return &ClubWithPlayers{
		Club:    dbClubToDomain(dbClub),
		Players: players,
	}, nil
}

// GetSquad returns a formatted string of the club's squad
func (c *Club) GetSquad(queries *db.Queries) (string, error) {
	ctx := context.Background()

	dbPlayers, err := queries.GetPlayersByClubID(ctx, c.ID)
	if err != nil {
		return "", fmt.Errorf("failed to get players: %w", err)
	}

	lineup := ""
	for _, player := range dbPlayers {
		lineup += fmt.Sprintf("%s (Q:%d)\n", player.Name, player.Quality)
	}
	return lineup, nil
}

// GetSquad returns a formatted string of the squad (for ClubWithPlayers)
func (cwp *ClubWithPlayers) GetSquad() string {
	lineup := ""
	for _, player := range cwp.Players {
		lineup += fmt.Sprintf("%s (Q:%d)\n", player.Name, player.Quality)
	}
	return lineup
}
