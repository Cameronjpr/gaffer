package game

// Club represents a football club with permanent attributes
type Club struct {
	Name     string
	Strength int // out of 20
}

// Clubs contains all Premier League clubs with realistic 2025/26 strength values
var Clubs = []Club{
	{Name: "Arsenal", Strength: 20},
	{Name: "Manchester City", Strength: 19},
	{Name: "Bournemouth", Strength: 16},
	{Name: "Liverpool", Strength: 18},
	{Name: "Chelsea", Strength: 17},
	{Name: "Tottenham", Strength: 17},
	{Name: "Sunderland", Strength: 15},
	{Name: "Crystal Palace", Strength: 15},
	{Name: "Manchester United", Strength: 14},
	{Name: "Brighton", Strength: 14},
	{Name: "Aston Villa", Strength: 13},
	{Name: "Newcastle United", Strength: 13},
	{Name: "West Ham", Strength: 12},
	{Name: "Everton", Strength: 12},
	{Name: "Fulham", Strength: 12},
	{Name: "Brentford", Strength: 11},
	{Name: "Burnley", Strength: 11},
	{Name: "Leeds United", Strength: 10},
	{Name: "Wolves", Strength: 10},
	{Name: "Nottingham Forest", Strength: 9},
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
