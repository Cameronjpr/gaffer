package domain

import "fmt"

// Club represents a football club with permanent attributes
type Club struct {
	Name       string
	Strength   int // out of 20
	Background string
	Foreground string
	Players    []Player
}

var Clubs = []Club{
	{
		Name:       "Arsenal",
		Strength:   20,
		Background: "#DB0007",
		Foreground: "#FFFFFF",
		Players: []Player{
			{Name: "Raya", Quality: 18},
			{Name: "Timber", Quality: 17},
			{Name: "Saliba", Quality: 18},
			{Name: "Gabriel", Quality: 18},
			{Name: "Calafiori", Quality: 17},
			{Name: "Zubimendi", Quality: 18},
			{Name: "Rice", Quality: 19},
			{Name: "Ødegaard", Quality: 18},
			{Name: "Saka", Quality: 19},
			{Name: "Gyokeres", Quality: 17},
			{Name: "Trossard", Quality: 17},
		},
	},
	{
		Name:       "Manchester City",
		Strength:   19,
		Background: "#6CABDD",
		Foreground: "#1C2C5B",
		Players: []Player{
			{Name: "Donnarumma", Quality: 18},
			{Name: "Lewis", Quality: 15},
			{Name: "Stones", Quality: 17},
			{Name: "Ruben Dias", Quality: 18},
			{Name: "Gvardiol", Quality: 17},
			{Name: "González", Quality: 18},
			{Name: "M. Nunes", Quality: 17},
			{Name: "B. Silva", Quality: 15},
			{Name: "Savinho", Quality: 18},
			{Name: "Haaland", Quality: 19},
			{Name: "Doku", Quality: 16},
		},
	},
}

// GetClubByName returns a pointer to a club by name, or nil if not found
func GetClubByName(name string) *Club {
	for i := range Clubs {
		if Clubs[i].Name == name {
			return &Clubs[i]
		}
	}
	return nil
}

func (c *Club) GetSquad() string {
	lineup := ""
	for _, player := range c.Players {
		lineup += fmt.Sprintf("%s (Q:%d)\n", player.Name, player.Quality)
	}
	return lineup
}
