package domain

import (
	"fmt"
	"testing"
)

// TestPitchTopologyStructure verifies the pitch graph is well-formed
func TestPitchTopologyStructure(t *testing.T) {
	// Verify all zones have transitions defined
	allZones := []PitchZone{
		DefLeft, DefCentre, DefRight,
		MidLeft, MidCentre, MidRight,
		AttLeft, AttCentre, AttRight,
	}

	for _, zone := range allZones {
		transitions := PitchTopology[zone]
		if len(transitions) == 0 {
			t.Errorf("Zone %s has no transitions defined", GetZoneName(zone))
		}
	}
}

// TestZoneDepthAndLanes verifies helper functions work correctly
func TestZoneDepthAndLanes(t *testing.T) {
	tests := []struct {
		zone         PitchZone
		expectedDepth int
		expectedLane  int
	}{
		{DefLeft, 1, 1},
		{DefCentre, 1, 2},
		{MidCentre, 2, 2},
		{AttRight, 3, 3},
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

// TestAttackingTransitions verifies teams can always progress forward (except from attacking zones)
func TestAttackingTransitions(t *testing.T) {
	// Defensive and midfield zones should have forward options
	shouldHaveForward := []PitchZone{
		DefLeft, DefCentre, DefRight,
		MidLeft, MidCentre, MidRight,
	}

	for _, zone := range shouldHaveForward {
		attacking := GetAttackingTransitions(zone)
		if len(attacking) == 0 {
			t.Errorf("Zone %s has no attacking transitions - teams need a way to progress!", GetZoneName(zone))
		}
	}

	// Attacking zones should NOT have forward options (they're already at the top)
	attackingZones := []PitchZone{AttLeft, AttCentre, AttRight}
	for _, zone := range attackingZones {
		attacking := GetAttackingTransitions(zone)
		if len(attacking) > 0 {
			t.Errorf("Zone %s has forward transitions but it's already in attacking third", GetZoneName(zone))
		}
	}
}

// Example_zoneProgression shows how to use the zone system in match simulation
func Example_zoneProgression() {
	// Start in defensive centre
	currentZone := DefCentre

	fmt.Printf("Starting in: %s\n", GetZoneName(currentZone))

	// Team has strong advantage, wants to attack
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
	// Starting in: Defensive Centre
	// Can attack to: Midfield Centre (value: 2), Midfield Left (value: 1), Midfield Right (value: 1)
	// Moved to: Midfield Centre
	// Can go lateral to: Midfield Left, Midfield Right
}

// Example_matchProgression shows simulating ball progression over multiple phases
func Example_matchProgression() {
	currentZone := DefLeft
	fmt.Printf("Starting zone: %s\n\n", GetZoneName(currentZone))

	// Simulate 5 phases of attacking play
	for phase := 1; phase <= 5; phase++ {
		fmt.Printf("Phase %d - In %s\n", phase, GetZoneName(currentZone))

		// Strong team, always picks best attacking option
		best := GetBestAttackingTransition(currentZone)
		if best != nil {
			currentZone = best.To
			fmt.Printf("  → Advanced to %s (attacking value: %d)\n", GetZoneName(currentZone), best.AttackingValue)
		} else {
			fmt.Println("  → No forward progress available (at attacking zone)")
			break
		}

		if IsAttackingZone(currentZone) {
			fmt.Println("  → In shooting position!")
		}
		fmt.Println()
	}

	// Output:
	// Starting zone: Defensive Left
	//
	// Phase 1 - In Defensive Left
	//   → Advanced to Midfield Left (attacking value: 2)
	//
	// Phase 2 - In Midfield Left
	//   → Advanced to Attacking Left (attacking value: 3)
	//   → In shooting position!
	//
	// Phase 3 - In Attacking Left
	//   → No forward progress available (at attacking zone)
}
