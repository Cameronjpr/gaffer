package components

import (
	"strings"
	"testing"

	"github.com/cameronjpr/gaffer/internal/domain"
)

// TestPitch verifies the 5x4 grid renders correctly (rotated to horizontal) with goals
func TestPitch(t *testing.T) {
	tests := []struct {
		name string
		zone domain.PitchZone
	}{
		{"East Centre", domain.EastCentre},
		{"West-Mid Left Wing", domain.WestMidLeftWing},
		{"West Right Wing", domain.WestRightWing},
		{"East-Mid Centre", domain.EastMidCentre},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Pitch(tt.zone, nil)
			lines := strings.Split(result, "\n")

			// With blank lines between rows: 5 lanes + 4 blank lines = 9 total lines
			if len(lines) != 9 {
				t.Errorf("Expected 9 lines (5 lanes + 4 blanks), got %d", len(lines))
			}

			// Verify exactly one ● exists (active zone)
			dotCount := strings.Count(result, "●")
			if dotCount != 1 {
				t.Errorf("Expected exactly 1 ●, got %d. Output:\n%s", dotCount, result)
			}

			// Verify goals are present (1 opening bracket, 1 closing bracket)
			westGoalCount := strings.Count(result, "[")
			eastGoalCount := strings.Count(result, "]")
			if westGoalCount != 1 {
				t.Errorf("Expected 1 West goal '[', got %d", westGoalCount)
			}
			if eastGoalCount != 1 {
				t.Errorf("Expected 1 East goal ']', got %d", eastGoalCount)
			}
		})
	}
}
