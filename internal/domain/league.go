package domain

type LeaguePosition struct {
	Club           *Club
	Played         int
	Won            int
	Drawn          int
	Lost           int
	GoalsFor       int
	GoalsAgainst   int
	GoalDifference int
	Points         int
}

type LeagueTable struct {
	Positions []LeaguePosition
}

type ByLeagueStanding []LeaguePosition

func (s ByLeagueStanding) Len() int      { return len(s) }
func (s ByLeagueStanding) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
func (s ByLeagueStanding) Less(i, j int) bool {
	// Points (descending - more points = better)
	if s[i].Points != s[j].Points {
		return s[i].Points > s[j].Points
	}
	// Goal difference (descending)
	if s[i].GoalDifference != s[j].GoalDifference {
		return s[i].GoalDifference > s[j].GoalDifference
	}
	// Goals scored (descending)
	if s[i].GoalsFor != s[j].GoalsFor {
		return s[i].GoalsFor > s[j].GoalsFor
	}
	// Alphabetical (ascending - for deterministic ordering)
	return s[i].Club.Name < s[j].Club.Name
}
