package domain

import (
	"fmt"
)

// Club represents a football club with permanent attributes
type Club struct {
	ID         int64
	Name       string
	Strength   int // out of 20
	Background string
	Foreground string
}

// ClubWithPlayers is a view model for when you need club + players together
type ClubWithPlayers struct {
	Club    *Club
	Players []Player
}

// GetSquad returns a formatted string of the squad (for ClubWithPlayers)
func (cwp *ClubWithPlayers) GetSquad() string {
	lineup := ""
	for _, player := range cwp.Players {
		lineup += fmt.Sprintf("%s (Q:%d)\n", player.Name, player.Quality)
	}
	return lineup
}
