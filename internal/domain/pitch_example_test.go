package domain

import (
	"fmt"
	"testing"
)

// TestPitchTopologyStructure verifies the pitch graph is well-formed
func TestPitchTopologyStructure(t *testing.T) {
	// Verify all zones have transitions defined (using GetValidTransitions)
	allZones := []PitchZone{
		WestLeftWing, WestLeftHalf, WestCentre, WestRightHalf, WestRightWing,
		WestMidLeftWing, WestMidLeftHalf, WestMidCentre, WestMidRightHalf, WestMidRightWing,
		EastMidLeftWing, EastMidLeftHalf, EastMidCentre, EastMidRightHalf, EastMidRightWing,
		EastLeftWing, EastLeftHalf, EastCentre, EastRightHalf, EastRightWing,
	}

	for _, zone := range allZones {
		transitions := GetValidTransitions(zone)
		if len(transitions) == 0 {
			t.Errorf("Zone %s has no transitions defined", GetZoneName(zone))
		}
	}
}

// TestZoneDepthAndLanes verifies helper functions work correctly
func TestZoneDepthAndLanes(t *testing.T) {
	tests := []struct {
		zone         PitchZone
		expectedDepth int // Row: 1=West, 2=West-Mid, 3=East-Mid, 4=East
		expectedLane  int // Col: 1=left wing, 2=left half, 3=centre, 4=right half, 5=right wing
	}{
		{WestLeftWing, 1, 1},   // West row, left wing
		{WestCentre, 1, 3},     // West row, centre
		{WestMidCentre, 2, 3},  // West-Mid row, centre
		{EastMidCentre, 3, 3},  // East-Mid row, centre
		{EastRightWing, 4, 5},  // East row, right wing
	}

	for _, tt := range tests {
		if depth := GetZoneDepth(tt.zone); depth != tt.expectedDepth {
			t.Errorf("GetZoneDepth(%s) = %d, want %d", GetZoneName(tt.zone), depth, tt.expectedDepth)
		}
		if lane := GetZoneLane(tt.zone); lane != tt.expectedLane {
			t.Errorf("GetZoneLane(%s) = %d, want %d", GetZoneName(tt.zone), lane, tt.expectedLane)
		}
	}
}

// TestAttackingTransitions verifies teams can always progress forward (except from East/West ends)
func TestAttackingTransitions(t *testing.T) {
	// Middle zones should have forward options (when attacking East)
	shouldHaveForward := []PitchZone{
		WestLeftWing, WestLeftHalf, WestCentre, WestRightHalf, WestRightWing,
		WestMidLeftWing, WestMidLeftHalf, WestMidCentre, WestMidRightHalf, WestMidRightWing,
		EastMidLeftWing, EastMidLeftHalf, EastMidCentre, EastMidRightHalf, EastMidRightWing,
	}

	for _, zone := range shouldHaveForward {
		attacking := GetAttackingTransitions(zone)
		if len(attacking) == 0 {
			t.Errorf("Zone %s has no attacking transitions - teams need a way to progress!", GetZoneName(zone))
		}
	}

	// East zones should NOT have forward options (they're already at East end when attacking East)
	eastZones := []PitchZone{EastLeftWing, EastLeftHalf, EastCentre, EastRightHalf, EastRightWing}
	for _, zone := range eastZones {
		attacking := GetAttackingTransitions(zone)
		if len(attacking) > 0 {
			t.Errorf("Zone %s has forward transitions but it's already at East end", GetZoneName(zone))
		}
	}
}

// Example_zoneProgression shows how to use the zone system in match simulation
func Example_zoneProgression() {
	// Start in West centre
	currentZone := WestCentre

	fmt.Printf("Starting in: %s\n", GetZoneName(currentZone))

	// Team has strong advantage, wants to attack (East)
	attacking := GetAttackingTransitions(currentZone)
	fmt.Printf("Can attack to: ")
	for i, t := range attacking {
		if i > 0 {
			fmt.Printf(", ")
		}
		fmt.Printf("%s (value: %d)", GetZoneName(t.To), t.AttackingValue)
	}
	fmt.Println()

	// Pick best attacking move
	best := GetBestAttackingTransition(currentZone)
	if best != nil {
		currentZone = best.To
		fmt.Printf("Moved to: %s\n", GetZoneName(currentZone))
	}

	// Team wants to play it safe, go lateral
	lateral := GetLateralTransitions(currentZone)
	fmt.Printf("Can go lateral to: ")
	for i, t := range lateral {
		if i > 0 {
			fmt.Printf(", ")
		}
		fmt.Printf("%s", GetZoneName(t.To))
	}
	fmt.Println()

	// Output:
	// Starting in: West Centre
	// Can attack to: West-Mid Left Half-Space (value: 1), West-Mid Centre (value: 1), West-Mid Right Half-Space (value: 1)
	// Moved to: West-Mid Centre
	// Can go lateral to: West-Mid Left Half-Space, West-Mid Right Half-Space
}

// Example_matchProgression shows simulating ball progression over multiple phases
func Example_matchProgression() {
	currentZone := WestLeftWing
	fmt.Printf("Starting zone: %s\n\n", GetZoneName(currentZone))

	// Simulate 5 phases of attacking play (attacking East)
	for phase := 1; phase <= 5; phase++ {
		fmt.Printf("Phase %d - In %s\n", phase, GetZoneName(currentZone))

		// Strong team, always picks best attacking option
		best := GetBestAttackingTransition(currentZone)
		if best != nil {
			currentZone = best.To
			fmt.Printf("  → Advanced to %s (attacking value: %d)\n", GetZoneName(currentZone), best.AttackingValue)
		} else {
			fmt.Println("  → No forward progress available (at East end)")
			break
		}

		if IsAttackingZone(currentZone) {
			fmt.Println("  → In shooting position!")
		}
		fmt.Println()
	}

	// Output:
	// Starting zone: West Left Wing
	//
	// Phase 1 - In West Left Wing
	//   → Advanced to West-Mid Left Wing (attacking value: 1)
	//
	// Phase 2 - In West-Mid Left Wing
	//   → Advanced to East-Mid Left Wing (attacking value: 1)
	//
	// Phase 3 - In East-Mid Left Wing
	//   → Advanced to East Left Wing (attacking value: 1)
	//   → In shooting position!
	//
	// Phase 4 - In East Left Wing
	//   → No forward progress available (at East end)
}
