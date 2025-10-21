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

func NewMatch() Match {
	return Match{
		Home:         createTeam("Leeds"),
		Away:         createTeam("Arsenal"),
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

	homePhaseStrength := m.Home.Strength + homeRoll
	awayPhaseStrength := m.Away.Strength + awayRoll

	powerDiff := math.Abs(float64(homePhaseStrength - awayPhaseStrength))
	goalsThisPhase := 0

	if homePhaseStrength > awayPhaseStrength {

		if m.Home.HasPossession {
			m.Commentary = append(m.Commentary, CommentaryMessage{Message: fmt.Sprintf("%v are keeping the ball", m.Home.Name), Flash: false})
		} else {
			m.Home.WinPossession()
			m.Away.LosePossession()
			m.Commentary = append(m.Commentary, CommentaryMessage{Message: fmt.Sprintf("%v takes the ball!", m.Home.Name), Flash: false})
		}

	} else if homePhaseStrength < awayPhaseStrength {
		if m.Away.HasPossession {
			m.Commentary = append(m.Commentary, CommentaryMessage{Message: fmt.Sprintf("%v are keeping the ball", m.Away.Name), Flash: false})
		} else {
			m.Away.WinPossession()
			m.Home.LosePossession()
			m.Commentary = append(m.Commentary, CommentaryMessage{Message: fmt.Sprintf("%v takes the ball!", m.Away.Name), Flash: false})
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
			m.Commentary = append(m.Commentary, CommentaryMessage{Message: fmt.Sprintf("GOAL! %v score!", m.Home.Name), Flash: true})
		}

		if homePhaseStrength < awayPhaseStrength {
			awayGoals = goalsThisPhase
			m.Commentary = append(m.Commentary, CommentaryMessage{Message: fmt.Sprintf("GOAL! %v score!", m.Away.Name), Flash: true})
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

	fmt.Fprintf(file, "%v [%v] – [%v] %v\n", match.Home.Name, match.Home.Score, match.Away.Score, match.Away.Name)
	fmt.Fprintf(file, "Phase %v\n", match.CurrentPhase)
}
