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
	BaseStaminaDrain   float64
}

var Positions = []Position{
	{Name: "Goalkeeper", AttackingEastZones: []PitchZone{}, AttackingWestZones: []PitchZone{}},
}
