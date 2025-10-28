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
	CurrentXI []*MatchPlayerParticipant
	Bench     []*MatchPlayerParticipant
	Formation string
	Score     int
}

// NewMatchParticipant creates a new match participant from a club and its players
func NewMatchParticipant(club *Club, players []Player) *MatchParticipant {
	// Assign players to 4-3-3 formation positions
	// Positions: GK, RB, CB, CB, LB, CM, CM, CM, RW, ST, LW
	positions := []string{"GK", "RB", "CB", "CB", "LB", "CM", "CM", "CM", "RW", "ST", "LW"}

	// Split into starting XI and bench
	currentXI := make([]*MatchPlayerParticipant, 0, 11)
	bench := make([]*MatchPlayerParticipant, 0, 7)

	for i := range players {
		if i < len(positions) {
			// First 11 players are the starting XI
			currentXI = append(currentXI, &MatchPlayerParticipant{
				Player:   &players[i],
				Position: positions[i],
			})
		} else {
			// Remaining players go on the bench (no specific position)
			bench = append(bench, &MatchPlayerParticipant{
				Player:   &players[i],
				Position: "", // Bench players don't have assigned positions
			})
		}
	}

	return &MatchParticipant{
		Club:      club,
		CurrentXI: currentXI,
		Bench:     bench,
		Formation: "4-3-3",
		Score:     0,
	}
}

func (p *MatchParticipant) MakeSubstitution(in, out *MatchPlayerParticipant) {
	// Find and replace the player coming off in CurrentXI
	for i, player := range p.CurrentXI {
		if player == out {
			// Give the substitute the position of the player coming off
			in.Position = out.Position
			p.CurrentXI[i] = in
			break
		}
	}

	// Remove the substitute from the bench
	for i, player := range p.Bench {
		if player == in {
			p.Bench = append(p.Bench[:i], p.Bench[i+1:]...)
			break
		}
	}

	// Add the substituted player to the bench (clear their position)
	out.Position = ""
	p.Bench = append(p.Bench, out)
}

func (p *MatchParticipant) GetStarPlayers() []*MatchPlayerParticipant {
	stars := make([]*MatchPlayerParticipant, 0)
	highestQuality := 0
	for _, player := range p.CurrentXI {
		if player.Player.Quality > highestQuality {
			highestQuality = player.Player.Quality
		}
	}
	for _, player := range p.CurrentXI {
		if player.Player.Quality == highestQuality {
			stars = append(stars, player)
		}
	}
	return stars
}

func (p *MatchParticipant) GetAverageQuality() float64 {
	totalQuality := 0
	for _, player := range p.CurrentXI {
		totalQuality += player.Player.Quality
	}
	return float64(totalQuality) / float64(len(p.CurrentXI))
}

func (p *MatchParticipant) GetRandomOutfielder() *MatchPlayerParticipant {
	outfielders := make([]*MatchPlayerParticipant, 0)
	for _, player := range p.CurrentXI {
		if player.Position != "GK" {
			outfielders = append(outfielders, player)
		}
	}
	if len(outfielders) == 0 {
		return nil
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
		for _, player := range p.CurrentXI {
			if len(player.Player.Name) > maxNameLen {
				maxNameLen = len(player.Player.Name)
			}
		}
	}

	for _, matchPlayer := range p.CurrentXI {
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
	for _, player := range p.CurrentXI {
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
