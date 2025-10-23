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
	Home             *MatchParticipant
	Away             *MatchParticipant
	TeamInPossession *MatchParticipant
	CurrentMinute    int
	CurrentHalf      Half
	ActiveZone       PitchZone
	PhaseHistory     []PhaseResult
	Events           []Event
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
