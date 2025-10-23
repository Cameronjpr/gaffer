package domain

import "fmt"

type Fixture struct {
	ID       int
	HomeTeam *Club
	AwayTeam *Club
	Result   Match
}

type Gameweek struct {
	ID       int
	Name     string
	Fixtures []Fixture
}

type Season struct {
	ID          int
	Name        string
	Clubs       []*Club
	Gameweeks   []Gameweek
	LeagueTable *LeagueTable
}

type LeaguePosition struct {
	Club   *Club
	Points int
}

type LeagueTable struct {
	Positions []LeaguePosition
}

func NewSeason(clubs []*Club) *Season {
	season := Season{
		ID:        1,
		Name:      "2025/26",
		Clubs:     clubs,
		Gameweeks: []Gameweek{},
		LeagueTable: &LeagueTable{
			Positions: []LeaguePosition{},
		},
	}

	return &season
}

func (s *Season) GenerateGameweeks() {
	for i := 1; i <= 38; i++ {
		gameweek := Gameweek{
			ID:       i,
			Name:     fmt.Sprintf("Gameweek %d", i),
			Fixtures: generateFixtures(s.Clubs),
		}
		s.Gameweeks = append(s.Gameweeks, gameweek)
	}
}

func generateFixtures(clubs []*Club) []Fixture {
	return []Fixture{
		{
			ID:       1,
			HomeTeam: clubs[0],
			AwayTeam: clubs[1],
			Result:   Match{},
		},
		{
			ID:       2,
			HomeTeam: clubs[1],
			AwayTeam: clubs[0],
			Result:   Match{},
		},
	}
}

func (s *Season) GenerateLeagueTable() {
	for _, club := range s.Clubs {
		position := LeaguePosition{
			Club:   club,
			Points: 0,
		}
		s.LeagueTable.Positions = append(s.LeagueTable.Positions, position)
	}
}

// GetNextFixture finds the next fixture in the season that hasn't been played.
func (s *Season) GetNextFixture() (*Fixture, error) {
	for i := range s.Gameweeks {
		for j := range s.Gameweeks[i].Fixtures {
			// A fixture is considered unplayed if its Result has no home team.
			if s.Gameweeks[i].Fixtures[j].Result.Home == nil {
				return &s.Gameweeks[i].Fixtures[j], nil
			}
		}
	}
	return nil, fmt.Errorf("no unplayed fixtures remaining in the season")
}
