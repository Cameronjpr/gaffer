# Zone Indicator UI Component

## Overview

A simple 3×3 grid visualization showing the current active zone of play during a match. Located in the center column, below the timeline.

## Visual Layout

The grid is oriented with **attacking zones at the top**:

```
· · ·    ← Attacking third (opponent's goal)
· ● ·    ← Midfield third (● = current active zone)
· · ·    ← Defensive third (own goal)
```

## Zone Mapping

```
Row 0 (Top):    AttLeft   AttCentre   AttRight
Row 1 (Middle): MidLeft   MidCentre   MidRight
Row 2 (Bottom): DefLeft   DefCentre   DefRight
```

## Examples

### Ball in Attacking Centre
```
· ● ·
· · ·
· · ·
```
Team is in prime shooting position!

### Ball in Midfield Left
```
· · ·
● · ·
· · ·
```
Team progressing down the left wing.

### Ball in Defensive Right
```
· · ·
· · ·
· · ●
```
Defending in own third, right side.

## Implementation Details

### Function: `buildZoneIndicator`
**Location**: `internal/tui/match.go:112-149`

```go
func buildZoneIndicator(zone game.PitchZone, teamInPossession *game.MatchParticipant) string
```

**Parameters**:
- `zone`: Current active pitch zone
- `teamInPossession`: Currently unused, but available for future enhancements (e.g., color-coding)

**Returns**: String representation of 3×3 grid with newlines

### Rendering

Added to ticker content in `View()`:

```go
tickerContent := lipgloss.JoinVertical(
    lipgloss.Center,
    scoreWidget,
    time,
    gap,
    timeline,
    gap,
    lipgloss.NewStyle().Faint(true).Render(zoneIndicator),  // ← Here
)
```

**Styling**: Uses `.Faint(true)` to make it subtle and unobtrusive.

## Match Screen Layout

```
┌─────────────────────────────────────────────────────────┐
│  Home Team          Score & Time         Away Team      │
│  Formation          (1:00)               Formation      │
│  Lineup             Timeline             Lineup         │
│                     · · ·    ← Zone indicator           │
│                     · ● ·                               │
│                     · · ·                               │
└─────────────────────────────────────────────────────────┘
│              Commentary Footer                          │
└─────────────────────────────────────────────────────────┘
```

## Match State Integration

### Initial State (Kickoff)
Match starts at `MidCentre` with home team in possession:

```go
func NewMatch(homeClub, awayClub *Club) Match {
    return Match{
        TeamInPossession: home,
        ActiveZone:       MidCentre,  // Kickoff from center circle
        // ...
    }
}
```

Indicator shows:
```
· · ·
· ● ·
· · ·
```

### During Play
Updates automatically as `match.ActiveZone` changes via `ProgressBall()` in `PlayPhase()`.

## Future Enhancements

### 1. Color-Coding by Team
```go
if teamInPossession == match.Home {
    grid[row][col] = "●"  // Could use team color
} else {
    grid[row][col] = "○"  // Different marker for away team
}
```

### 2. Direction Indicator
Show attacking direction with arrow:
```
· · ·  ↑
· ● ·  ↑ Attacking this way
· · ·  ↑
```

### 3. Heat Map
Track zone activity over time:
```
2 4 1    ← Number of times ball was in each zone
5 8 3
1 2 0
```

### 4. Pressure Visualization
Show defensive pressure with multiple markers:
```
· ● ●    ← Multiple dots = congested area
● ● ·
· · ·
```

## Testing

**Test**: `TestZoneIndicator` in `internal/tui/match_test.go:306-349`

Verifies:
- Grid is exactly 3 lines tall
- Active zone marked with ●
- Only one ● appears
- Correct [row, col] positioning

Run tests:
```bash
go test ./internal/tui -v -run TestZoneIndicator
```

## Related Files

- **Pitch topology**: `internal/game/pitch.go`
- **Match state**: `internal/game/match.go`
- **UI rendering**: `internal/tui/match.go`
- **Tests**: `internal/tui/match_test.go`
