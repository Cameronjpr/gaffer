package game

import (
	"math"
	"math/rand/v2"
)

var (
	phases = 90
)

const goalscoringThreshold = 75

type Half int

const (
	FirstHalf Half = iota + 1
	SecondHalf
)

type Match struct {
	Home          *MatchParticipant
	Away          *MatchParticipant
	CurrentMinute int
	CurrentHalf   Half
	PhaseHistory  []PhaseResult
	Events        []Event
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
	return Match{
		Home:          NewMatchParticipant(homeClub),
		Away:          NewMatchParticipant(awayClub),
		CurrentMinute: 1,
		CurrentHalf:   FirstHalf,
		PhaseHistory:  make([]PhaseResult, 0),
		Events:        make([]Event, 0),
	}
}

func (m *Match) StartFirstHalf() {
	m.CurrentHalf = FirstHalf
	m.CurrentMinute = 1
}

func (m *Match) StartSecondHalf() {
	m.CurrentHalf = SecondHalf
	m.CurrentMinute = 46
}

// AddEvent adds an event to the match
func (m *Match) AddEvent(event Event) {
	m.Events = append(m.Events, event)
}

func (m *Match) GetAddedTime(half Half) int {
	addedTime := 0
	for _, event := range m.Events {
		if event.Minute > 45 && half == SecondHalf {
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

func (m *Match) PlayPhase() PhaseResult {
	homeRoll := rand.IntN(100)
	awayRoll := rand.IntN(100)

	homePhasePower := m.Home.Club.Strength + homeRoll
	awayPhasePower := m.Away.Club.Strength + awayRoll
	powerDiff := math.Abs(float64(homePhasePower - awayPhasePower))

	if homePhasePower > awayPhasePower {
		if m.Home.HasPossession {
			m.AddEvent(NewEvent(PossessionRetainedEvent, m.CurrentMinute, m.Home, nil))
		} else {
			m.Home.WinPossession()
			m.Away.LosePossession()
			m.AddEvent(NewEvent(PossessionChangedEvent, m.CurrentMinute, m.Home, nil))
		}
	} else if homePhasePower < awayPhasePower {
		if m.Away.HasPossession {
			m.AddEvent(NewEvent(PossessionRetainedEvent, m.CurrentMinute, m.Away, nil))
		} else {
			m.Away.WinPossession()
			m.Home.LosePossession()
			m.AddEvent(NewEvent(PossessionChangedEvent, m.CurrentMinute, m.Away, nil))
		}
	}

	if powerDiff < goalscoringThreshold {
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

func (m Match) IsHalfTime() bool {
	return m.CurrentMinute == 45+m.GetAddedTime(FirstHalf)
}

func (m Match) IsFullTime() bool {
	return m.CurrentMinute == 90+m.GetAddedTime(SecondHalf)
}
