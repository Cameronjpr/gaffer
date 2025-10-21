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

// MatchParticipant represents a club participating in a specific match
type MatchParticipant struct {
	Club          *Club
	Players       []string
	Score         int
	HasPossession bool
}

// NewMatchParticipant creates a new match participant from a club
func NewMatchParticipant(club *Club) *MatchParticipant {
	return &MatchParticipant{
		Club:          club,
		Players:       []string{},
		Score:         0,
		HasPossession: false,
	}
}

func (p *MatchParticipant) GetLineup() string {
	lineup := ""
	for _, player := range p.Players {
		lineup += fmt.Sprintf("%s\n", player)
	}
	return lineup
}

func (p *MatchParticipant) AddPlayer(player string) {
	p.Players = append(p.Players, player)
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
	IsHalfTime   bool
	PhaseHistory []PhaseResult
	Commentary   []CommentaryMessage
}

type PhaseResult struct {
	Phase             int
	HomeRoll          int
	AwayRoll          int
	HomePhaseStrength int
	AwayPhaseStrength int
	HomeGoals         int
	AwayGoals         int
}

func NewMatch(homeClub, awayClub *Club) Match {
	return Match{
		Home:         NewMatchParticipant(homeClub),
		Away:         NewMatchParticipant(awayClub),
		CurrentPhase: 1,
		CurrentHalf:  1,
		IsComplete:   false,
		IsHalfTime:   false,
		PhaseHistory: make([]PhaseResult, 0),
		Commentary:   make([]CommentaryMessage, 0),
	}
}

func PlayMatch() {
	// Example match for testing
	homeClub := GetClubByName("Leeds United")
	awayClub := GetClubByName("Arsenal")
	m := NewMatch(homeClub, awayClub)

	file, err := os.Create("game.log")
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer file.Close()

	coinToss := rand.IntN(2) == 0
	if coinToss {
		m.Home.HasPossession = true
	} else {
		m.Away.HasPossession = true
	}

	for m.CurrentPhase <= 45 {
		m.PlayPhase()
		m.CurrentPhase++
	}

	m.CurrentHalf++

	for m.CurrentPhase <= 90 {
		m.PlayPhase()
		m.CurrentPhase++
	}

	writeMatchLog(m)

}

func (m *Match) PlayPhase() PhaseResult {
	homeRoll := rand.IntN(100)
	awayRoll := rand.IntN(100)

	homePhaseStrength := m.Home.Club.Strength + homeRoll
	awayPhaseStrength := m.Away.Club.Strength + awayRoll

	powerDiff := math.Abs(float64(homePhaseStrength - awayPhaseStrength))
	goalsThisPhase := 0

	if homePhaseStrength > awayPhaseStrength {

		if m.Home.HasPossession {
			m.Commentary = append(m.Commentary, getCommentaryForEvent(PossessionRetainedEvent, m.Home, m))
		} else {
			m.Home.WinPossession()
			m.Away.LosePossession()
			m.Commentary = append(m.Commentary, getCommentaryForEvent(PossessionChangedEvent, m.Home, m))
		}

	} else if homePhaseStrength < awayPhaseStrength {
		if m.Away.HasPossession {
			m.Commentary = append(m.Commentary, getCommentaryForEvent(PossessionRetainedEvent, m.Away, m))
		} else {
			m.Away.WinPossession()
			m.Home.LosePossession()
			m.Commentary = append(m.Commentary, getCommentaryForEvent(PossessionChangedEvent, m.Away, m))
		}
	}

	if powerDiff < 80 {
		return PhaseResult{
			HomeRoll:          homeRoll,
			AwayRoll:          awayRoll,
			HomePhaseStrength: homePhaseStrength,
			AwayPhaseStrength: awayPhaseStrength,
			HomeGoals:         0,
			AwayGoals:         0,
		}
	}
	if powerDiff > 80 {
		goalsThisPhase = rand.IntN(2)
	}

	homeGoals := 0
	awayGoals := 0

	if goalsThisPhase > 0 {
		if homePhaseStrength > awayPhaseStrength {
			homeGoals = goalsThisPhase
			m.Commentary = append(m.Commentary, getCommentaryForEvent(HomeGoalScoredEvent, m.Home, m))
		}

		if homePhaseStrength < awayPhaseStrength {
			awayGoals = goalsThisPhase
			m.Commentary = append(m.Commentary, getCommentaryForEvent(AwayGoalScoredEvent, m.Away, m))
		}
	}

	return PhaseResult{
		HomeRoll:          homeRoll,
		AwayRoll:          awayRoll,
		HomePhaseStrength: homePhaseStrength,
		AwayPhaseStrength: awayPhaseStrength,
		HomeGoals:         homeGoals,
		AwayGoals:         awayGoals,
	}
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
