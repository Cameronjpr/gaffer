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

type Team struct {
	Name     string
	Strength int
	Score    int
}

type Match struct {
	Home         Team
	Away         Team
	CurrentPhase int
	CurrentHalf  int
	IsComplete   bool
	IsHalfTime   bool
	PhaseHistory []PhaseResult
	Commentary   []CommentaryMessage
}

type CommentaryMessage struct {
	Message string
	Flash   bool
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

func createTeam(name string) Team {
	return Team{
		Name:     name,
		Strength: rand.IntN(20),
		Score:    0,
	}
}

func NewMatch() Match {
	return Match{
		Home:         createTeam("Home"),
		Away:         createTeam("Away"),
		CurrentPhase: 1,
		CurrentHalf:  1,
		IsComplete:   false,
		IsHalfTime:   false,
		PhaseHistory: make([]PhaseResult, 0),
		Commentary:   make([]CommentaryMessage, 0),
	}
}

func PlayMatch() {
	m := NewMatch()

	file, err := os.Create("game.log")
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer file.Close()

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

	homePhaseStrength := m.Home.Strength + homeRoll
	awayPhaseStrength := m.Away.Strength + awayRoll

	powerDiff := math.Abs(float64(homePhaseStrength - awayPhaseStrength))
	goalsThisPhase := 0

	if powerDiff < 80 {
		m.Commentary = append(m.Commentary, CommentaryMessage{Message: "...", Flash: false})
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
			m.Commentary = append(m.Commentary, CommentaryMessage{Message: "Home team scores!", Flash: true})
		}

		if homePhaseStrength < awayPhaseStrength {
			awayGoals = goalsThisPhase
			m.Commentary = append(m.Commentary, CommentaryMessage{Message: "Away team scores!", Flash: true})
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

	fmt.Fprintf(file, "Home [%v] – [%v] Away\n", match.Home.Score, match.Away.Score)
	fmt.Fprintf(file, "Phase %v\n", match.CurrentPhase)
}
