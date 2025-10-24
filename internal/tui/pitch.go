package tui

import (
	"github.com/cameronjpr/gaffer/internal/domain"
	"github.com/charmbracelet/lipgloss"
)

// buildZoneIndicator creates a 5x4 grid showing current active zone with goals
// Layout: 5 rows (lanes: LW, LH, C, RH, RW) × 4 columns (West to East)
func buildZoneIndicator(zone domain.PitchZone, match *domain.Match) string {
	// Map zones to grid positions (horizontal pitch: West to East, left to right)
	// Row represents lane (0=Left Wing, 4=Right Wing)
	// Column represents depth (0=West, 3=East)
	zoneMap := map[domain.PitchZone][2]int{
		// Column 0: West end
		domain.WestLeftWing:  {0, 0}, // row=Left Wing lane, col=West
		domain.WestLeftHalf:  {1, 0}, // row=Left Half lane, col=West
		domain.WestCentre:    {2, 0}, // row=Centre lane, col=West
		domain.WestRightHalf: {3, 0}, // row=Right Half lane, col=West
		domain.WestRightWing: {4, 0}, // row=Right Wing lane, col=West

		// Column 1: West-Mid
		domain.WestMidLeftWing:  {0, 1},
		domain.WestMidLeftHalf:  {1, 1},
		domain.WestMidCentre:    {2, 1},
		domain.WestMidRightHalf: {3, 1},
		domain.WestMidRightWing: {4, 1},

		// Column 2: East-Mid
		domain.EastMidLeftWing:  {0, 2},
		domain.EastMidLeftHalf:  {1, 2},
		domain.EastMidCentre:    {2, 2},
		domain.EastMidRightHalf: {3, 2},
		domain.EastMidRightWing: {4, 2},

		// Column 3: East end
		domain.EastLeftWing:  {0, 3},
		domain.EastLeftHalf:  {1, 3},
		domain.EastCentre:    {2, 3},
		domain.EastRightHalf: {3, 3},
		domain.EastRightWing: {4, 3},
	}

	// Build 5x4 grid (5 rows for lanes, 4 columns for depth)
	grid := [5][4]string{}
	for z, pos := range zoneMap {
		row, col := pos[0], pos[1]
		if z == zone {
			grid[row][col] = "   ●   " // Active zone
		} else {
			grid[row][col] = "   ·   " // Inactive zone
		}
	}

	// Determine goal colors based on which half we're in
	// First half: Home attacks East (]), Away attacks West ([)
	// Second half: Teams switch sides
	var westGoalColor, eastGoalColor lipgloss.Color
	if match != nil {
		if match.HomeAttackingDirection == domain.AttackingEast {
			// First half: Home attacks East, Away defends West
			westGoalColor = lipgloss.Color(match.Home.Club.Background)
			eastGoalColor = lipgloss.Color(match.Away.Club.Background)
		} else {
			// First half: Away attacks West, Home defends East
			westGoalColor = lipgloss.Color(match.Away.Club.Background)
			eastGoalColor = lipgloss.Color(match.Home.Club.Background)
		}
	} else {
		// Default colors if match is nil
		westGoalColor = lipgloss.Color("240")
		eastGoalColor = lipgloss.Color("240")
	}

	westGoalStyle := lipgloss.NewStyle().Foreground(westGoalColor)
	eastGoalStyle := lipgloss.NewStyle().Foreground(eastGoalColor)

	// Render grid (5 rows × 4 columns, West to East left-to-right) with goals
	result := ""
	for row := range 5 {
		// West goal bracket
		if row == 2 {
			result += westGoalStyle.Render("[")
		} else {
			result += " "
		}

		// Pitch zones
		for col := range 4 {
			result += grid[row][col]
		}

		// East goal bracket
		if row == 2 {
			result += eastGoalStyle.Render("]")
		} else {
			result += ""
		}

		if row < 4 {
			result += "\n\n"
		}
	}

	return result
}
