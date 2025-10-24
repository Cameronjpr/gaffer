package domain

type PitchZone int

// 4x5 grid: 4 horizontal rows (West to East) × 5 vertical lanes (left wing to right wing)
// Home team attacks East (Row 4), Away team attacks West (Row 1)
const (
	// Row 1: West end of pitch
	WestLeftWing PitchZone = iota + 1
	WestLeftHalf
	WestCentre
	WestRightHalf
	WestRightWing

	// Row 2: West-Mid
	WestMidLeftWing
	WestMidLeftHalf
	WestMidCentre
	WestMidRightHalf
	WestMidRightWing

	// Row 3: East-Mid
	EastMidLeftWing
	EastMidLeftHalf
	EastMidCentre
	EastMidRightHalf
	EastMidRightWing

	// Row 4: East end of pitch
	EastLeftWing
	EastLeftHalf
	EastCentre
	EastRightHalf
	EastRightWing
)

// ZoneTransition represents a possible move from one zone to another
type ZoneTransition struct {
	To             PitchZone
	AttackingValue int // Higher = more attacking (0=lateral, positive=forward, negative=backward)
	IsLateral      bool
	IsForward      bool
	IsBackward     bool
}

// GetZoneRow returns the row number (1-4) for a zone, where 1=West, 4=East
func GetZoneRow(zone PitchZone) int {
	if zone >= WestLeftWing && zone <= WestRightWing {
		return 1
	} else if zone >= WestMidLeftWing && zone <= WestMidRightWing {
		return 2
	} else if zone >= EastMidLeftWing && zone <= EastMidRightWing {
		return 3
	} else if zone >= EastLeftWing && zone <= EastRightWing {
		return 4
	}
	return 0
}

// GetZoneCol returns the column number (1-5) for a zone, where 1=left wing, 3=centre, 5=right wing
func GetZoneCol(zone PitchZone) int {
	// Each row has 5 zones, so we can use modulo arithmetic
	col := int(zone-1) % 5
	return col + 1
}

// GetZoneFromRowCol returns the PitchZone for given row and column (returns 0 if invalid)
func GetZoneFromRowCol(row, col int) PitchZone {
	if row < 1 || row > 4 || col < 1 || col > 5 {
		return 0
	}
	return PitchZone((row-1)*5 + col)
}

// GetValidTransitions returns all valid transitions from a zone
// Rules: 8-directional movement (adjacent cells) + can always move to any lower row (backward passing)
func GetValidTransitions(from PitchZone) []ZoneTransition {
	fromRow := GetZoneRow(from)
	fromCol := GetZoneCol(from)

	transitions := make([]ZoneTransition, 0)
	seen := make(map[PitchZone]bool)

	// Rule 1: 8-directional adjacency (Queen-like movement)
	for rowDelta := -1; rowDelta <= 1; rowDelta++ {
		for colDelta := -1; colDelta <= 1; colDelta++ {
			if rowDelta == 0 && colDelta == 0 {
				continue // Skip self
			}

			toRow := fromRow + rowDelta
			toCol := fromCol + colDelta
			toZone := GetZoneFromRowCol(toRow, toCol)

			if toZone != 0 && !seen[toZone] {
				seen[toZone] = true
				transitions = append(transitions, createTransition(from, toZone))
			}
		}
	}

	// Rule 2: Can always move backward (to any zone in lower rows)
	// if fromRow > 1 {
	// 	for backRow := 1; backRow < fromRow; backRow++ {
	// 		for col := 1; col <= 5; col++ {
	// 			toZone := GetZoneFromRowCol(backRow, col)
	// 			if toZone != 0 && !seen[toZone] {
	// 				seen[toZone] = true
	// 				transitions = append(transitions, createTransition(from, toZone))
	// 			}
	// 		}
	// 	}
	// }

	return transitions
}

// createTransition creates a ZoneTransition with auto-calculated properties
func createTransition(from, to PitchZone) ZoneTransition {
	fromRow := GetZoneRow(from)
	toRow := GetZoneRow(to)

	rowDiff := toRow - fromRow

	return ZoneTransition{
		To:             to,
		AttackingValue: rowDiff,
		IsForward:      rowDiff > 0,
		IsLateral:      rowDiff == 0,
		IsBackward:     rowDiff < 0,
	}
}

// GetAttackingTransitions returns all forward-moving transitions from current zone.
// DEPRECATED: Use GetAttackingTransitionsForDirection instead to account for attacking direction.
// This function assumes attacking East (Row 1→4) for backward compatibility.
func GetAttackingTransitions(from PitchZone) []ZoneTransition {
	return GetAttackingTransitionsForDirection(from, AttackingEast)
}

