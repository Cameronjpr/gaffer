package simulation

import (
	"math"
	"math/rand/v2"

	"github.com/cameronjpr/gaffer/internal/domain"
)

const goalscoringThreshold = 13
const shotThreshold = 8

// Engine runs the match simulation
type Engine struct {
	Match *domain.Match
}

// NewEngine creates a new match engine
func NewEngine(match *domain.Match) *Engine {
	return &Engine{Match: match}
}

// ProgressBall attempts to move the ball to a better zone based on power difference
func (e *Engine) ProgressBall(powerDiff int) {
	// Strong advantage - try to attack
	if powerDiff >= 10 {
		attacking := domain.GetAttackingTransitions(e.Match.ActiveZone)
		if len(attacking) > 0 {
			// Pick best attacking transition or random if multiple equal
			best := domain.GetBestAttackingTransition(e.Match.ActiveZone)
			if best != nil {
				e.Match.ActiveZone = best.To
			}
		}
		return
	}

	// Moderate advantage - mix of lateral and forward
	if powerDiff >= 5 {
		allMoves := domain.PitchTopology[e.Match.ActiveZone]
		// Prefer forward, but allow lateral
		var validMoves []domain.ZoneTransition
		for _, move := range allMoves {
			if move.IsForward || move.IsLateral {
				validMoves = append(validMoves, move)
			}
		}
		if len(validMoves) > 0 {
			// Weight towards forward moves
			forwardMoves := domain.GetAttackingTransitions(e.Match.ActiveZone)
			if len(forwardMoves) > 0 && rand.IntN(100) < 70 {
				e.Match.ActiveZone = forwardMoves[rand.IntN(len(forwardMoves))].To
			} else if len(validMoves) > 0 {
				e.Match.ActiveZone = validMoves[rand.IntN(len(validMoves))].To
			}
		}
		return
	}

	// Small advantage - mostly lateral, some backward
	lateralMoves := domain.GetLateralTransitions(e.Match.ActiveZone)
	if len(lateralMoves) > 0 && rand.IntN(100) < 60 {
		e.Match.ActiveZone = lateralMoves[rand.IntN(len(lateralMoves))].To
	}
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
		e.Match.AddEvent(domain.NewEvent(domain.PossessionChangedEvent, e.Match.CurrentMinute, morePowerfulTeam, nil))
		e.Match.TeamInPossession = morePowerfulTeam

		// When possession changes, ball likely moves backward for new team
		defensiveMoves := domain.GetDefensiveTransitions(e.Match.ActiveZone)
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
	e.ProgressBall(powerDiff)

	// Team kept the ball, no action taken
	if powerDiff < shotThreshold {
		return domain.PhaseResult{
			HomeRoll:       homeRoll,
			AwayRoll:       awayRoll,
			HomePhasePower: homePhasePower,
			AwayPhasePower: awayPhasePower,
			HomeGoals:      0,
			AwayGoals:      0,
		}
	}

	// Team has ball, but didn't take a shot
	if powerDiff < goalscoringThreshold {
		if powerDiff%2 == 0 {
			e.Match.AddEvent(domain.NewEvent(domain.SavedShotEvent, e.Match.CurrentMinute, morePowerfulTeam, nil))
		} else {
			e.Match.AddEvent(domain.NewEvent(domain.MissedShotEvent, e.Match.CurrentMinute, morePowerfulTeam, nil))
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

	goalsThisPhase := 0
	if powerDiff > goalscoringThreshold {
		goalsThisPhase = rand.IntN(2)
	}

	homeGoals := 0
	awayGoals := 0

	if goalsThisPhase > 0 {
		if homePhasePower > awayPhasePower {
			homeGoals = goalsThisPhase
			scorer := e.Match.Home.GetRandomOutfielder()
			e.Match.AddEvent(domain.NewEvent(domain.GoalEvent, e.Match.CurrentMinute, e.Match.Home, &scorer))
		}

		if homePhasePower < awayPhasePower {
			awayGoals = goalsThisPhase
			scorer := e.Match.Away.GetRandomOutfielder()
			e.Match.AddEvent(domain.NewEvent(domain.GoalEvent, e.Match.CurrentMinute, e.Match.Away, &scorer))
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
