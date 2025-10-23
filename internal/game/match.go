package game

import (
	"math"
	"math/rand/v2"
)

const goalscoringThreshold = 15
const shotThreshold = 10

type Half int

const (
	FirstHalf Half = iota + 1
	SecondHalf
)

type Match struct {
	Home             *MatchParticipant
	Away             *MatchParticipant
	TeamInPossession *MatchParticipant
	CurrentMinute    int
	CurrentHalf      Half
	ActiveZone       PitchZone
	PhaseHistory     []PhaseResult
	Events           []Event
}

type PhaseResult struct {
	Phase          int
	HomeRoll       int
	AwayRoll       int
	HomePhasePower int
	AwayPhasePower int
	HomeGoals      int
	AwayGoals      int
}

func NewMatch(homeClub, awayClub *Club) Match {
	home := NewMatchParticipant(homeClub)
	return Match{
		Home:             home,
		Away:             NewMatchParticipant(awayClub),
		TeamInPossession: home, // Home team starts with kickoff
		CurrentMinute:    1,
		CurrentHalf:      FirstHalf,
		ActiveZone:       MidCentre, // Kickoff from center
		PhaseHistory:     make([]PhaseResult, 0),
		Events:           make([]Event, 0),
	}
}

func (m *Match) StartFirstHalf() {
	m.CurrentHalf = FirstHalf
	m.CurrentMinute = 1
}

func (m *Match) StartSecondHalf() {
	m.CurrentHalf = SecondHalf
	m.CurrentMinute = 45 // Will be incremented to 46 on first tick
}

func (m *Match) AddEvent(event Event) {
	m.Events = append(m.Events, event)
}

func (m *Match) GetAddedTime(half Half) int {
	addedTime := 0
	for _, event := range m.Events {
		if half == FirstHalf && event.Minute > 45 {
			continue
		}
		if half == SecondHalf && event.Minute <= 45 && event.Minute > 90 {
			continue
		}

		if event.Type == GoalEvent || event.Type == InjuryEvent || event.Type == RedCardEvent {
			addedTime++
		}
	}
	return addedTime
}

func (m *Match) IsInAddedTime() bool {
	return (m.CurrentMinute > 45 && m.CurrentMinute <= 45+m.GetAddedTime(FirstHalf)) ||
		(m.CurrentMinute > 90 && m.CurrentMinute <= 90+m.GetAddedTime(SecondHalf))
}

func (m Match) IsHalfTime() bool {
	return m.CurrentHalf == FirstHalf && m.CurrentMinute > 45+m.GetAddedTime(FirstHalf)
}

func (m Match) IsFullTime() bool {
	return m.CurrentMinute >= 90+m.GetAddedTime(SecondHalf)
}

func (m *Match) GetMaxPlayerNameLength() int {
	maxNameLen := 0

	// Check home team
	for _, player := range m.Home.Players {
		if len(player.Player.Name) > maxNameLen {
			maxNameLen = len(player.Player.Name)
		}
	}

	// Check away team
	for _, player := range m.Away.Players {
		if len(player.Player.Name) > maxNameLen {
			maxNameLen = len(player.Player.Name)
		}
	}

	return maxNameLen
}

// ProgressBall attempts to move the ball to a better zone based on power difference
func (m *Match) ProgressBall(powerDiff int) {
	// Strong advantage - try to attack
	if powerDiff >= 10 {
		attacking := GetAttackingTransitions(m.ActiveZone)
		if len(attacking) > 0 {
			// Pick best attacking transition or random if multiple equal
			best := GetBestAttackingTransition(m.ActiveZone)
			if best != nil {
				m.ActiveZone = best.To
			}
		}
		return
	}

	// Moderate advantage - mix of lateral and forward
	if powerDiff >= 5 {
		allMoves := PitchTopology[m.ActiveZone]
		// Prefer forward, but allow lateral
		var validMoves []ZoneTransition
		for _, move := range allMoves {
			if move.IsForward || move.IsLateral {
				validMoves = append(validMoves, move)
			}
		}
		if len(validMoves) > 0 {
			// Weight towards forward moves
			forwardMoves := GetAttackingTransitions(m.ActiveZone)
			if len(forwardMoves) > 0 && rand.IntN(100) < 70 {
				m.ActiveZone = forwardMoves[rand.IntN(len(forwardMoves))].To
			} else if len(validMoves) > 0 {
				m.ActiveZone = validMoves[rand.IntN(len(validMoves))].To
			}
		}
		return
	}

	// Small advantage - mostly lateral, some backward
	lateralMoves := GetLateralTransitions(m.ActiveZone)
	if len(lateralMoves) > 0 && rand.IntN(100) < 60 {
		m.ActiveZone = lateralMoves[rand.IntN(len(lateralMoves))].To
	}
}

func (m *Match) PlayPhase() PhaseResult {
	homeRoll := rand.IntN(20)
	awayRoll := rand.IntN(20)

	homePhasePower := m.Home.Club.Strength + homeRoll
	awayPhasePower := m.Away.Club.Strength + awayRoll
	powerDiff := int(math.Abs(float64(homePhasePower - awayPhasePower)))

	morePowerfulTeam := m.Home
	if homePhasePower < awayPhasePower {
		morePowerfulTeam = m.Away
	}

	// Possibility to change possession
	if morePowerfulTeam != m.TeamInPossession {
		m.AddEvent(NewEvent(PossessionChangedEvent, m.CurrentMinute, morePowerfulTeam, nil))
		m.TeamInPossession = morePowerfulTeam

		// When possession changes, ball likely moves backward for new team
		defensiveMoves := GetDefensiveTransitions(m.ActiveZone)
		if len(defensiveMoves) > 0 {
			m.ActiveZone = defensiveMoves[rand.IntN(len(defensiveMoves))].To
		}

		return PhaseResult{
			HomeRoll:       homeRoll,
			AwayRoll:       awayRoll,
			HomePhasePower: homePhasePower,
			AwayPhasePower: awayPhasePower,
			HomeGoals:      0,
			AwayGoals:      0,
		}
	}

	// Team kept the ball, try to progress it
	m.ProgressBall(powerDiff)

	// Team kept the ball, no action taken
	if powerDiff < shotThreshold {
		return PhaseResult{
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
			m.AddEvent(NewEvent(SavedShotEvent, m.CurrentMinute, morePowerfulTeam, nil))
		} else {
			m.AddEvent(NewEvent(MissedShotEvent, m.CurrentMinute, morePowerfulTeam, nil))
		}
		return PhaseResult{
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
			scorer := m.Home.GetRandomOutfielder()
			m.AddEvent(NewEvent(GoalEvent, m.CurrentMinute, m.Home, &scorer))
		}

		if homePhasePower < awayPhasePower {
			awayGoals = goalsThisPhase
			scorer := m.Away.GetRandomOutfielder()
			m.AddEvent(NewEvent(GoalEvent, m.CurrentMinute, m.Away, &scorer))
		}
	}

	return PhaseResult{
		HomeRoll:       homeRoll,
		AwayRoll:       awayRoll,
		HomePhasePower: homePhasePower,
		AwayPhasePower: awayPhasePower,
		HomeGoals:      homeGoals,
		AwayGoals:      awayGoals,
	}
}
