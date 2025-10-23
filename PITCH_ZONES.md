# Pitch Zone System

## Overview

The pitch is divided into a **3x3 grid** representing different areas of the field:

```
┌─────────────┬─────────────┬─────────────┐
│  AttLeft    │  AttCentre  │  AttRight   │  ← Attacking Third (closest to opponent goal)
├─────────────┼─────────────┼─────────────┤
│  MidLeft    │  MidCentre  │  MidRight   │  ← Midfield Third
├─────────────┼─────────────┼─────────────┤
│  DefLeft    │  DefCentre  │  DefRight   │  ← Defensive Third (closest to own goal)
└─────────────┴─────────────┴─────────────┘
```

## Zone Topology (Graph Structure)

Each zone has defined **transitions** (edges) to adjacent zones with the following properties:

- **AttackingValue**: Higher = more forward progress (0=lateral, positive=forward, negative=backward)
- **IsForward**: Boolean flag for attacking moves
- **IsLateral**: Boolean flag for sideways moves
- **IsBackward**: Boolean flag for defensive moves

### Example: DefCentre Transitions

```go
DefCentre: {
    {To: DefLeft, AttackingValue: 0, IsLateral: true},        // Sideways
    {To: DefRight, AttackingValue: 0, IsLateral: true},       // Sideways
    {To: MidCentre, AttackingValue: 2, IsForward: true},      // Best attack
    {To: MidLeft, AttackingValue: 1, IsForward: true},        // Diagonal attack
    {To: MidRight, AttackingValue: 1, IsForward: true},       // Diagonal attack
}
```

## Data Structure

### Core Data (pitch.go)

```go
// PitchTopology is a map defining all valid zone transitions
var PitchTopology = map[PitchZone][]ZoneTransition{
    DefLeft: {...},
    DefCentre: {...},
    // ... all 9 zones
}
```

### Helper Functions

```go
// Get specific transition types
GetAttackingTransitions(zone)   // Forward moves only
GetLateralTransitions(zone)     // Sideways moves only
GetDefensiveTransitions(zone)   // Backward moves only

// Find optimal moves
GetBestAttackingTransition(zone) // Highest attacking value

// Zone properties
GetZoneDepth(zone)              // 1=def, 2=mid, 3=att
GetZoneLane(zone)               // 1=left, 2=centre, 3=right
IsAttackingZone(zone)           // true if in attacking third
GetZoneName(zone)               // "Defensive Centre", etc.
```

## Usage in Match Simulation

### 1. Basic Zone Progression

```go
// Team in DefCentre wants to attack
currentZone := DefCentre

// Get all forward options
attacking := GetAttackingTransitions(currentZone)
// Returns: [MidCentre (value:2), MidLeft (value:1), MidRight (value:1)]

// Pick best option
best := GetBestAttackingTransition(currentZone)
currentZone = best.To  // Now in MidCentre
```

### 2. Power-Based Decision Making (match.go)

The `ProgressBall()` function uses power difference to decide movement:

```go
func (m *Match) ProgressBall(powerDiff int) {
    if powerDiff >= 10 {
        // Strong advantage - attack aggressively
        best := GetBestAttackingTransition(m.ActiveZone)
        if best != nil {
            m.ActiveZone = best.To
        }
    } else if powerDiff >= 5 {
        // Moderate advantage - mix forward/lateral (70% forward)
        // ... weighted random between forward and lateral moves
    } else {
        // Weak advantage - play safe, go lateral
        lateral := GetLateralTransitions(m.ActiveZone)
        // ... random lateral move
    }
}
```

### 3. Possession Changes

When possession changes, ball typically moves backward:

```go
if morePowerfulTeam != m.TeamInPossession {
    m.TeamInPossession = morePowerfulTeam

    // Ball goes backward for new possessing team
    defensiveMoves := GetDefensiveTransitions(m.ActiveZone)
    if len(defensiveMoves) > 0 {
        m.ActiveZone = defensiveMoves[rand.IntN(len(defensiveMoves))].To
    }
}
```

## Design Benefits

### 1. **Declarative Topology**
All zone relationships defined in one place (`PitchTopology` map). Easy to understand and modify.

### 2. **Type-Safe Transitions**
Can't move to invalid zones - only transitions defined in the graph are possible.

### 3. **Semantic Queries**
`GetAttackingTransitions()` is more readable than manual checks.

### 4. **Extensible**
Easy to add:
- Weighted transitions (e.g., centre preferred over wings)
- Formation-specific zones
- Tactical instructions (e.g., "always use wings")

### 5. **AI-Friendly**
Simple for AI to evaluate:
```go
// How many moves to goal?
depth := 3 - GetZoneDepth(currentZone)

// Is this a good zone?
if IsAttackingZone(currentZone) {
    // Take shot
}
```

## Example: Full Attack Sequence

```go
zone := DefLeft
fmt.Println(GetZoneName(zone))  // "Defensive Left"

// Phase 1: Strong team attacks
best := GetBestAttackingTransition(zone)
zone = best.To
fmt.Println(GetZoneName(zone))  // "Midfield Left"

// Phase 2: Continue attacking
best = GetBestAttackingTransition(zone)
zone = best.To
fmt.Println(GetZoneName(zone))  // "Attacking Left"

// Phase 3: In shooting position
if IsAttackingZone(zone) {
    fmt.Println("Take shot!")
}
```

## Future Enhancements

### 1. Zone-Specific Probabilities
```go
type ZoneTransition struct {
    To             PitchZone
    AttackingValue int
    Probability    float64  // Weight for random selection
}
```

### 2. Player Positioning
```go
// Which players are in which zones?
GetPlayersInZone(team, zone) []MatchPlayerParticipant
```

### 3. Formation Integration
```go
// 4-3-3 has more players in MidCentre
// 4-4-2 has more players in MidLeft/MidRight
GetZoneStrength(team, zone) int
```

### 4. Tactical Instructions
```go
type Tactic struct {
    PreferredLane int  // 1=left, 2=centre, 3=right
    Directness    int  // 0=patient buildup, 100=direct to goal
}
```

## Testing

Run tests to verify topology:
```bash
go test ./internal/game -v -run "TestPitch|TestZone|TestAttacking|Example"
```

Tests verify:
- All zones have transitions defined
- Defensive/midfield zones can progress forward
- Attacking zones have no forward options (already at goal)
- Helper functions return correct values
