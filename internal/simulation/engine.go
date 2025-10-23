package simulation

import (
	"math"
	"math/rand/v2"

	"github.com/cameronjpr/gaffer/internal/domain"
)

// Zone progression scaling factor
// Higher values make it harder to progress into dangerous zones
// Lower values make progression easier
// With team strengths around 15-20 and d20 rolls, max powerDiff is ~40
// With 4x5 grid and granular threat values (0.02-1.0), this is tuned so:
// - Defensive zones (0.02-0.05): need ~0-0.1 power (free backward passing)
// - Def-Mid zones (0.08-0.15): need ~0.16-0.30 power
// - Att-Mid zones (0.25-0.45): need ~0.50-0.90 power
// - Attacking zones (0.55-1.0): need ~1.1-2.0 power
// This creates realistic penetration difficulty - reaching the box requires dominance
const zoneProgressionScaling = 2.0

// Shot quality thresholds
// These are applied AFTER zone threat is factored in
const shotOnTargetThreshold = 0.1  // 10% of shots miss the target entirely
const saveThreshold = 0.4          // 40% of on-target shots are saved (30% overall)

// Engine runs the match simulation
type Engine struct {
	Match *domain.Match
}

func NewEngine(match *domain.Match) *Engine {
	return &Engine{Match: match}
}

// ProgressBall attempts to move the ball to a better zone based on power difference.
// Returns true if the ball progressed to a new zone, false if stuck (forcing a shot or turnover).
// Uses zone threat to determine difficulty - higher threat zones require more power to reach.
func (e *Engine) ProgressBall(powerDiff int) bool {
	attackingDirection := e.Match.GetAttackingDirection()
	allTransitions := domain.GetValidTransitions(e.Match.ActiveZone)

	// Evaluate all possible moves
	var forwardMoves []domain.ZoneTransition
	var lateralMoves []domain.ZoneTransition
	var backwardMoves []domain.ZoneTransition

	for _, transition := range allTransitions {
		targetThreat := domain.GetShotThreatForDirection(transition.To, attackingDirection)
		// Required power is proportional to target zone threat
		requiredPower := int(targetThreat * zoneProgressionScaling)

		// Can we make this move with current power advantage?
		if powerDiff >= requiredPower {
			// Determine if this is forward/lateral/backward relative to attacking direction
			isForward := (attackingDirection == domain.AttackingEast && transition.IsForward) ||
				(attackingDirection == domain.AttackingWest && transition.IsBackward)
			isBackward := (attackingDirection == domain.AttackingEast && transition.IsBackward) ||
				(attackingDirection == domain.AttackingWest && transition.IsForward)

			if isForward {
				forwardMoves = append(forwardMoves, transition)
			} else if transition.IsLateral {
				lateralMoves = append(lateralMoves, transition)
			} else if isBackward {
				backwardMoves = append(backwardMoves, transition)
			}
		}
	}

	// Decision priority: Forward > Lateral > Backward
	// This creates realistic build-up play patterns

	// Try to go forward (70% chance if available)
	if len(forwardMoves) > 0 && rand.IntN(100) < 70 {
		// Prefer moves to higher-threat zones
		best := forwardMoves[0]
		for _, move := range forwardMoves {
			if domain.GetShotThreatForDirection(move.To, attackingDirection) > domain.GetShotThreatForDirection(best.To, attackingDirection) {
				best = move
			}
		}
		e.Match.ActiveZone = best.To
		return true
	}

	// Try lateral movement if forward blocked (60% chance)
	if len(lateralMoves) > 0 && rand.IntN(100) < 60 {
		e.Match.ActiveZone = lateralMoves[rand.IntN(len(lateralMoves))].To
		return true
	}

	// Try backward if necessary (40% chance)
	if len(backwardMoves) > 0 && rand.IntN(100) < 40 {
		e.Match.ActiveZone = backwardMoves[rand.IntN(len(backwardMoves))].To
		return true
	}

	// Can't progress - stuck in current position
	// This will force a shot attempt or potential turnover
	return false
}

