package domain

import (
	"fmt"
	"math/rand/v2"
	"strings"
)

// MatchPlayerParticipant represents a player participating in a specific match
type MatchPlayerParticipant struct {
	Player   *Player
	Position string
	Stamina  int
}

// MatchParticipant represents a club participating in a specific match
type MatchParticipant struct {
	Club      *Club
	Players   []MatchPlayerParticipant
	Formation string
	Score     int
}

// NewMatchParticipant creates a new match participant from a club and its players
func NewMatchParticipant(club *Club, players []Player) *MatchParticipant {
	// Assign players to 4-3-3 formation positions
	// Positions: GK, RB, CB, CB, LB, CM, CM, CM, RW, ST, LW
	positions := []string{"GK", "RB", "CB", "CB", "LB", "CM", "CM", "CM", "RW", "ST", "LW"}

	matchPlayers := make([]MatchPlayerParticipant, 0, len(players))
	for i := range players {
		if i < len(positions) {
			matchPlayers = append(matchPlayers, MatchPlayerParticipant{
				Player:   &players[i],
				Position: positions[i],
			})
		}
	}

	return &MatchParticipant{
		Club:      club,
		Players:   matchPlayers,
		Formation: "4-3-3",
		Score:     0,
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

func (p *MatchParticipant) GetLineup(match *Match) string {
	lineup := ""
	stars := p.GetStarPlayers()

	// Calculate max player name length for consistent formatting
	// If we have a match, use global max across both teams to prevent layout shift
	maxNameLen := 0
	if match != nil {
		maxNameLen = match.GetMaxPlayerNameLength()
	} else {
		// Pre-match: calculate for this team only
		for _, player := range p.Players {
			if len(player.Player.Name) > maxNameLen {
				maxNameLen = len(player.Player.Name)
			}
		}
	}

	for _, matchPlayer := range p.Players {
		// Star indicator with brackets
		starIndicator := "   " // 3 spaces
		for _, star := range stars {
			if star.Player.Name == matchPlayer.Player.Name {
				starIndicator = "[★]"
				break
			}
		}

		// Count goals for this player
		goalCount := 0
		if match != nil {
			for _, event := range match.Events {
				if event.Type == GoalEvent &&
					event.For == p &&
					event.Player != nil &&
					event.Player.Player != nil &&
					event.Player.Player.Name == matchPlayer.Player.Name {
					goalCount++
				}
			}
		}

		// Build goal indicator string
		goalIndicator := ""
		if goalCount > 0 {
			goalIndicator = " " + strings.Repeat("●", goalCount)
		}

		// Fixed-width formatting: position (2 chars) - name (padded) star (3 chars) goals
		row := fmt.Sprintf("%-2s - %-*s %s%s\n",
			matchPlayer.Position,
			maxNameLen,
			matchPlayer.Player.Name,
			starIndicator,
			goalIndicator)
		lineup += row
	}
	return lineup
}

func (p *MatchParticipant) DrainStamina(hasPossession bool) {
	amount := 1
	if !hasPossession {
		amount = 2
	}
	for _, player := range p.Players {
		player.DrainStamina(amount)
	}

}

func (p *MatchPlayerParticipant) GetStamina() int {
	return p.Stamina
}

func (p *MatchPlayerParticipant) DrainStamina(amount int) {
	p.Stamina -= amount
	if p.Stamina < 0 {
		p.Stamina = 0
	}
}