// GetAttackingTransitionsForDirection returns all transitions toward the target goal.
// direction: AttackingEast (toward Row 4) or AttackingWest (toward Row 1)
func GetAttackingTransitionsForDirection(from PitchZone, direction AttackingDirection) []ZoneTransition {
	transitions := GetValidTransitions(from)
	attacking := make([]ZoneTransition, 0)

	for _, t := range transitions {
		// Check if this transition moves toward the target goal
		isAttacking := false
		if direction == AttackingEast {
			isAttacking = t.IsForward // Moving East (Row 1→4)
		} else {
			isAttacking = t.IsBackward // Moving West (Row 4→1)
		}

		if isAttacking {
			attacking = append(attacking, t)
		}
	}
	return attacking
}

// GetLateralTransitions returns all lateral transitions from current zone
func GetLateralTransitions(from PitchZone) []ZoneTransition {
	transitions := GetValidTransitions(from)
	lateral := make([]ZoneTransition, 0)
	for _, t := range transitions {
		if t.IsLateral {
			lateral = append(lateral, t)
		}
	}
	return lateral
}

// GetDefensiveTransitions returns all backward-moving transitions from current zone.
// DEPRECATED: Use GetDefensiveTransitionsForDirection instead to account for attacking direction.
// This function assumes attacking East (Row 1→4) for backward compatibility.
func GetDefensiveTransitions(from PitchZone) []ZoneTransition {
	return GetDefensiveTransitionsForDirection(from, AttackingEast)
}

// GetDefensiveTransitionsForDirection returns all transitions away from the target goal (toward own goal).
// direction: AttackingEast (toward Row 4) or AttackingWest (toward Row 1)
func GetDefensiveTransitionsForDirection(from PitchZone, direction AttackingDirection) []ZoneTransition {
	transitions := GetValidTransitions(from)
	defensive := make([]ZoneTransition, 0)

	for _, t := range transitions {
		// Check if this transition moves away from target goal
		isDefensive := false
		if direction == AttackingEast {
			isDefensive = t.IsBackward // Moving West (Row 4→1)
		} else {
			isDefensive = t.IsForward // Moving East (Row 1→4)
		}

		if isDefensive {
			defensive = append(defensive, t)
		}
	}
	return defensive
}

// GetBestAttackingTransition returns the transition with highest attacking value.
// DEPRECATED: Use GetBestAttackingTransitionForDirection instead to account for attacking direction.
// This function assumes attacking East (Row 1→4) for backward compatibility.
func GetBestAttackingTransition(from PitchZone) *ZoneTransition {
	return GetBestAttackingTransitionForDirection(from, AttackingEast)
}

// GetBestAttackingTransitionForDirection returns the transition with highest attacking value toward target goal.
// When multiple transitions have equal attacking value, it prefers maintaining
// the current lane (column), then closest lane, then most central.
// direction: AttackingEast (toward Row 4) or AttackingWest (toward Row 1)
func GetBestAttackingTransitionForDirection(from PitchZone, direction AttackingDirection) *ZoneTransition {
	attacking := GetAttackingTransitionsForDirection(from, direction)
	if len(attacking) == 0 {
		return nil
	}

	// Find the maximum attacking value
	maxValue := attacking[0].AttackingValue
	for i := range attacking {
		if attacking[i].AttackingValue > maxValue {
			maxValue = attacking[i].AttackingValue
		}
	}

	// Collect all transitions with max value
	var candidates []ZoneTransition
	for i := range attacking {
		if attacking[i].AttackingValue == maxValue {
			candidates = append(candidates, attacking[i])
		}
	}

	// If only one candidate, return it
	if len(candidates) == 1 {
		return &candidates[0]
	}

	// Tie-breaker: prefer maintaining current lane (column)
	fromCol := GetZoneCol(from)

	// First, check if any candidate is in the same column
	for i := range candidates {
		if GetZoneCol(candidates[i].To) == fromCol {
			return &candidates[i]
		}
	}

	// Second, find the closest column
	minColDist := 100
	var closest []ZoneTransition
	for i := range candidates {
		toCol := GetZoneCol(candidates[i].To)
		dist := toCol - fromCol
		if dist < 0 {
			dist = -dist
		}
		if dist < minColDist {
			minColDist = dist
			closest = []ZoneTransition{candidates[i]}
		} else if dist == minColDist {
			closest = append(closest, candidates[i])
		}
	}

	// If still multiple options, pick randomly
	if len(closest) > 1 {
		// Use a simple deterministic selection for now (can add proper random later if needed)
		// For determinism in tests, prefer center column, then left, then right
		bestCol := GetZoneCol(closest[0].To)
		bestIdx := 0
		for i := 1; i < len(closest); i++ {
			col := GetZoneCol(closest[i].To)
			// Prefer column 3 (center), then columns closer to center
			centerDist := col - 3
			if centerDist < 0 {
				centerDist = -centerDist
			}
			bestCenterDist := bestCol - 3
			if bestCenterDist < 0 {
				bestCenterDist = -bestCenterDist
			}
			if centerDist < bestCenterDist {
				bestCol = col
				bestIdx = i
			}
		}
		return &closest[bestIdx]
	}

	return &closest[0]
}

