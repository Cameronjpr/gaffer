package tui

import (
	"strings"
	"testing"

	"github.com/cameronjpr/gaffer/internal/domain"
	"github.com/charmbracelet/lipgloss"
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
			result := buildScoreWidget(tt.homeTeam, tt.awayTeam, tt.homeScore, tt.awayScore, tt.width)

			// Strip ANSI codes and styling to get raw text
			plain := stripANSI(result)

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

// TestBuildTimelineFromEvents_Alignment tests timeline alignment and width with new Event system
func TestBuildTimelineFromEvents_Alignment(t *testing.T) {
	// Create dummy match participants for testing
	homeClub := &domain.Club{Name: "Home Team"}
	awayClub := &domain.Club{Name: "Away Team"}
	homeParticipant := &domain.MatchParticipant{Club: homeClub}
	awayParticipant := &domain.MatchParticipant{Club: awayClub}

	player1 := &domain.MatchPlayerParticipant{Player: &domain.Player{Name: "Player1", Quality: 18}}
	player2 := &domain.MatchPlayerParticipant{Player: &domain.Player{Name: "Player2", Quality: 19}}
	player3 := &domain.MatchPlayerParticipant{Player: &domain.Player{Name: "Player3", Quality: 17}}

	tests := []struct {
		name       string
		homeEvents []domain.Event
		awayEvents []domain.Event
		colWidth   int
	}{
		{
			name: "balanced events",
			homeEvents: []domain.Event{
				{Type: domain.GoalEvent, Minute: 10, For: homeParticipant, Player: player1},
				{Type: domain.GoalEvent, Minute: 25, For: homeParticipant, Player: player2},
			},
			awayEvents: []domain.Event{
				{Type: domain.GoalEvent, Minute: 15, For: awayParticipant, Player: player3},
			},
			colWidth: 60,
		},
		{
			name: "unbalanced events - more home goals",
			homeEvents: []domain.Event{
				{Type: domain.GoalEvent, Minute: 5, For: homeParticipant, Player: player1},
				{Type: domain.GoalEvent, Minute: 15, For: homeParticipant, Player: player2},
				{Type: domain.GoalEvent, Minute: 45, For: homeParticipant, Player: player3},
			},
			awayEvents: []domain.Event{},
			colWidth:   60,
		},
		{
			name:       "no events",
			homeEvents: []domain.Event{},
			awayEvents: []domain.Event{},
			colWidth:   60,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := buildTimelineFromEvents(tt.homeEvents, tt.awayEvents, tt.colWidth)

			// Strip ANSI for analysis
			plain := stripANSI(result)

			// Verify the timeline has content (even if empty)
			if result == "" {
				t.Error("Timeline should not be empty string")
			}

			// Verify width matches colWidth (with some tolerance for lipgloss padding)
			lines := strings.Split(plain, "\n")
			if len(lines) > 0 {
				var firstLine string
				for _, line := range lines {
					if len(line) > 0 {
						firstLine = line
						break
					}
				}
				if firstLine != "" && (len(firstLine) < tt.colWidth-5 || len(firstLine) > tt.colWidth+5) {
					t.Errorf("Expected timeline width around %d, got %d", tt.colWidth, len(firstLine))
				}
			}
		})
	}
}

// TestBuildTimelineColumnFromEvents_Width tests that timeline columns have fixed widths
func TestBuildTimelineColumnFromEvents_Width(t *testing.T) {
	participant := &domain.MatchParticipant{Club: &domain.Club{Name: "Test Team"}}
	player1 := &domain.MatchPlayerParticipant{Player: &domain.Player{Name: "PlayerA", Quality: 18}}
	player2 := &domain.MatchPlayerParticipant{Player: &domain.Player{Name: "PlayerB", Quality: 19}}

	events := []domain.Event{
		{Type: domain.GoalEvent, Minute: 10, For: participant, Player: player1},
		{Type: domain.GoalEvent, Minute: 20, For: participant, Player: player2},
	}

	widths := []int{20, 30, 40}

	for _, width := range widths {
		t.Run("width_"+string(rune(width+'0')), func(t *testing.T) {
			// Test right-aligned
			rightResult := buildTimelineColumnFromEvents(events, width, lipgloss.Right)
			rightPlain := stripANSI(rightResult)
			rightLines := strings.Split(rightPlain, "\n")

			// Test left-aligned
			leftResult := buildTimelineColumnFromEvents(events, width, lipgloss.Left)
			leftPlain := stripANSI(leftResult)
			leftLines := strings.Split(leftPlain, "\n")

			// Both should have the same number of lines
			if len(rightLines) != len(leftLines) {
				t.Errorf("Right and left columns should have same line count")
			}

			// Each line should match the specified width (with tolerance for lipgloss)
			for i, line := range rightLines {
				if line != "" && (len(line) < width-5 || len(line) > width+5) {
					t.Errorf("Right-aligned line %d: expected width around %d, got %d", i, width, len(line))
				}
			}

			for i, line := range leftLines {
				if line != "" && (len(line) < width-5 || len(line) > width+5) {
					t.Errorf("Left-aligned line %d: expected width around %d, got %d", i, width, len(line))
				}
			}
		})
	}
}

