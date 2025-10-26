package domain

import (
	"testing"
)

// Test helper to create a ClubWithPlayers for testing
func makeTestClubWithPlayers(name string) *ClubWithPlayers {
	club := &Club{Name: name}
	players := []Player{
		{Name: "Player 1", Quality: 15},
		{Name: "Player 2", Quality: 15},
		{Name: "Player 3", Quality: 15},
		{Name: "Player 4", Quality: 15},
		{Name: "Player 5", Quality: 15},
		{Name: "Player 6", Quality: 15},
		{Name: "Player 7", Quality: 15},
		{Name: "Player 8", Quality: 15},
		{Name: "Player 9", Quality: 15},
		{Name: "Player 10", Quality: 15},
		{Name: "Player 11", Quality: 15},
	}
	return &ClubWithPlayers{
		Club:    club,
		Players: players,
	}
}

func TestGetNextFixture_AdvancesThroughGameweeks(t *testing.T) {
	// Setup: Create a season with three clubs for more fixtures
	club1 := makeTestClubWithPlayers("Arsenal")
	club2 := makeTestClubWithPlayers("Chelsea")
	club3 := makeTestClubWithPlayers("Liverpool")
	clubs := []*ClubWithPlayers{club1, club2, club3}

	season := NewSeason(clubs)
	season.GenerateAllFixtures()

	// With 3 clubs: 3*2 = 6 total fixtures distributed across 38 gameweeks
	// First gameweek should have at least 1 fixture

	// Get the first fixture
	firstFixture, err := season.GetNextFixture()
	if err != nil {
		t.Fatalf("Failed to get first fixture: %v", err)
	}

	if firstFixture == nil {
		t.Fatal("First fixture is nil")
	}

	// Simulate playing the match by populating the Result
	firstFixture.Result = NewMatchFromFixture(firstFixture)

	// Get the next fixture
	secondFixture, err := season.GetNextFixture()
	if err != nil {
		t.Fatalf("Failed to get second fixture: %v", err)
	}

	if secondFixture == nil {
		t.Fatal("Second fixture is nil")
	}

	// Verify it's a different fixture
	if firstFixture.ID == secondFixture.ID {
		t.Errorf("Second fixture should be different from first fixture, but got same ID: %d", secondFixture.ID)
	}

	// Mark second fixture as played
	secondFixture.Result = NewMatchFromFixture(secondFixture)

	// Get the third fixture
	thirdFixture, err := season.GetNextFixture()
	if err != nil {
		t.Fatalf("Failed to get third fixture: %v", err)
	}

	if thirdFixture == nil {
		t.Fatal("Third fixture is nil")
	}

	// Verify it's a different fixture
	if thirdFixture.ID == firstFixture.ID || thirdFixture.ID == secondFixture.ID {
		t.Errorf("Third fixture should be different from first two fixtures")
	}
}

func TestGetFixturesForClub_ReturnsFixturesForChosenClub(t *testing.T) {
	// Setup: Create a season with two clubs
	club1 := makeTestClubWithPlayers("Arsenal")
	club2 := makeTestClubWithPlayers("Chelsea")
	clubs := []*ClubWithPlayers{club1, club2}

	season := NewSeason(clubs)
	season.GenerateAllFixtures()

	// Verify fixtures for chosen club
	// With 2 clubs across 38 gameweeks: each plays the other twice per gameweek (home and away)
	// = 2 fixtures per gameweek * 38 gameweeks = 76 total fixtures
	fixtures := season.GetFixturesForClub( club1.Club)
	expectedFixtures := 76 // 38 home + 38 away
	if len(fixtures) != expectedFixtures {
		t.Errorf("Expected %d fixtures for chosen club, got %d", expectedFixtures, len(fixtures))
	}

	for _, fixture := range fixtures {
		if fixture.HomeTeam.Club != club1.Club && fixture.AwayTeam.Club != club1.Club {
			t.Errorf("Expected fixture to be for chosen club, got fixture with home team %s and away team %s", fixture.HomeTeam.Club.Name, fixture.AwayTeam.Club.Name)
		}
	}

	// Verify equal distribution of home and away fixtures
	homeCount := 0
	awayCount := 0
	for _, fixture := range fixtures {
		if fixture.HomeTeam.Club == club1.Club {
			homeCount++
		}
		if fixture.AwayTeam.Club == club1.Club {
			awayCount++
		}
	}
	expectedHomeAway := 38
	if homeCount != expectedHomeAway || awayCount != expectedHomeAway {
		t.Errorf("Expected %d home and %d away fixtures, got %d home and %d away", expectedHomeAway, expectedHomeAway, homeCount, awayCount)
	}
}

