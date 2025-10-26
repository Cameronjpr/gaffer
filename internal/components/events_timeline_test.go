package components

import (
	"strings"
	"testing"

	"github.com/cameronjpr/gaffer/internal/domain"
	"github.com/charmbracelet/lipgloss"
)

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
			result := EventsTimeline(tt.homeEvents, tt.awayEvents, tt.colWidth)

			// Strip ANSI for analysis
			plain := StripANSI(result)

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
			rightResult := EventsTimelineForTeam(events, width, lipgloss.Right)
			rightPlain := StripANSI(rightResult)
			rightLines := strings.Split(rightPlain, "\n")

			// Test left-aligned
			leftResult := EventsTimelineForTeam(events, width, lipgloss.Left)
			leftPlain := StripANSI(leftResult)
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
	rightResult := EventsTimelineForTeam(events, width, lipgloss.Right)
	rightPlain := StripANSI(rightResult)
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
	leftResult := EventsTimelineForTeam(events, width, lipgloss.Left)
	leftPlain := StripANSI(leftResult)
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
			score := Scoreboard("Home", "Away", 1, 0, colWidth)
			scorePlain := StripANSI(score)
			scoreLines := strings.Split(scorePlain, "\n")

			// Build a timeline
			participant := &domain.MatchParticipant{Club: &domain.Club{Name: "Test Team"}}
			player := &domain.MatchPlayerParticipant{Player: &domain.Player{Name: "Test", Quality: 18}}
			events := []domain.Event{
				{Type: domain.GoalEvent, Minute: 10, For: participant, Player: player},
			}
			timeline := EventsTimeline(events, []domain.Event{}, colWidth)
			timelinePlain := StripANSI(timeline)
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