// AttemptShot simulates a shot attempt based on zone threat and power advantage.
// Returns the number of goals scored (0 or 1).
// Shot quality is determined by: zone threat * power modifier
func (e *Engine) AttemptShot(powerDiff int) int {
	// Get base goal probability from zone threat (direction-aware)
	attackingDirection := e.Match.GetAttackingDirection()
	zoneThreat := domain.GetShotThreatForDirection(e.Match.ActiveZone, attackingDirection)

	// Power advantage modifies shot quality (each point of power adds 2% to chance)
	powerModifier := 1.0 + (float64(powerDiff) * 0.02)

	// Calculate final goal probability (capped at 0.9 to keep some realism)
	goalProbability := math.Min(zoneThreat*powerModifier, 0.9)

	// Roll for shot outcome
	shotRoll := rand.Float64()

	// Shot misses target entirely (10% of shots)
	if shotRoll < shotOnTargetThreshold {
		e.Match.AddEvent(domain.NewEvent(
			domain.MissedShotEvent,
			e.Match.CurrentMinute,
			e.Match.TeamInPossession,
			nil,
		))
		return 0
	}

	// Shot is on target - now check if it goes in
	// Adjust roll to 0-1 range for on-target shots
	onTargetRoll := (shotRoll - shotOnTargetThreshold) / (1.0 - shotOnTargetThreshold)

	// Shot is saved
	if onTargetRoll > goalProbability {
		e.Match.AddEvent(domain.NewEvent(
			domain.SavedShotEvent,
			e.Match.CurrentMinute,
			e.Match.TeamInPossession,
			nil,
		))
		return 0
	}

	// Goal!
	scorer := e.Match.TeamInPossession.GetRandomOutfielder()
	e.Match.AddEvent(domain.NewEvent(
		domain.GoalEvent,
		e.Match.CurrentMinute,
		e.Match.TeamInPossession,
		&scorer,
	))
	return 1
}

// PlayPhase simulates one phase of play (roughly one minute)
func (e *Engine) PlayPhase() domain.PhaseResult {
	homeRoll := rand.IntN(20)
	awayRoll := rand.IntN(20)

	homePhasePower := e.Match.Home.Club.Strength + homeRoll
	awayPhasePower := e.Match.Away.Club.Strength + awayRoll
	powerDiff := int(math.Abs(float64(homePhasePower - awayPhasePower)))

	morePowerfulTeam := e.Match.Home
	if homePhasePower < awayPhasePower {
		morePowerfulTeam = e.Match.Away
	}

	// Possibility to change possession
	if morePowerfulTeam != e.Match.TeamInPossession {
		e.Match.AddEvent(domain.NewEvent(
			domain.PossessionChangedEvent,
			e.Match.CurrentMinute,
			morePowerfulTeam,
			nil,
		))
		e.Match.TeamInPossession = morePowerfulTeam

		// When possession changes, ball likely moves backward for new team
		// The new team's attacking direction determines what "backward" means
		newAttackingDirection := e.Match.GetAttackingDirection()
		defensiveMoves := domain.GetDefensiveTransitionsForDirection(e.Match.ActiveZone, newAttackingDirection)
		if len(defensiveMoves) > 0 {
			e.Match.ActiveZone = defensiveMoves[rand.IntN(len(defensiveMoves))].To
		}

		return domain.PhaseResult{
			HomeRoll:       homeRoll,
			AwayRoll:       awayRoll,
			HomePhasePower: homePhasePower,
			AwayPhasePower: awayPhasePower,
			HomeGoals:      0,
			AwayGoals:      0,
		}
	}

	// Team kept the ball, try to progress it
	ballProgressed := e.ProgressBall(powerDiff)

	// If ball can't progress further (in attacking position), attempt a shot
	var homeGoals, awayGoals int
	if !ballProgressed {
		goals := e.AttemptShot(powerDiff)
		if e.Match.TeamInPossession == e.Match.Home {
			homeGoals = goals
		} else {
			awayGoals = goals
		}
	}

	return domain.PhaseResult{
		HomeRoll:       homeRoll,
		AwayRoll:       awayRoll,
		HomePhasePower: homePhasePower,
		AwayPhasePower: awayPhasePower,
		HomeGoals:      homeGoals,
		AwayGoals:      awayGoals,
	}
}
