package game

import "math/rand/v2"

type Team struct {
	Name          string
	Strength      int
	Score         int
	HasPossession bool
}

func createTeam(name string) Team {
	return Team{
		Name:          name,
		Strength:      rand.IntN(20),
		Score:         0,
		HasPossession: false,
	}
}

func (t *Team) IncreaseScore() {
	t.Score++
}

func (t *Team) WinPossession() {
	t.HasPossession = true
}

func (t *Team) LosePossession() {
	t.HasPossession = false
}
