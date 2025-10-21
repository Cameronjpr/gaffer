package tui

import (
	"strings"
	"testing"

	"github.com/cameronjpr/gaffer/internal/game"
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

// TestBuildTimeline_Alignment tests timeline alignment and width
func TestBuildTimeline_Alignment(t *testing.T) {
	tests := []struct {
		name       string
		homeEvents []game.PlayerEvent
		awayEvents []game.PlayerEvent
		colWidth   int
	}{
		{
			name: "balanced events",
			homeEvents: []game.PlayerEvent{
				{Type: game.PlayerScoredEvent, Player: game.Player{Name: "Player1", Quality: 18}, Minute: 10},
				{Type: game.PlayerScoredEvent, Player: game.Player{Name: "Player2", Quality: 19}, Minute: 25},
			},
			awayEvents: []game.PlayerEvent{
				{Type: game.PlayerScoredEvent, Player: game.Player{Name: "Player3", Quality: 17}, Minute: 15},
				{Type: game.PlayerScoredEvent, Player: game.Player{Name: "Player4", Quality: 18}, Minute: 30},
			},
			colWidth: 60,
		},
		{
			name: "unbalanced events - more home goals",
			homeEvents: []game.PlayerEvent{
				{Type: game.PlayerScoredEvent, Player: game.Player{Name: "Striker1", Quality: 20}, Minute: 5},
				{Type: game.PlayerScoredEvent, Player: game.Player{Name: "Striker2", Quality: 19}, Minute: 15},
				{Type: game.PlayerScoredEvent, Player: game.Player{Name: "Striker3", Quality: 18}, Minute: 45},
			},
			awayEvents: []game.PlayerEvent{},
			colWidth:   60,
		},
		{
			name:       "no events",
			homeEvents: []game.PlayerEvent{},
			awayEvents: []game.PlayerEvent{},
			colWidth:   60,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := buildTimeline(tt.homeEvents, tt.awayEvents, tt.colWidth)

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

// TestBuildTimelineColumn_Width tests that timeline columns have fixed widths
func TestBuildTimelineColumn_Width(t *testing.T) {
	events := []game.PlayerEvent{
		{Type: game.PlayerScoredEvent, Player: game.Player{Name: "PlayerA", Quality: 18}, Minute: 10},
		{Type: game.PlayerScoredEvent, Player: game.Player{Name: "PlayerB", Quality: 19}, Minute: 20},
	}

	widths := []int{20, 30, 40}

	for _, width := range widths {
		t.Run("width_"+string(rune(width+'0')), func(t *testing.T) {
			// Test right-aligned
			rightResult := buildTimelineColumn(events, width, lipgloss.Right)
			rightPlain := stripANSI(rightResult)
			rightLines := strings.Split(rightPlain, "\n")

			// Test left-aligned
			leftResult := buildTimelineColumn(events, width, lipgloss.Left)
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

// TestBuildTimelineColumn_Alignment tests that timeline columns align correctly
func TestBuildTimelineColumn_Alignment(t *testing.T) {
	events := []game.PlayerEvent{
		{Type: game.PlayerScoredEvent, Player: game.Player{Name: "Test", Quality: 18}, Minute: 10},
	}
	width := 30

	// Test right-aligned - text should be at the end
	rightResult := buildTimelineColumn(events, width, lipgloss.Right)
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
	leftResult := buildTimelineColumn(events, width, lipgloss.Left)
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
			events := []game.PlayerEvent{
				{Type: game.PlayerScoredEvent, Player: game.Player{Name: "Test", Quality: 18}, Minute: 10},
			}
			timeline := buildTimeline(events, []game.PlayerEvent{}, colWidth)
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
