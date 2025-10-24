package domain

import (
	"fmt"
	"sort"
)

type Gameweek struct {
	ID       int
	Name     string
	Fixtures []*Fixture
}

type Season struct {
	ID        int
	Name      string
	Clubs     []*Club
	Gameweeks []*Gameweek
}

func NewSeason(clubs []*Club) *Season {

	season := Season{
		ID:        1,
		Name:      "2025/26",
		Clubs:     clubs,
		Gameweeks: []*Gameweek{},
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
		s.Gameweeks = append(s.Gameweeks, &gameweek)
	}
}

func generateFixtures(clubs []*Club) []*Fixture {
	var fixtures []*Fixture

	for _, club := range clubs {
		for _, opponent := range clubs {
			if club == opponent {
				continue
			}
			fixture := Fixture{
				ID:       len(clubs)*len(clubs) + 1,
				HomeTeam: club,
				AwayTeam: opponent,
				Result:   nil,
			}
			fixtures = append(fixtures, &fixture)
		}
	}
	return fixtures
}

func (s *Season) GetLeagueTable() LeagueTable {
	table := LeagueTable{
		Positions: []LeaguePosition{},
	}

	for _, club := range s.Clubs {
		fixtures := s.GetFixturesForClub(club)

		var points int
		for _, fixture := range fixtures {
			// Skip fixtures that haven't been played yet
			if fixture.Result == nil {
				continue
			}

			winner := fixture.Result.GetWinner()
			if winner == nil {
				points += 1 // Draw
			} else if winner == club {
				points += 3 // Win
			}
		}

		position := LeaguePosition{
			Club:   club,
			Points: points,
		}
		table.Positions = append(table.Positions, position)
	}

	sort.Sort(ByLeagueStanding(table.Positions))

	return table
}

func (s *Season) GetFixturesForClub(club *Club) []*Fixture {
	var fixtures []*Fixture
	for _, gameweek := range s.Gameweeks {
		for _, fixture := range gameweek.Fixtures {
			if fixture.HomeTeam == club || fixture.AwayTeam == club {
				fixtures = append(fixtures, fixture)
			}
		}
	}
	return fixtures
}

// GetNextFixture finds the next fixture in the season that hasn't been played.
func (s *Season) GetNextFixture() (*Fixture, error) {
	for i := range s.Gameweeks {
		for j := range s.Gameweeks[i].Fixtures {
			f := s.Gameweeks[i].Fixtures[j]
			if f == nil {
				continue
			}

			if f.Result == nil {
				return f, nil
			}

			// A fixture is considered unplayed if its Result has no home team.
			if f.Result.Home == nil {
				return f, nil
			}
		}
	}
	return nil, fmt.Errorf("no unplayed fixtures remaining in the season")
}
