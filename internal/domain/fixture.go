package domain

type Fixture struct {
	ID       int
	Gameweek int
	HomeTeam *ClubWithPlayers
	AwayTeam *ClubWithPlayers
	Result   *Match
}
