package game

import (
	"fmt"
	"math"
	"math/rand/v2"
	"os"
)

var (
	phases = 90
)

const goalscoringThreshold = 75

// MatchPlayerParticipant represents a player participating in a specific match
type MatchPlayerParticipant struct {
	Player   *Player
	Position string
}

// MatchParticipant represents a club participating in a specific match
type MatchParticipant struct {
	Club          *Club
	Players       []MatchPlayerParticipant
	Formation     string
	Score         int
	HasPossession bool
}

// NewMatchParticipant creates a new match participant from a club
func NewMatchParticipant(club *Club) *MatchParticipant {
	// Assign players to 4-3-3 formation positions
	// Positions: GK, RB, CB, CB, LB, CM, CM, CM, RW, ST, LW
	positions := []string{"GK", "RB", "CB", "CB", "LB", "CM", "CM", "CM", "RW", "ST", "LW"}

	matchPlayers := make([]MatchPlayerParticipant, 0, len(club.Players))
	for i := range club.Players {
		if i < len(positions) {
			matchPlayers = append(matchPlayers, MatchPlayerParticipant{
				Player:   &club.Players[i],
				Position: positions[i],
			})
		}
	}

	return &MatchParticipant{
		Club:          club,
		Players:       matchPlayers,
		Formation:     "4-3-3",
		Score:         0,
		HasPossession: false,
	}
}

func (p *MatchParticipant) GetStarPlayers() []MatchPlayerParticipant {
	stars := make([]MatchPlayerParticipant, 0)
	highestQuality := 0
	for _, player := range p.Players {
		if player.Player.Quality > highestQuality {
			highestQuality = player.Player.Quality
		}
	}
	for _, player := range p.Players {
		if player.Player.Quality == highestQuality {
			stars = append(stars, player)
		}
	}
	return stars
}

func (p *MatchParticipant) GetRandomOutfielder() MatchPlayerParticipant {
	outfielders := make([]MatchPlayerParticipant, 0)
	for _, player := range p.Players {
		if player.Position != "GK" {
			outfielders = append(outfielders, player)
		}
	}
	if len(outfielders) == 0 {
		return MatchPlayerParticipant{}
	}
	return outfielders[rand.IntN(len(outfielders))]
}

func (p *MatchParticipant) GetLineup() string {
	lineup := ""
	stars := p.GetStarPlayers()
	for _, matchPlayer := range p.Players {
		suffix := ""
		for _, star := range stars {
			if star.Player.Name == matchPlayer.Player.Name {
				suffix += " ★"
			}
		}

		row := fmt.Sprintf("%s - %s%s\n", matchPlayer.Position, matchPlayer.Player.Name, suffix)
		lineup += row

	}
	return lineup
}

func (p *MatchParticipant) AddPlayer(matchPlayer MatchPlayerParticipant) {
	p.Players = append(p.Players, matchPlayer)
}

func (p *MatchParticipant) IncreaseScore() {
	p.Score++
}

func (p *MatchParticipant) WinPossession() {
	p.HasPossession = true
}

func (p *MatchParticipant) LosePossession() {
	p.HasPossession = false
}

type Match struct {
	Home         *MatchParticipant
	Away         *MatchParticipant
	CurrentPhase int
	CurrentHalf  int
	IsComplete   bool
	PhaseHistory []PhaseResult
	Events       []Event
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
		Home:         NewMatchParticipant(homeClub),
		Away:         NewMatchParticipant(awayClub),
		CurrentPhase: 1,
		CurrentHalf:  1,
		PhaseHistory: make([]PhaseResult, 0),
		Events:       make([]Event, 0),
	}
}

// AddEvent adds an event to the match
func (m *Match) AddEvent(event Event) {
	m.Events = append(m.Events, event)
}

func (m *Match) PlayPhase() PhaseResult {
	homeRoll := rand.IntN(100)
	awayRoll := rand.IntN(100)

	homePhasePower := m.Home.Club.Strength + homeRoll
	awayPhasePower := m.Away.Club.Strength + awayRoll
	powerDiff := math.Abs(float64(homePhasePower - awayPhasePower))

	if homePhasePower > awayPhasePower {
		if m.Home.HasPossession {
			m.AddEvent(NewEvent(PossessionRetainedEvent, m.CurrentPhase, m.Home, nil))
		} else {
			m.Home.WinPossession()
			m.Away.LosePossession()
			m.AddEvent(NewEvent(PossessionChangedEvent, m.CurrentPhase, m.Home, nil))
		}
	} else if homePhasePower < awayPhasePower {
		if m.Away.HasPossession {
			m.AddEvent(NewEvent(PossessionRetainedEvent, m.CurrentPhase, m.Away, nil))
		} else {
			m.Away.WinPossession()
			m.Home.LosePossession()
			m.AddEvent(NewEvent(PossessionChangedEvent, m.CurrentPhase, m.Away, nil))
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
			m.AddEvent(NewEvent(GoalEvent, m.CurrentPhase, m.Home, &scorer))
		}

		if homePhasePower < awayPhasePower {
			awayGoals = goalsThisPhase
			scorer := m.Away.GetRandomOutfielder()
			m.AddEvent(NewEvent(GoalEvent, m.CurrentPhase, m.Away, &scorer))
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
	return m.CurrentPhase == 45
}

func (m Match) IsFullTime() bool {
	return m.CurrentPhase == 90
}

func writeMatchLog(match Match) {
	file, err := os.Create("game.log")
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer file.Close()

	fmt.Fprintf(file, "%v [%v] – [%v] %v\n", match.Home.Club.Name, match.Home.Score, match.Away.Score, match.Away.Club.Name)
	fmt.Fprintf(file, "Phase %v\n", match.CurrentPhase)
}
