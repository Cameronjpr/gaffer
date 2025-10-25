package db

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
)

// ClubSeed represents the JSON structure for seeding clubs
type ClubSeed struct {
	Name       string        `json:"Name"`
	Strength   int64         `json:"Strength"`
	Background string        `json:"Background"`
	Foreground string        `json:"Foreground"`
	Players    []PlayerSeed  `json:"Players"`
}

// PlayerSeed represents the JSON structure for seeding players
type PlayerSeed struct {
	Name    string `json:"Name"`
	Quality int64  `json:"Quality"`
}

// SeedDatabase loads clubs and players from clubs.json into the database
func SeedDatabase(db *sql.DB, clubsJSONPath string) error {
	// Read the JSON file
	data, err := os.ReadFile(clubsJSONPath)
	if err != nil {
		return fmt.Errorf("failed to read clubs.json: %w", err)
	}

	// Parse JSON
	var clubs []ClubSeed
	if err := json.Unmarshal(data, &clubs); err != nil {
		return fmt.Errorf("failed to parse clubs.json: %w", err)
	}

	// Create queries instance
	queries := New(db)
	ctx := context.Background()

	// Check if database is already seeded
	existingClubs, err := queries.GetAllClubs(ctx)
	if err != nil {
		return fmt.Errorf("failed to check existing clubs: %w", err)
	}
	if len(existingClubs) > 0 {
		// Database already seeded
		return nil
	}

	// Seed clubs and players
	for _, clubSeed := range clubs {
		// Create club
		club, err := queries.CreateClub(ctx, CreateClubParams{
			Name:            clubSeed.Name,
			Strength:        clubSeed.Strength,
			BackgroundColor: clubSeed.Background,
			ForegroundColor: clubSeed.Foreground,
		})
		if err != nil {
			return fmt.Errorf("failed to create club %s: %w", clubSeed.Name, err)
		}

		// Create players for this club
		for _, playerSeed := range clubSeed.Players {
			_, err := queries.CreatePlayer(ctx, CreatePlayerParams{
				ClubID:  club.ID,
				Name:    playerSeed.Name,
				Quality: playerSeed.Quality,
			})
			if err != nil {
				return fmt.Errorf("failed to create player %s for club %s: %w", playerSeed.Name, clubSeed.Name, err)
			}
		}
	}

	return nil
}
