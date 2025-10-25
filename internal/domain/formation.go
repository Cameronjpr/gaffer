package domain

type Formation int

const (
	FourFourTwo Formation = iota
	FourThreeThree
)

type FormationData struct {
	Formation Formation
	Positions []Position
}

type Position struct {
	Name               string
	AttackingEastZones []PitchZone
	AttackingWestZones []PitchZone
}

var Positions = []Position{
	{Name: "Goalkeeper", AttackingEastZones: []PitchZone{}, AttackingWestZones: []PitchZone{}},
}