// GetZoneDepth returns how far up the pitch a zone is (1-4, where 1=defensive, 4=attacking)
// This is just an alias for GetZoneRow for backward compatibility
func GetZoneDepth(zone PitchZone) int {
	return GetZoneRow(zone)
}

// GetZoneLane returns the horizontal lane (1=left wing, 3=centre, 5=right wing)
// This is just an alias for GetZoneCol for backward compatibility
func GetZoneLane(zone PitchZone) int {
	return GetZoneCol(zone)
}

// IsAttackingZone returns true if zone is in attacking third (row 4)
func IsAttackingZone(zone PitchZone) bool {
	return GetZoneRow(zone) == 4
}

// IsMidfieldZone returns true if zone is in midfield (rows 2 or 3)
func IsMidfieldZone(zone PitchZone) bool {
	row := GetZoneRow(zone)
	return row == 2 || row == 3
}

// IsDefensiveZone returns true if zone is in defensive third (row 1)
func IsDefensiveZone(zone PitchZone) bool {
	return GetZoneRow(zone) == 1
}

// GetZoneName returns a human-readable name for the zone
func GetZoneName(zone PitchZone) string {
	names := map[PitchZone]string{
		WestLeftWing:  "West Left Wing",
		WestLeftHalf:  "West Left Half-Space",
		WestCentre:    "West Centre",
		WestRightHalf: "West Right Half-Space",
		WestRightWing: "West Right Wing",

		WestMidLeftWing:  "West-Mid Left Wing",
		WestMidLeftHalf:  "West-Mid Left Half-Space",
		WestMidCentre:    "West-Mid Centre",
		WestMidRightHalf: "West-Mid Right Half-Space",
		WestMidRightWing: "West-Mid Right Wing",

		EastMidLeftWing:  "East-Mid Left Wing",
		EastMidLeftHalf:  "East-Mid Left Half-Space",
		EastMidCentre:    "East-Mid Centre",
		EastMidRightHalf: "East-Mid Right Half-Space",
		EastMidRightWing: "East-Mid Right Wing",

		EastLeftWing:  "East Left Wing",
		EastLeftHalf:  "East Left Half-Space",
		EastCentre:    "East Centre",
		EastRightHalf: "East Right Half-Space",
		EastRightWing: "East Right Wing",
	}
	return names[zone]
}

// AttackingDirection represents which goal a team is attacking toward
type AttackingDirection int

const (
	AttackingEast AttackingDirection = iota // Attacking toward East end (Row 4)
	AttackingWest                           // Attacking toward West end (Row 1)
)

// GetShotThreat returns the base goal probability for shots from this zone (0.0-1.0).
// DEPRECATED: Use GetShotThreatForDirection instead to account for attacking direction.
// This function assumes attacking East for backward compatibility.
func GetShotThreat(zone PitchZone) float64 {
	return GetShotThreatForDirection(zone, AttackingEast)
}

// GetShotThreatForDirection returns the base goal probability for shots from this zone (0.0-1.0).
// Based on xG (expected goals) data - shots from closer to goal have higher conversion rates.
// direction: AttackingEast (toward Row 4) or AttackingWest (toward Row 1)
func GetShotThreatForDirection(zone PitchZone, direction AttackingDirection) float64 {
	row := GetZoneRow(zone)
	col := GetZoneCol(zone)

	// Calculate distance from target goal (1 = at goal, 4 = own half)
	var distanceFromGoal int
	if direction == AttackingEast {
		distanceFromGoal = 5 - row // Row 4 (East) = 1, Row 1 (West) = 4
	} else {
		distanceFromGoal = row // Row 1 (West) = 1, Row 4 (East) = 4
	}

	// Base threat by distance from goal
	var baseThreat float64
	switch distanceFromGoal {
	case 1: // At goal (penalty area)
		baseThreat = 1.0
	case 2: // Edge of box
		baseThreat = 0.45
	case 3: // Midfield
		baseThreat = 0.15
	case 4: // Own half
		baseThreat = 0.05
	}

	// Modify by column (centre is best, wings are worse)
	var columnModifier float64
	if col == 3 { // Centre
		columnModifier = 1.0
	} else if col == 2 || col == 4 { // Half-spaces
		columnModifier = 0.8
	} else { // Wings (1 or 5)
		columnModifier = 0.55
	}

	return baseThreat * columnModifier
}
