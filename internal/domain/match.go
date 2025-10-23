package domain

type Half int

const (
	FirstHalf Half = iota + 1
	SecondHalf
)

type PhaseResult struct {
	Phase          int
	HomeRoll       int
	AwayRoll       int
	HomePhasePower int
	AwayPhasePower int
	HomeGoals      int
	AwayGoals      int
}

type Match struct {
	ForFixture             *Fixture
	Home                   *MatchParticipant
	Away                   *MatchParticipant
	TeamInPossession       *MatchParticipant
	CurrentMinute          int
	CurrentHalf            Half
	ActiveZone             PitchZone
	HomeAttackingDirection AttackingDirection // Which goal Home attacks (switches at halftime)
	PhaseHistory           []PhaseResult
	Events                 []Event
}

func NewMatchFromFixture(f *Fixture) *Match {
	home := NewMatchParticipant(f.HomeTeam)
	away := NewMatchParticipant(f.AwayTeam)
	return &Match{
		ForFixture:             f,
		Home:                   home,
		Away:                   away,
		TeamInPossession:       home, // Home team starts with kickoff
		CurrentMinute:          1,
		CurrentHalf:            FirstHalf,
		ActiveZone:             WestMidCentre, // Kickoff from center
		HomeAttackingDirection: AttackingEast, // Home attacks East in first half
		PhaseHistory:           make([]PhaseResult, 0),
		Events:                 make([]Event, 0),
	}
}

func (m *Match) StartFirstHalf() {
	m.CurrentHalf = FirstHalf
	m.CurrentMinute = 1
	m.HomeAttackingDirection = AttackingEast // Home attacks East in first half
}

func (m *Match) StartSecondHalf() {
	m.CurrentHalf = SecondHalf
	m.CurrentMinute = 45                     // Will be incremented to 46 on first tick
	m.HomeAttackingDirection = AttackingWest // Teams switch sides at halftime
}

func (m *Match) AddEvent(event Event) {
	m.Events = append(m.Events, event)
}

// GetAttackingDirection returns the attacking direction for the team in possession
func (m *Match) GetAttackingDirection() AttackingDirection {
	if m.TeamInPossession == m.Home {
		return m.HomeAttackingDirection
	}
	// Away team attacks opposite direction
	if m.HomeAttackingDirection == AttackingEast {
		return AttackingWest
	}
	return AttackingEast
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

func (m *Match) IsFirstHalf() bool {
	return m.CurrentHalf == FirstHalf
}

func (m *Match) IsSecondHalf() bool {
	return m.CurrentHalf == SecondHalf
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

// GetScore returns the current score by counting goal events.
// Returns (homeScore, awayScore).
func (m *Match) GetScore() (int, int) {
	homeScore := 0
	awayScore := 0

	for _, event := range m.Events {
		if event.Type == GoalEvent {
			if event.For == m.Home {
				homeScore++
			} else if event.For == m.Away {
				awayScore++
			}
		}
	}

	return homeScore, awayScore
}

func (m *Match) GetWinner() *Club {
	// Safety check - should never happen, but prevents crash
	if m == nil || m.Home == nil || m.Away == nil {
		return nil
	}

	if m.Home.Score > m.Away.Score {
		return m.Home.Club
	} else if m.Away.Score > m.Home.Score {
		return m.Away.Club
	}
	return nil
}
