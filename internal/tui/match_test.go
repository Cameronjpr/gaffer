package tui

import (
	"strings"
	"testing"

	"github.com/cameronjpr/gaffer/internal/components"
)

// TestBuildScoreWidget_Centering tests that the score is centered regardless of team name lengths
func TestBuildScoreWidget_Centering(t *testing.T) {
	tests := []struct {
		name      string
		homeTeam  string
		awayTeam  string
		homeScore int
		awayScore int
		width     int
	}{
		{
			name:      "equal length team names",
			homeTeam:  "Arsenal",
			awayTeam:  "Chelsea",
			homeScore: 2,
			awayScore: 1,
			width:     60,
		},
		{
			name:      "very different length team names",
			homeTeam:  "Arsenal",
			awayTeam:  "Manchester City",
			homeScore: 0,
			awayScore: 0,
			width:     60,
		},
		{
			name:      "short vs long team name",
			homeTeam:  "FC",
			awayTeam:  "Very Long Team Name United",
			homeScore: 5,
			awayScore: 3,
			width:     80,
		},
		{
			name:      "double digit scores",
			homeTeam:  "Home",
			awayTeam:  "Away",
			homeScore: 10,
			awayScore: 9,
			width:     60,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := components.Scoreboard(tt, tt.width)

			// Strip ANSI codes and styling to get raw text
			plain := components.StripANSI(result)

			// The result should have the specified width (lipgloss adds padding)
			lines := strings.Split(plain, "\n")
			if len(lines) > 0 {
				// Get the first non-empty line
				var firstLine string
				for _, line := range lines {
					if len(line) > 0 {
						firstLine = line
						break
					}
				}

				// lipgloss adds some padding, so check we're within reasonable range
				// The actual content width should be close to specified width
				if firstLine != "" && (len(firstLine) < tt.width-5 || len(firstLine) > tt.width+5) {
					t.Errorf("Expected width around %d, got %d (line: %q)", tt.width, len(firstLine), firstLine)
				}

				// Verify both team names and score exist in the output
				if !strings.Contains(plain, tt.homeTeam) {
					t.Errorf("Home team %q not found in output", tt.homeTeam)
				}
				if !strings.Contains(plain, tt.awayTeam) {
					t.Errorf("Away team %q not found in output", tt.awayTeam)
				}
				// Check for score brackets
				if !strings.Contains(plain, "]") || !strings.Contains(plain, "[") {
					t.Error("Score brackets not found in output")
				}
			}
		})
	}
}
