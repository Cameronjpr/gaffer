package domain

import (
	"fmt"
	"sort"
)

type Season struct {
	ID       int
	Name     string
	Clubs    []*ClubWithPlayers
	Fixtures []*Fixture
}

func NewSeason(clubs []*ClubWithPlayers) *Season {
	season := Season{
		ID:       1,
		Name:     "2025/26",
		Clubs:    clubs,
		Fixtures: []*Fixture{},
	}

	return &season
}

func (s *Season) GenerateAllFixtures() {
	// Create 38 empty gameweeks
	for i := 1; i <= 38; i++ {
		fixtures := generateAllFixturesForGameweek(s.Clubs, i)
		s.Fixtures = append(s.Fixtures, fixtures...)
	}

}

func generateAllFixturesForGameweek(clubs []*ClubWithPlayers, gameweek int) []*Fixture {
	var fixtures []*Fixture
	fixtureID := 1

	// Generate all fixtures (each pair plays home and away)
	for i, homeClub := range clubs {
		for j, awayClub := range clubs {
			if i == j {
				continue
			}
			fixture := Fixture{
				ID:       fixtureID,
				Gameweek: gameweek,
				HomeTeam: homeClub,
				AwayTeam: awayClub,
				Result:   nil,
			}
			fixtures = append(fixtures, &fixture)
			fixtureID++
		}
	}

	return fixtures
}

func (s *Season) GetLeagueTable() LeagueTable {
	table := LeagueTable{
		Positions: []LeaguePosition{},
	}

	for _, club := range s.Clubs {
		// Get fixtures for this club from the season's fixtures
		fixtures := s.GetFixturesForClub(club.Club)

		var points int
		for _, fixture := range fixtures {
			// Skip fixtures that haven't been played yet
			if fixture.Result == nil {
				continue
			}

			winner := fixture.Result.GetWinner()
			if winner == nil {
				points += 1 // Draw
			} else if winner == club.Club {
				points += 3 // Win
			}
		}

		position := LeaguePosition{
			Club:   club.Club,
			Points: points,
		}
		table.Positions = append(table.Positions, position)
	}

	sort.Sort(ByLeagueStanding(table.Positions))

	return table
}

// GetFixturesForClub returns all fixtures (home and away) for a club, sorted by gameweek
func (s *Season) GetFixturesForClub(club *Club) []*Fixture {
	var fixtures []*Fixture
	for _, fixture := range s.Fixtures {
		if fixture.HomeTeam.Club == club || fixture.AwayTeam.Club == club {
			fixtures = append(fixtures, fixture)
		}
	}
	// Sort fixtures by gameweek
	sort.Slice(fixtures, func(i, j int) bool {
		return fixtures[i].Gameweek < fixtures[j].Gameweek
	})
	return fixtures
}

// GetNextFixture finds the next unplayed fixture in the season
func (s *Season) GetNextFixture() (*Fixture, error) {
	for _, fixture := range s.Fixtures {
		if fixture.Result == nil {
			return fixture, nil
		}
	}
	return nil, fmt.Errorf("no unplayed fixtures remaining in the season")
}

// GetNextFixtureForClub finds the next unplayed fixture for a specific club
func (s *Season) GetNextFixtureForClub(club *Club) (*Fixture, error) {
	clubFixtures := s.GetFixturesForClub(club)
	for _, fixture := range clubFixtures {
		if fixture.Result == nil {
			return fixture, nil
		}
	}
	return nil, fmt.Errorf("no unplayed fixtures remaining for club %s", club.Name)
}

// GetFixturesByGameweek returns all fixtures for a specific gameweek
func (s *Season) GetFixturesByGameweek(gameweek int) []*Fixture {
	var fixtures []*Fixture
	for _, fixture := range s.Fixtures {
		if fixture.Gameweek == gameweek {
			fixtures = append(fixtures, fixture)
		}
	}
	return fixtures
}
