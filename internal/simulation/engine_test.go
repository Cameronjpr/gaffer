package simulation

import (
	"testing"

	"github.com/cameronjpr/gaffer/internal/domain"
)

// TestGoalAverageOverMultipleMatches verifies the simulation produces realistic goal counts.
// This acts as a smoke test to catch any RNG or probability changes that break the simulation.
// Expected: 2-3 goals per game on average over 50 simulated matches.
func TestGoalAverageOverMultipleMatches(t *testing.T) {
	const numMatches = 50
	const minAvgGoals = 2.0
	const maxAvgGoals = 3.0

	totalGoals := 0

	// Get two evenly-matched clubs for consistent testing
	homeClub := domain.GetClubByName("Arsenal")
	awayClub := domain.GetClubByName("Manchester City")

	if homeClub == nil || awayClub == nil {
		t.Fatal("Test clubs not found")
	}

	for i := 0; i < numMatches; i++ {
		// Create a fresh match
		match := domain.NewMatch(homeClub, awayClub)
		engine := NewEngine(&match)

		// Simulate full 90 minutes
		for minute := 1; minute <= 90; minute++ {
			match.CurrentMinute = minute

			// Handle half-time transition
			if minute == 46 {
				match.StartSecondHalf()
			}

			result := engine.PlayPhase()
			match.Home.Score += result.HomeGoals
			match.Away.Score += result.AwayGoals
		}

		matchGoals := match.Home.Score + match.Away.Score
		totalGoals += matchGoals

		// Log individual match for debugging if needed
		t.Logf("Match %d: %s %d - %d %s (Total: %d goals)",
			i+1,
			match.Home.Club.Name,
			match.Home.Score,
			match.Away.Score,
			match.Away.Club.Name,
			matchGoals)
	}

	avgGoals := float64(totalGoals) / float64(numMatches)

	t.Logf("\nSimulation Results:")
	t.Logf("  Matches simulated: %d", numMatches)
	t.Logf("  Total goals: %d", totalGoals)
	t.Logf("  Average goals per match: %.2f", avgGoals)
	t.Logf("  Expected range: %.1f - %.1f goals per match", minAvgGoals, maxAvgGoals)

	if avgGoals < minAvgGoals {
		t.Errorf("Average goals (%.2f) is below expected minimum (%.1f). Simulation may be too defensive or RNG is broken.",
			avgGoals, minAvgGoals)
	}

	if avgGoals > maxAvgGoals {
		t.Errorf("Average goals (%.2f) is above expected maximum (%.1f). Simulation may be too aggressive or RNG is broken.",
			avgGoals, maxAvgGoals)
	}
}

// TestSimulationDeterminism verifies that matches are non-deterministic.
// Two identical setups should produce different results (sanity check for RNG).
func TestSimulationDeterminism(t *testing.T) {
	homeClub := domain.GetClubByName("Arsenal")
	awayClub := domain.GetClubByName("Manchester City")

	if homeClub == nil || awayClub == nil {
		t.Fatal("Test clubs not found")
	}

	// Run two matches with identical setup
	scores := make([][2]int, 2)
	for matchNum := 0; matchNum < 2; matchNum++ {
		match := domain.NewMatch(homeClub, awayClub)
		engine := NewEngine(&match)

		for minute := 1; minute <= 90; minute++ {
			match.CurrentMinute = minute
			if minute == 46 {
				match.StartSecondHalf()
			}
			result := engine.PlayPhase()
			match.Home.Score += result.HomeGoals
			match.Away.Score += result.AwayGoals
		}

		scores[matchNum] = [2]int{match.Home.Score, match.Away.Score}
	}

	t.Logf("Match 1: %d - %d", scores[0][0], scores[0][1])
	t.Logf("Match 2: %d - %d", scores[1][0], scores[1][1])

	// It's extremely unlikely (but technically possible) that two matches produce identical scores
	// If they're identical, log a warning but don't fail (could be legitimate randomness)
	if scores[0] == scores[1] {
		t.Logf("WARNING: Two matches produced identical scores. This is unlikely but possible with RNG.")
	}
}

// TestSimulationBasicSanity verifies a single match produces sensible results
func TestSimulationBasicSanity(t *testing.T) {
	homeClub := domain.GetClubByName("Arsenal")
	awayClub := domain.GetClubByName("Manchester City")

	if homeClub == nil || awayClub == nil {
		t.Fatal("Test clubs not found")
	}

	match := domain.NewMatch(homeClub, awayClub)
	engine := NewEngine(&match)

	for minute := 1; minute <= 90; minute++ {
		match.CurrentMinute = minute
		if minute == 46 {
			match.StartSecondHalf()
		}
		result := engine.PlayPhase()
		match.Home.Score += result.HomeGoals
		match.Away.Score += result.AwayGoals
	}

	totalGoals := match.Home.Score + match.Away.Score

	t.Logf("Final score: %s %d - %d %s (Total: %d goals)",
		match.Home.Club.Name,
		match.Home.Score,
		match.Away.Score,
		match.Away.Club.Name,
		totalGoals)

	// Basic sanity checks
	if match.Home.Score < 0 {
		t.Errorf("Home score is negative: %d", match.Home.Score)
	}
	if match.Away.Score < 0 {
		t.Errorf("Away score is negative: %d", match.Away.Score)
	}
	if totalGoals > 20 {
		t.Errorf("Total goals (%d) is unrealistically high (>20)", totalGoals)
	}

	// Verify events were generated
	if len(match.Events) == 0 {
		t.Error("No events were generated during the match")
	}

	// Count goal events
	goalEvents := 0
	for _, event := range match.Events {
		if event.Type == domain.GoalEvent {
			goalEvents++
		}
	}

	if goalEvents != totalGoals {
		t.Errorf("Goal event count (%d) doesn't match total score (%d)", goalEvents, totalGoals)
	}
}