// TestBuildTimelineColumnFromEvents_Alignment tests that timeline columns align correctly
func TestBuildTimelineColumnFromEvents_Alignment(t *testing.T) {
	participant := &domain.MatchParticipant{Club: &domain.Club{Name: "Test Team"}}
	player := &domain.MatchPlayerParticipant{Player: &domain.Player{Name: "Test", Quality: 18}}

	events := []domain.Event{
		{Type: domain.GoalEvent, Minute: 10, For: participant, Player: player},
	}
	width := 30

	// Test right-aligned - text should be at the end
	rightResult := buildTimelineColumnFromEvents(events, width, lipgloss.Right)
	rightPlain := stripANSI(rightResult)
	rightLines := strings.Split(rightPlain, "\n")
	if len(rightLines) > 0 {
		firstLine := rightLines[0]
		trimmed := strings.TrimLeft(firstLine, " ")
		// Right-aligned means leading spaces
		if len(trimmed) == len(firstLine) {
			t.Error("Right-aligned timeline should have leading spaces")
		}
	}

	// Test left-aligned - text should be at the start
	leftResult := buildTimelineColumnFromEvents(events, width, lipgloss.Left)
	leftPlain := stripANSI(leftResult)
	leftLines := strings.Split(leftPlain, "\n")
	if len(leftLines) > 0 {
		firstLine := leftLines[0]
		trimmed := strings.TrimRight(firstLine, " ")
		// Left-aligned means trailing spaces
		if len(trimmed) == len(firstLine) {
			t.Error("Left-aligned timeline should have trailing spaces")
		}
	}
}

// TestColumnLayout tests the 3-column layout proportions
func TestColumnLayout(t *testing.T) {
	widths := []int{90, 120, 150, 180}

	for _, totalWidth := range widths {
		t.Run("width_"+string(rune(totalWidth/10+'0')), func(t *testing.T) {
			colWidth := totalWidth / 3

			// Build a score widget
			score := buildScoreWidget("Home", "Away", 1, 0, colWidth)
			scorePlain := stripANSI(score)
			scoreLines := strings.Split(scorePlain, "\n")

			// Build a timeline
			participant := &domain.MatchParticipant{Club: &domain.Club{Name: "Test Team"}}
			player := &domain.MatchPlayerParticipant{Player: &domain.Player{Name: "Test", Quality: 18}}
			events := []domain.Event{
				{Type: domain.GoalEvent, Minute: 10, For: participant, Player: player},
			}
			timeline := buildTimelineFromEvents(events, []domain.Event{}, colWidth)
			timelinePlain := stripANSI(timeline)
			timelineLines := strings.Split(timelinePlain, "\n")

			// Verify each column has width around totalWidth / 3 (with tolerance)
			var scoreFirstLine string
			for _, line := range scoreLines {
				if len(line) > 0 {
					scoreFirstLine = line
					break
				}
			}
			if scoreFirstLine != "" && (len(scoreFirstLine) < colWidth-5 || len(scoreFirstLine) > colWidth+5) {
				t.Errorf("Score column width: expected around %d, got %d", colWidth, len(scoreFirstLine))
			}

			var timelineFirstLine string
			for _, line := range timelineLines {
				if len(line) > 0 {
					timelineFirstLine = line
					break
				}
			}
			if timelineFirstLine != "" && (len(timelineFirstLine) < colWidth-5 || len(timelineFirstLine) > colWidth+5) {
				t.Errorf("Timeline column width: expected around %d, got %d", colWidth, len(timelineFirstLine))
			}
		})
	}
}

// TestZoneIndicator verifies the 3x3 grid renders correctly
func TestZoneIndicator(t *testing.T) {
	tests := []struct {
		name        string
		zone        domain.PitchZone
		expectedDot [2]int // [row, col] where ● should appear
	}{
		{"Attacking Centre", domain.AttCentre, [2]int{0, 1}}, // top middle
		{"Midfield Left", domain.MidLeft, [2]int{1, 0}},      // middle left
		{"Defensive Right", domain.DefRight, [2]int{2, 2}},   // bottom right
		{"Midfield Centre", domain.MidCentre, [2]int{1, 1}},  // center
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := buildZoneIndicator(tt.zone, nil)
			lines := strings.Split(result, "\n")

			if len(lines) != 3 {
				t.Errorf("Expected 3 lines, got %d", len(lines))
			}

			// Check that ● appears in expected position
			expectedRow := tt.expectedDot[0]
			expectedCol := tt.expectedDot[1]

			// Parse the line (format: "X X X" where X is · or ●)
			if expectedRow < len(lines) {
				cols := strings.Split(lines[expectedRow], " ")
				if expectedCol < len(cols) {
					if cols[expectedCol] != "●" {
						t.Errorf("Expected ● at [%d,%d], got %q. Full output:\n%s",
							expectedRow, expectedCol, cols[expectedCol], result)
					}
				}
			}

			// Verify only one ● exists
			dotCount := strings.Count(result, "●")
			if dotCount != 1 {
				t.Errorf("Expected exactly 1 ●, got %d. Output:\n%s", dotCount, result)
			}
		})
	}
}

// stripANSI removes ANSI escape codes from a string for testing
func stripANSI(s string) string {
	// Simple ANSI stripper - removes common escape sequences
	result := s
	// Remove CSI sequences (most common)
	for strings.Contains(result, "\x1b[") {
		start := strings.Index(result, "\x1b[")
		end := start + 2
		for end < len(result) && !((result[end] >= 'A' && result[end] <= 'Z') || (result[end] >= 'a' && result[end] <= 'z')) {
			end++
		}
		if end < len(result) {
			end++
		}
		result = result[:start] + result[end:]
	}
	return result
}
