package domain

type PitchZone int

const (
	DefLeft PitchZone = iota + 1
	DefCentre
	DefRight
	MidLeft
	MidCentre
	MidRight
	AttLeft
	AttCentre
	AttRight
)

// ZoneTransition represents a possible move from one zone to another
type ZoneTransition struct {
	To             PitchZone
	AttackingValue int // Higher = more attacking (0=lateral, positive=forward, negative=backward)
	IsLateral      bool
	IsForward      bool
	IsBackward     bool
}

// PitchTopology defines the structure of the pitch and valid transitions
var PitchTopology = map[PitchZone][]ZoneTransition{
	// Defensive Third
	DefLeft: {
		{To: DefCentre, AttackingValue: 0, IsLateral: true},
		{To: MidLeft, AttackingValue: 2, IsForward: true},
		{To: MidCentre, AttackingValue: 2, IsForward: true},
	},
	DefCentre: {
		{To: DefLeft, AttackingValue: 0, IsLateral: true},
		{To: DefRight, AttackingValue: 0, IsLateral: true},
		{To: MidCentre, AttackingValue: 2, IsForward: true},
		{To: MidLeft, AttackingValue: 1, IsForward: true},
		{To: MidRight, AttackingValue: 1, IsForward: true},
	},
	DefRight: {
		{To: DefCentre, AttackingValue: 0, IsLateral: true},
		{To: MidRight, AttackingValue: 2, IsForward: true},
		{To: MidCentre, AttackingValue: 2, IsForward: true},
	},

	// Midfield Third
	MidLeft: {
		{To: MidCentre, AttackingValue: 0, IsLateral: true},
		{To: DefLeft, AttackingValue: -2, IsBackward: true},
		{To: AttLeft, AttackingValue: 3, IsForward: true},
		{To: AttCentre, AttackingValue: 3, IsForward: true},
	},
	MidCentre: {
		{To: MidLeft, AttackingValue: 0, IsLateral: true},
		{To: MidRight, AttackingValue: 0, IsLateral: true},
		{To: DefCentre, AttackingValue: -2, IsBackward: true},
		{To: AttCentre, AttackingValue: 3, IsForward: true},
		{To: AttLeft, AttackingValue: 2, IsForward: true},
		{To: AttRight, AttackingValue: 2, IsForward: true},
	},
	MidRight: {
		{To: MidCentre, AttackingValue: 0, IsLateral: true},
		{To: DefRight, AttackingValue: -2, IsBackward: true},
		{To: AttRight, AttackingValue: 3, IsForward: true},
		{To: AttCentre, AttackingValue: 3, IsForward: true},
	},

	// Attacking Third
	AttLeft: {
		{To: AttCentre, AttackingValue: 0, IsLateral: true},
		{To: MidLeft, AttackingValue: -3, IsBackward: true},
	},
	AttCentre: {
		{To: AttLeft, AttackingValue: 0, IsLateral: true},
		{To: AttRight, AttackingValue: 0, IsLateral: true},
		{To: MidCentre, AttackingValue: -3, IsBackward: true},
	},
	AttRight: {
		{To: AttCentre, AttackingValue: 0, IsLateral: true},
		{To: MidRight, AttackingValue: -3, IsBackward: true},
	},
}

// GetAttackingTransitions returns all forward-moving transitions from current zone
func GetAttackingTransitions(from PitchZone) []ZoneTransition {
	transitions := PitchTopology[from]
	attacking := make([]ZoneTransition, 0)
	for _, t := range transitions {
		if t.IsForward {
			attacking = append(attacking, t)
		}
	}
	return attacking
}

// GetLateralTransitions returns all lateral transitions from current zone
func GetLateralTransitions(from PitchZone) []ZoneTransition {
	transitions := PitchTopology[from]
	lateral := make([]ZoneTransition, 0)
	for _, t := range transitions {
		if t.IsLateral {
			lateral = append(lateral, t)
		}
	}
	return lateral
}

// GetDefensiveTransitions returns all backward-moving transitions from current zone
func GetDefensiveTransitions(from PitchZone) []ZoneTransition {
	transitions := PitchTopology[from]
	defensive := make([]ZoneTransition, 0)
	for _, t := range transitions {
		if t.IsBackward {
			defensive = append(defensive, t)
		}
	}
	return defensive
}

// GetBestAttackingTransition returns the transition with highest attacking value
func GetBestAttackingTransition(from PitchZone) *ZoneTransition {
	attacking := GetAttackingTransitions(from)
	if len(attacking) == 0 {
		return nil
	}

	best := &attacking[0]
	for i := range attacking {
		if attacking[i].AttackingValue > best.AttackingValue {
			best = &attacking[i]
		}
	}
	return best
}

// GetZoneDepth returns how far up the pitch a zone is (1=defensive, 2=midfield, 3=attacking)
func GetZoneDepth(zone PitchZone) int {
	switch zone {
	case DefLeft, DefCentre, DefRight:
		return 1
	case MidLeft, MidCentre, MidRight:
		return 2
	case AttLeft, AttCentre, AttRight:
		return 3
	default:
		return 0
	}
}

// GetZoneLane returns the horizontal lane (1=left, 2=centre, 3=right)
func GetZoneLane(zone PitchZone) int {
	switch zone {
	case DefLeft, MidLeft, AttLeft:
		return 1
	case DefCentre, MidCentre, AttCentre:
		return 2
	case DefRight, MidRight, AttRight:
		return 3
	default:
		return 0
	}
}

// IsAttackingZone returns true if zone is in attacking third
func IsAttackingZone(zone PitchZone) bool {
	return GetZoneDepth(zone) == 3
}

// IsMidfieldZone returns true if zone is in midfield third
func IsMidfieldZone(zone PitchZone) bool {
	return GetZoneDepth(zone) == 2
}

// IsDefensiveZone returns true if zone is in defensive third
func IsDefensiveZone(zone PitchZone) bool {
	return GetZoneDepth(zone) == 1
}

// GetZoneName returns a human-readable name for the zone
func GetZoneName(zone PitchZone) string {
	names := map[PitchZone]string{
		DefLeft:   "Defensive Left",
		DefCentre: "Defensive Centre",
		DefRight:  "Defensive Right",
		MidLeft:   "Midfield Left",
		MidCentre: "Midfield Centre",
		MidRight:  "Midfield Right",
		AttLeft:   "Attacking Left",
		AttCentre: "Attacking Centre",
		AttRight:  "Attacking Right",
	}
	return names[zone]
}
