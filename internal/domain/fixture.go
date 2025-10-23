package domain

type Fixture struct {
	ID       int
	HomeTeam *Club
	AwayTeam *Club
	Result   *Match
}
