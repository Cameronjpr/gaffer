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
	Clubs     []*ClubWithPlayers
	Gameweeks []*Gameweek
}

func NewSeason(clubs []*ClubWithPlayers) *Season {
	season := Season{
		ID:        1,
		Name:      "2025/26",
		Clubs:     clubs,
		Gameweeks: []*Gameweek{},
	}

	return &season
}

func (s *Season) GenerateGameweeks() {
	// Generate all fixtures for the season (each pair of teams plays home and away)
	allFixtures := generateAllFixtures(s.Clubs)

	// Create 38 empty gameweeks
	for i := 1; i <= 38; i++ {
		gameweek := Gameweek{
			ID:       i,
			Name:     fmt.Sprintf("Gameweek %d", i),
			Fixtures: []*Fixture{},
		}
		s.Gameweeks = append(s.Gameweeks, &gameweek)
	}

	// Distribute fixtures across gameweeks
	fixturesPerGameweek := len(allFixtures) / 38
	remainder := len(allFixtures) % 38

	fixtureIndex := 0
	for i := 0; i < 38; i++ {
		// Some gameweeks get one extra fixture if there's a remainder
		numFixtures := fixturesPerGameweek
		if i < remainder {
			numFixtures++
		}

		for j := 0; j < numFixtures && fixtureIndex < len(allFixtures); j++ {
			s.Gameweeks[i].Fixtures = append(s.Gameweeks[i].Fixtures, allFixtures[fixtureIndex])
			fixtureIndex++
		}
	}
}

func generateAllFixtures(clubs []*ClubWithPlayers) []*Fixture {
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

func (s *Season) GetFixturesForClub(club *Club) []*Fixture {
	var fixtures []*Fixture
	for _, gameweek := range s.Gameweeks {
		for _, fixture := range gameweek.Fixtures {
			if fixture.HomeTeam.Club == club || fixture.AwayTeam.Club == club {
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
