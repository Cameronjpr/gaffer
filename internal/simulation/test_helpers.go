package simulation

import (
	"context"
	"database/sql"
	"testing"

	"github.com/cameronjpr/gaffer/internal/db"
	"github.com/cameronjpr/gaffer/internal/domain"
	"github.com/cameronjpr/gaffer/internal/repository"
)

// setupTestDB creates an in-memory SQLite database for testing
func setupTestDB(t *testing.T) (*sql.DB, *db.Queries) {
	t.Helper()

	// Create in-memory database
	database, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatalf("failed to open test database: %v", err)
	}

	// Run migrations
	schema := `
		CREATE TABLE IF NOT EXISTS clubs (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL UNIQUE,
			strength INTEGER NOT NULL CHECK(strength >= 0 AND strength <= 20),
			background_color TEXT NOT NULL,
			foreground_color TEXT NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		);

		CREATE TABLE IF NOT EXISTS players (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			club_id INTEGER NOT NULL,
			name TEXT NOT NULL,
			quality INTEGER NOT NULL CHECK(quality >= 0 AND quality <= 20),
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (club_id) REFERENCES clubs(id) ON DELETE CASCADE
		);
	`

	if _, err := database.Exec(schema); err != nil {
		t.Fatalf("failed to run migrations: %v", err)
	}

	queries := db.New(database)
	ctx := context.Background()

	// Seed test clubs
	arsenal, err := queries.CreateClub(ctx, db.CreateClubParams{
		Name:            "Arsenal",
		Strength:        20,
		BackgroundColor: "#DB0007",
		ForegroundColor: "#FFFFFF",
	})
	if err != nil {
		t.Fatalf("failed to create Arsenal: %v", err)
	}

	city, err := queries.CreateClub(ctx, db.CreateClubParams{
		Name:            "Manchester City",
		Strength:        19,
		BackgroundColor: "#6CABDD",
		ForegroundColor: "#1C2C5B",
	})
	if err != nil {
		t.Fatalf("failed to create Manchester City: %v", err)
	}

	// Add players for Arsenal
	arsenalPlayers := []struct {
		name    string
		quality int64
	}{
		{"Raya", 18},
		{"Timber", 17},
		{"Saliba", 18},
		{"Gabriel", 18},
		{"Calafiori", 17},
		{"Zubimendi", 18},
		{"Rice", 19},
		{"Ødegaard", 18},
		{"Saka", 19},
		{"Gyokeres", 17},
		{"Trossard", 17},
	}

	for _, p := range arsenalPlayers {
		_, err := queries.CreatePlayer(ctx, db.CreatePlayerParams{
			ClubID:  arsenal.ID,
			Name:    p.name,
			Quality: p.quality,
		})
		if err != nil {
			t.Fatalf("failed to create player %s: %v", p.name, err)
		}
	}

	// Add players for Manchester City
	cityPlayers := []struct {
		name    string
		quality int64
	}{
		{"Donnarumma", 18},
		{"Lewis", 15},
		{"Stones", 17},
		{"Ruben Dias", 18},
		{"Gvardiol", 17},
		{"González", 18},
		{"M. Nunes", 17},
		{"B. Silva", 15},
		{"Savinho", 18},
		{"Haaland", 19},
		{"Doku", 16},
	}

	for _, p := range cityPlayers {
		_, err := queries.CreatePlayer(ctx, db.CreatePlayerParams{
			ClubID:  city.ID,
			Name:    p.name,
			Quality: p.quality,
		})
		if err != nil {
			t.Fatalf("failed to create player %s: %v", p.name, err)
		}
	}

	return database, queries
}

// getTestClubs returns the two test clubs for testing
func getTestClubs(t *testing.T, queries *db.Queries) (*domain.ClubWithPlayers, *domain.ClubWithPlayers) {
	t.Helper()

	repo := repository.NewClubRepository(queries)

	arsenal, err := repo.GetByID(1) // Arsenal is ID 1
	if err != nil {
		t.Fatalf("failed to get Arsenal: %v", err)
	}

	city, err := repo.GetByID(2) // Manchester City is ID 2
	if err != nil {
		t.Fatalf("failed to get Manchester City: %v", err)
	}

	return arsenal, city
}