func TestGetNextFixture_ReturnsErrorWhenNoFixturesRemain(t *testing.T) {
	// Setup: Create a season with two clubs
	club1 := makeTestClubWithPlayers("Arsenal")
	club2 := makeTestClubWithPlayers("Chelsea")
	clubs := []*ClubWithPlayers{club1, club2}

	season := NewSeason(clubs)
	season.GenerateAllFixtures()

	// Mark all fixtures as played
	for i := range season.Fixtures {
		season.Fixtures[i].Result = NewMatchFromFixture(season.Fixtures[i])
	}

	// Try to get next fixture - should return error
	fixture, err := season.GetNextFixture()
	if err == nil {
		t.Error("Expected error when no fixtures remain, but got nil")
	}

	if fixture != nil {
		t.Errorf("Expected nil fixture when season is complete, got fixture with ID: %d", fixture.ID)
	}

	expectedError := "no unplayed fixtures remaining in the season"
	if err.Error() != expectedError {
		t.Errorf("Expected error message '%s', got '%s'", expectedError, err.Error())
	}
}

func TestGetNextFixture_WithoutMarkingAsPlayed_ReturnsSameFixture(t *testing.T) {
	// This test verifies the BUG: if we don't mark fixtures as played,
	// GetNextFixture will keep returning the same fixture
	club1 := makeTestClubWithPlayers("Arsenal")
	club2 := makeTestClubWithPlayers("Chelsea")
	clubs := []*ClubWithPlayers{club1, club2}

	season := NewSeason(clubs)
	season.GenerateAllFixtures()

	// Get the first fixture
	firstFixture, err := season.GetNextFixture()
	if err != nil {
		t.Fatalf("Failed to get first fixture: %v", err)
	}

	// WITHOUT marking it as played, get the next fixture again
	nextFixture, err := season.GetNextFixture()
	if err != nil {
		t.Fatalf("Failed to get next fixture: %v", err)
	}

	// This is the BUG: we get the same fixture because we never marked it as played
	if firstFixture.ID != nextFixture.ID {
		t.Errorf("BUG NOT REPRODUCED: Expected same fixture ID, got first=%d, next=%d", firstFixture.ID, nextFixture.ID)
	} else {
		t.Logf("BUG CONFIRMED: GetNextFixture returns the same fixture (ID=%d) when not marked as played", firstFixture.ID)
	}
}

func TestGetLeagueTable_HandlesUnplayedFixtures(t *testing.T) {
	// This test ensures GetLeagueTable doesn't crash when there are unplayed fixtures
	// Regression test for nil pointer panic when fixture.Result is nil
	club1 := makeTestClubWithPlayers("Arsenal")
	club2 := makeTestClubWithPlayers("Chelsea")
	club3 := makeTestClubWithPlayers("Liverpool")
	clubs := []*ClubWithPlayers{club1, club2, club3}

	season := NewSeason(clubs)
	season.GenerateAllFixtures()

	// Play only the first fixture
	firstFixture, err := season.GetNextFixture()
	if err != nil {
		t.Fatalf("Failed to get first fixture: %v", err)
	}
	firstFixture.Result = NewMatchFromFixture(firstFixture)
	// Simulate a win for home team
	firstFixture.Result.Home.Score = 2
	firstFixture.Result.Away.Score = 1

	// Debug: Check fixture results
	allFixtures := season.GetFixturesForClub( club1.Club)
	playedCount := 0
	for _, f := range allFixtures {
		if f.Result != nil {
			playedCount++
			if f.Result.Home == nil || f.Result.Away == nil {
				t.Errorf("Fixture %d has Result but nil Home or Away", f.ID)
			}
		}
	}
	t.Logf("Total fixtures: %d, Played: %d", len(allFixtures), playedCount)

	// Call GetLeagueTable - should NOT crash even though most fixtures are unplayed
	table := season.GetLeagueTable()

	// Verify table was generated
	if len(table.Positions) != 3 {
		t.Errorf("Expected 3 positions in table, got %d", len(table.Positions))
	}

	// Verify the winner has 3 points, others have 0
	totalPoints := 0
	for _, pos := range table.Positions {
		totalPoints += pos.Points
		t.Logf("Club %s: %d points", pos.Club.Name, pos.Points)
		if pos.Club == firstFixture.Result.GetWinner() && pos.Points != 3 {
			t.Errorf("Winner %s should have 3 points, got %d", pos.Club.Name, pos.Points)
		}
	}

	if totalPoints != 3 {
		t.Errorf("Expected total of 3 points (1 match played), got %d", totalPoints)
	}
}
