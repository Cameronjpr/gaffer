package domain

type Fixture struct {
	ID       int
	HomeTeam *ClubWithPlayers
	AwayTeam *ClubWithPlayers
	Result   *Match
}
