package domain

import (
	"testing"
)

func TestGetNextFixture_AdvancesThroughGameweeks(t *testing.T) {
	// Setup: Create a season with two clubs
	club1 := &Club{Name: "Arsenal"}
	club2 := &Club{Name: "Chelsea"}
	clubs := []*Club{club1, club2}

	season := NewSeason(clubs)
	season.GenerateGameweeks()

	// Get the first fixture (should be from Gameweek 1)
	firstFixture, err := season.GetNextFixture()
	if err != nil {
		t.Fatalf("Failed to get first fixture: %v", err)
	}

	if firstFixture == nil {
		t.Fatal("First fixture is nil")
	}

	// Verify it's from Gameweek 1
	if season.Gameweeks[0].Fixtures[0].ID != firstFixture.ID {
		t.Errorf("First fixture should be from Gameweek 1")
	}

	// Simulate playing the match by populating the Result
	firstFixture.Result = NewMatchFromFixture(firstFixture)

	// Get the next fixture (should be from Gameweek 1, fixture 2)
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

	// Verify it's the second fixture from Gameweek 1
	if season.Gameweeks[0].Fixtures[1].ID != secondFixture.ID {
		t.Errorf("Second fixture should be Gameweek 1, Fixture 2")
	}

	// Mark second fixture as played
	secondFixture.Result = NewMatchFromFixture(secondFixture)

	// Get the next fixture (should be from Gameweek 2)
	thirdFixture, err := season.GetNextFixture()
	if err != nil {
		t.Fatalf("Failed to get third fixture: %v", err)
	}

	if thirdFixture == nil {
		t.Fatal("Third fixture is nil")
	}

	// Verify it's from Gameweek 2
	if season.Gameweeks[1].Fixtures[0].ID != thirdFixture.ID {
		t.Errorf("Third fixture should be from Gameweek 2, got fixture from different gameweek")
	}
}

func TestGetNextFixture_ReturnsErrorWhenNoFixturesRemain(t *testing.T) {
	// Setup: Create a season with two clubs
	club1 := &Club{Name: "Arsenal"}
	club2 := &Club{Name: "Chelsea"}
	clubs := []*Club{club1, club2}

	season := NewSeason(clubs)
	season.GenerateGameweeks()

	// Mark all fixtures as played
	for i := range season.Gameweeks {
		for j := range season.Gameweeks[i].Fixtures {
			season.Gameweeks[i].Fixtures[j].Result = NewMatchFromFixture(&season.Gameweeks[i].Fixtures[j])
		}
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
	club1 := &Club{Name: "Arsenal"}
	club2 := &Club{Name: "Chelsea"}
	clubs := []*Club{club1, club2}

	season := NewSeason(clubs)
	season.GenerateGameweeks()

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
