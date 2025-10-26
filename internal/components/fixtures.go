package components

import (
	"fmt"

	"github.com/cameronjpr/gaffer/internal/domain"
)

func Fixtures(fixtures []*domain.Fixture) string {
	fixturesStr := "Fixtures:\n"

	numToShow := min(len(fixtures), 5)

	if numToShow == 0 {
		fixturesStr += "No fixtures scheduled\n"
	} else {
		for _, fixture := range fixtures[:numToShow] {
			fixturesStr += fmt.Sprintf("GW%d: %s vs %s\n", fixture.Gameweek, fixture.HomeTeam.Club.Name, fixture.AwayTeam.Club.Name)
		}
	}

	return fixturesStr
}
